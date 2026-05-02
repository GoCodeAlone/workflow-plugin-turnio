package internal

import "testing"

func TestResolveIntDefaultCurrentZeroOverridesConfig(t *testing.T) {
	got := resolveIntDefault("limit", 3, map[string]any{"limit": 0}, map[string]any{"limit": 7})
	if got != 0 {
		t.Fatalf("expected explicit current zero to override config, got %d", got)
	}
}
