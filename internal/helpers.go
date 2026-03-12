package internal

import "fmt"

// getModuleName returns the "module" key from a step config map, defaulting to "turnio".
func getModuleName(config map[string]any) string {
	if v, ok := config["module"].(string); ok && v != "" {
		return v
	}
	return "turnio"
}

// resolveValue looks up key in current first, then config.
// Returns "" if not found.
func resolveValue(key string, current, config map[string]any) string {
	if v, ok := current[key].(string); ok && v != "" {
		return v
	}
	if v, ok := config[key].(string); ok && v != "" {
		return v
	}
	return ""
}

// resolveInt64 looks up key in current first, then config as int64.
func resolveInt64(key string, current, config map[string]any) int64 {
	if v := toInt64(current[key]); v != 0 {
		return v
	}
	return toInt64(config[key])
}

// resolveFloat64 looks up key in current first, then config as float64.
func resolveFloat64(key string, current, config map[string]any) float64 {
	if v := toFloat64(current[key]); v != 0 {
		return v
	}
	return toFloat64(config[key])
}

// resolveStringSlice looks up key in current first, then config as []string.
func resolveStringSlice(key string, current, config map[string]any) []string {
	if v, ok := current[key].([]string); ok {
		return v
	}
	if v, ok := current[key].([]any); ok {
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	if v, ok := config[key].([]string); ok {
		return v
	}
	if v, ok := config[key].([]any); ok {
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}

// resolveMap looks up key in current first, then config as map[string]any.
func resolveMap(key string, current, config map[string]any) map[string]any {
	if v, ok := current[key].(map[string]any); ok {
		return v
	}
	if v, ok := config[key].(map[string]any); ok {
		return v
	}
	return nil
}

// resolveBool looks up key in current first, then config as bool.
func resolveBool(key string, current, config map[string]any) bool {
	if v, ok := current[key].(bool); ok {
		return v
	}
	if v, ok := config[key].(bool); ok {
		return v
	}
	return false
}

// resolveInt looks up key in current first, then config as int.
func resolveInt(key string, current, config map[string]any) int {
	return int(resolveInt64(key, current, config))
}

func toInt64(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case int32:
		return int64(t)
	case float64:
		return int64(t)
	case float32:
		return int64(t)
	case string:
		var n int64
		fmt.Sscanf(t, "%d", &n)
		return n
	}
	return 0
}

func toFloat64(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int64:
		return float64(t)
	case int:
		return float64(t)
	case string:
		var f float64
		fmt.Sscanf(t, "%f", &f)
		return f
	}
	return 0
}
