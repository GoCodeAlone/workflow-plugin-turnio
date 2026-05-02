package internal

import "testing"

func TestResolveIntDefaultCurrentZeroOverridesConfig(t *testing.T) {
	got := resolveIntDefault("limit", 3, map[string]any{"limit": 0}, map[string]any{"limit": 7})
	if got != 0 {
		t.Fatalf("expected explicit current zero to override config, got %d", got)
	}
}

func TestResolveIntDefaultFallbacks(t *testing.T) {
	tests := []struct {
		name    string
		current map[string]any
		config  map[string]any
		want    int
	}{
		{
			name:    "absent returns default",
			current: map[string]any{},
			config:  map[string]any{},
			want:    3,
		},
		{
			name:    "config used when current absent",
			current: map[string]any{},
			config:  map[string]any{"limit": 7},
			want:    7,
		},
		{
			name:    "invalid current falls back to config",
			current: map[string]any{"limit": ""},
			config:  map[string]any{"limit": 7},
			want:    7,
		},
		{
			name:    "invalid values return default",
			current: map[string]any{"limit": "bad"},
			config:  map[string]any{"limit": "also-bad"},
			want:    3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveIntDefault("limit", 3, tt.current, tt.config)
			if got != tt.want {
				t.Fatalf("resolveIntDefault() = %d, want %d", got, tt.want)
			}
		})
	}
}
