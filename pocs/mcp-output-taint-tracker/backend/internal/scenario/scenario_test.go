package scenario

import (
	"testing"

	"mcp-output-taint-tracker/internal/policy"
	"mcp-output-taint-tracker/internal/session"
)

func TestEveryScenarioBlocksExactlyTheAttackAndAllowsTheRest(t *testing.T) {
	store := session.NewStore(policy.Default())
	for _, sc := range All() {
		_, events := Run(store, sc)
		for i, st := range sc.Steps {
			if st.Kind != "call" {
				continue
			}
			d := events[i].Decision
			if st.Expect == "" {
				if !d.Allow {
					t.Errorf("%q step %d: legit call %s blocked: %+v", sc.Name, i+1, st.Tool, d.Flows)
				}
				continue
			}
			if d.Allow {
				t.Errorf("%q step %d: attack %s allowed", sc.Name, i+1, st.Tool)
				continue
			}
			found := false
			for _, f := range d.Flows {
				found = found || f.Label == st.Expect
			}
			if !found {
				t.Errorf("%q step %d: wanted a %s flow, got %+v", sc.Name, i+1, st.Expect, d.Flows)
			}
		}
	}
}

func TestBlockedFlowPointsBackToTheToolResultThatCausedIt(t *testing.T) {
	store := session.NewStore(policy.Default())
	s, _ := Run(store, All()[0])
	last := s.Events[len(s.Events)-1]
	src := s.Events[last.Decision.Flows[0].SourceID-1]
	if src.Kind != "result" || src.Tool != "WebFetch" {
		t.Fatalf("flow source must be the fetched page, got %+v", src)
	}
}
