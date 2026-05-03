package internal

import "strconv"

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

// hasKey reports whether key exists in current or config (regardless of value).
func hasKey(key string, current, config map[string]any) bool {
	if _, ok := current[key]; ok {
		return true
	}
	_, ok := config[key]
	return ok
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
	if v, ok := resolveIntValue(key, current, config); ok {
		return v
	}
	return 0
}

// resolveIntDefault looks up key in current first, then config as int.
// If the key is absent from both maps, def is returned.
func resolveIntDefault(key string, def int, current, config map[string]any) int {
	if v, ok := current[key]; ok {
		if n, valid := toInt(v); valid {
			return n
		}
	}
	if v, ok := config[key]; ok {
		if n, valid := toInt(v); valid {
			return n
		}
	}
	return def
}

func resolveIntValue(key string, current, config map[string]any) (int, bool) {
	if v, ok := current[key]; ok {
		if n, valid := toInt(v); valid {
			return n, true
		}
	}
	if v, ok := config[key]; ok {
		if n, valid := toInt(v); valid {
			return n, true
		}
	}
	return 0, false
}

func toInt(v any) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int64:
		return atoiOK(strconv.FormatInt(t, 10))
	case int32:
		return atoiOK(strconv.FormatInt(int64(t), 10))
	case float64:
		return atoiOK(strconv.FormatInt(int64(t), 10))
	case float32:
		return atoiOK(strconv.FormatInt(int64(t), 10))
	case string:
		return atoiOK(t)
	}
	return 0, false
}

func atoiOK(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	return n, err == nil
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
		f, _ := strconv.ParseFloat(t, 64)
		return f
	}
	return 0
}
