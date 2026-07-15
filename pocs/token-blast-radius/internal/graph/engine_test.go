package graph

import (
	"strings"
	"testing"

	"token-blast-radius/internal/scenario"
	"token-blast-radius/internal/token"
)

func TestContainedAndWideScenarios(t *testing.T) {
	tokens, err := token.New(strings.Repeat("r", 32))
	if err != nil {
		t.Fatal(err)
	}
	scenarios, err := scenario.Build(tokens)
	if err != nil {
		t.Fatal(err)
	}
	engine := New(tokens, scenarios)
	contained, err := engine.Analyze(scenarios[0].ID, scenarios[0].Token)
	if err != nil {
		t.Fatal(err)
	}
	wide, err := engine.Analyze(scenarios[1].ID, scenarios[1].Token)
	if err != nil {
		t.Fatal(err)
	}
	if contained.Metrics.CriticalResources != 0 || contained.Metrics.RiskLevel != "contained" {
		t.Fatalf("unexpected contained result: %+v", contained.Metrics)
	}
	if wide.Metrics.CriticalResources < 4 || wide.Metrics.RiskLevel != "critical" {
		t.Fatalf("unexpected wide result: %+v", wide.Metrics)
	}
	if wide.Metrics.ReachableResources <= contained.Metrics.ReachableResources {
		t.Fatal("wide scenario must reach more resources")
	}
}

func TestModifiedTokenFails(t *testing.T) {
	tokens, _ := token.New(strings.Repeat("r", 32))
	scenarios, _ := scenario.Build(tokens)
	engine := New(tokens, scenarios)
	value := scenarios[0].Token
	value = value[:len(value)-1] + "A"
	if _, err := engine.Analyze(scenarios[0].ID, value); err == nil {
		t.Fatal("modified token must fail")
	}
}
