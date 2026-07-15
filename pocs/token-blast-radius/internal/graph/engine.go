package graph

import (
	"errors"
	"fmt"
	"sort"

	"token-blast-radius/internal/domain"
	"token-blast-radius/internal/token"
)

type Engine struct {
	tokens    *token.Service
	scenarios map[string]domain.Scenario
}

type parentStep struct {
	from       string
	permission string
}

func New(tokens *token.Service, scenarios []domain.Scenario) *Engine {
	indexed := make(map[string]domain.Scenario, len(scenarios))
	for _, scenario := range scenarios {
		indexed[scenario.ID] = scenario
	}
	return &Engine{tokens: tokens, scenarios: indexed}
}

func (e *Engine) Summaries() []domain.ScenarioSummary {
	result := make([]domain.ScenarioSummary, 0, len(e.scenarios))
	for _, scenario := range e.scenarios {
		result = append(result, domain.ScenarioSummary{
			ID: scenario.ID, Name: scenario.Name, Summary: scenario.Summary,
			Question: scenario.Question, Token: scenario.Token, Claims: scenario.Claims,
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (e *Engine) Analyze(scenarioID string, value string) (domain.Analysis, error) {
	scenario, exists := e.scenarios[scenarioID]
	if !exists {
		return domain.Analysis{}, errors.New("scenario does not exist")
	}
	claims, err := e.tokens.Parse(value)
	if err != nil {
		return domain.Analysis{}, err
	}
	if claims.Audience != scenario.Claims.Audience {
		return domain.Analysis{}, errors.New("token audience does not match the selected system")
	}
	if claims.Subject != scenario.Claims.Subject {
		return domain.Analysis{}, errors.New("token subject does not exist in the selected graph")
	}
	return calculate(scenario, claims), nil
}

func calculate(scenario domain.Scenario, claims domain.Claims) domain.Analysis {
	scopes := make(map[string]bool, len(claims.Scopes))
	for _, scope := range claims.Scopes {
		scopes[scope] = true
	}
	depths := map[string]int{claims.Subject: 0, claims.Actor: 0}
	parents := make(map[string]parentStep)
	queue := []string{claims.Subject}
	active := make(map[int]bool)
	blocked := make(map[int]string)
	maximumReached := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		depth := depths[current]
		for index, relationship := range scenario.Relationships {
			if relationship.From != current {
				continue
			}
			if depth >= claims.MaximumDepth {
				blocked[index] = "token delegation depth reached"
				continue
			}
			if relationship.RequiredScope != "" && !scopes[relationship.RequiredScope] {
				blocked[index] = "missing " + relationship.RequiredScope + " scope"
				continue
			}
			active[index] = true
			nextDepth := depth + 1
			if existing, visited := depths[relationship.To]; visited && existing <= nextDepth {
				continue
			}
			depths[relationship.To] = nextDepth
			parents[relationship.To] = parentStep{from: current, permission: relationship.Permission}
			queue = append(queue, relationship.To)
			if nextDepth > maximumReached {
				maximumReached = nextDepth
			}
		}
	}
	nodes := make([]domain.AnalyzedNode, 0, len(scenario.Nodes))
	reachableResources := 0
	criticalResources := 0
	for _, node := range scenario.Nodes {
		depth, reachable := depths[node.ID]
		nodes = append(nodes, domain.AnalyzedNode{Node: node, Reachable: reachable, Depth: depth})
		if reachable && isResource(node.Type) {
			reachableResources++
			if node.Sensitivity == "critical" {
				criticalResources++
			}
		}
	}
	edges := make([]domain.AnalyzedRelationship, 0, len(scenario.Relationships)+1)
	edges = append(edges, domain.AnalyzedRelationship{
		Relationship: domain.Relationship{From: claims.Actor, To: claims.Subject, Permission: "acts for", Rationale: "signed actor chain in the token"},
		Active:       true,
	})
	activeDelegations := 0
	privilegedEdges := 0
	for index, relationship := range scenario.Relationships {
		reason := blocked[index]
		if reason == "" && !active[index] {
			reason = "source is outside the reachable graph"
		}
		edges = append(edges, domain.AnalyzedRelationship{Relationship: relationship, Active: active[index], BlockedReason: reason})
		if active[index] {
			activeDelegations++
			if isPrivileged(relationship.Permission) {
				privilegedEdges++
			}
		}
	}
	paths := buildPaths(nodes, parents, claims)
	score := reachableResources*5 + criticalResources*18 + activeDelegations*2 + privilegedEdges*9 + maximumReached*3
	if score > 100 {
		score = 100
	}
	level, verdict := risk(score, reachableResources, criticalResources)
	return domain.Analysis{
		ScenarioID: scenario.ID, Claims: claims, Nodes: nodes, Edges: edges, Paths: paths,
		Metrics: domain.Metrics{ReachableResources: reachableResources, CriticalResources: criticalResources, ActiveDelegations: activeDelegations, MaximumDepth: maximumReached, RiskScore: score, RiskLevel: level},
		Verdict: verdict,
	}
}

func buildPaths(nodes []domain.AnalyzedNode, parents map[string]parentStep, claims domain.Claims) []domain.AccessPath {
	paths := make([]domain.AccessPath, 0)
	for _, node := range nodes {
		if !node.Reachable || node.ID == claims.Subject || node.ID == claims.Actor {
			continue
		}
		ids := []string{node.ID}
		permissions := make([]string, 0)
		current := node.ID
		for current != claims.Subject {
			step, exists := parents[current]
			if !exists {
				break
			}
			ids = append(ids, step.from)
			permissions = append(permissions, step.permission)
			current = step.from
		}
		reverse(ids)
		reverse(permissions)
		paths = append(paths, domain.AccessPath{Target: node.ID, NodeIDs: ids, Permissions: permissions})
	}
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i].NodeIDs) == len(paths[j].NodeIDs) {
			return paths[i].Target < paths[j].Target
		}
		return len(paths[i].NodeIDs) < len(paths[j].NodeIDs)
	})
	return paths
}

func reverse(values []string) {
	for left, right := 0, len(values)-1; left < right; left, right = left+1, right-1 {
		values[left], values[right] = values[right], values[left]
	}
}

func isResource(kind string) bool {
	return kind == "document" || kind == "vault" || kind == "service" || kind == "account" || kind == "database"
}

func isPrivileged(permission string) bool {
	return permission == "admin" || permission == "export" || permission == "impersonate" || permission == "reset" || permission == "delegate"
}

func risk(score int, reachable int, critical int) (string, string) {
	if score >= 65 {
		return "critical", fmt.Sprintf("This credential reaches %d resources, including %d critical targets. Revoke it and reduce scopes or delegation depth.", reachable, critical)
	}
	if score >= 35 {
		return "elevated", fmt.Sprintf("This credential crosses several trust boundaries and reaches %d resources. Remove unused paths before issuance.", reachable)
	}
	return "contained", fmt.Sprintf("This credential reaches %d resources through a narrow, explainable path.", reachable)
}
