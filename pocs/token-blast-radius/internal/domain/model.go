package domain

type Claims struct {
	Issuer       string   `json:"iss"`
	Subject      string   `json:"sub"`
	Actor        string   `json:"act"`
	Audience     string   `json:"aud"`
	Scopes       []string `json:"scope"`
	MaximumDepth int      `json:"maximum_depth"`
	ExpiresAt    int64    `json:"exp"`
	ID           string   `json:"jti"`
}

type Node struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Sensitivity string `json:"sensitivity"`
	Description string `json:"description"`
}

type Relationship struct {
	From          string `json:"from"`
	To            string `json:"to"`
	Permission    string `json:"permission"`
	RequiredScope string `json:"required_scope"`
	Rationale     string `json:"rationale"`
}

type Scenario struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Summary       string         `json:"summary"`
	Question      string         `json:"question"`
	Token         string         `json:"token"`
	Claims        Claims         `json:"claims"`
	Nodes         []Node         `json:"nodes"`
	Relationships []Relationship `json:"relationships"`
}

type ScenarioSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Summary  string `json:"summary"`
	Question string `json:"question"`
	Token    string `json:"token"`
	Claims   Claims `json:"claims"`
}

type AnalysisRequest struct {
	ScenarioID string `json:"scenario_id"`
	Token      string `json:"token"`
}

type AnalyzedNode struct {
	Node
	Reachable bool `json:"reachable"`
	Depth     int  `json:"depth"`
}

type AnalyzedRelationship struct {
	Relationship
	Active        bool   `json:"active"`
	BlockedReason string `json:"blocked_reason,omitempty"`
}

type AccessPath struct {
	Target      string   `json:"target"`
	NodeIDs     []string `json:"node_ids"`
	Permissions []string `json:"permissions"`
}

type Metrics struct {
	ReachableResources int    `json:"reachable_resources"`
	CriticalResources  int    `json:"critical_resources"`
	ActiveDelegations  int    `json:"active_delegations"`
	MaximumDepth       int    `json:"maximum_depth"`
	RiskScore          int    `json:"risk_score"`
	RiskLevel          string `json:"risk_level"`
}

type Analysis struct {
	ScenarioID string                 `json:"scenario_id"`
	Claims     Claims                 `json:"claims"`
	Nodes      []AnalyzedNode         `json:"nodes"`
	Edges      []AnalyzedRelationship `json:"edges"`
	Paths      []AccessPath           `json:"paths"`
	Metrics    Metrics                `json:"metrics"`
	Verdict    string                 `json:"verdict"`
}
