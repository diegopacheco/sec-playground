package pin

import (
	"encoding/json"
	"testing"

	"mcp-rug-pull-proxy/internal/mcp"
)

func TestHashIgnoresSchemaKeyOrderSoHarmlessReserializationIsNotARugPull(t *testing.T) {
	a := mcp.Tool{Name: "get_weather", Description: "Returns the weather.", InputSchema: json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"}}}`)}
	b := mcp.Tool{Name: "get_weather", Description: "Returns the weather.", InputSchema: json.RawMessage(`{"properties":{"city":{"type":"string"}},"type":"object"}`)}
	if Hash(a) != Hash(b) {
		t.Fatalf("reordered keys changed the hash: %s vs %s", Hash(a), Hash(b))
	}
}

func TestHashChangesOnSingleCharacterEditSoSilentDescriptionSwapsAreCaught(t *testing.T) {
	a := mcp.Tool{Name: "get_weather", Description: "Returns the weather."}
	b := mcp.Tool{Name: "get_weather", Description: "Returns the weather!"}
	if Hash(a) == Hash(b) {
		t.Fatal("a one character description edit kept the same hash")
	}
}

func TestHashCoversSchemaAndAnnotationsNotOnlyTheDescription(t *testing.T) {
	base := mcp.Tool{Name: "add_numbers", Description: "Adds numbers.", InputSchema: json.RawMessage(`{"type":"object"}`)}
	withParam := base
	withParam.InputSchema = json.RawMessage(`{"type":"object","properties":{"sidenote":{"type":"string"}}}`)
	withAnnotation := base
	withAnnotation.Annotations = json.RawMessage(`{"readOnlyHint":true}`)
	if Hash(base) == Hash(withParam) {
		t.Fatal("adding a hidden parameter did not change the hash")
	}
	if Hash(base) == Hash(withAnnotation) {
		t.Fatal("changing annotations did not change the hash")
	}
}
