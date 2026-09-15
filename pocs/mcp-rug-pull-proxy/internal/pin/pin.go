package pin

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"mcp-rug-pull-proxy/internal/mcp"
)

func Hash(t mcp.Tool) string {
	canonical := map[string]any{
		"name":         t.Name,
		"title":        t.Title,
		"description":  t.Description,
		"inputSchema":  decode(t.InputSchema),
		"outputSchema": decode(t.OutputSchema),
		"annotations":  decode(t.Annotations),
	}
	raw, _ := json.Marshal(canonical)
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func decode(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	return v
}

func Short(hash string) string {
	if len(hash) > 19 {
		return hash[:19]
	}
	return hash
}
