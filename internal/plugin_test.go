package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestModuleSchemas verifies that the plugin implements sdk.SchemaProvider
// and returns a correct contract descriptor for the turnio.provider module.
func TestModuleSchemas(t *testing.T) {
	p := &turnIOPlugin{}
	schemas := p.ModuleSchemas()

	if len(schemas) != 1 {
		t.Fatalf("ModuleSchemas: expected 1 schema, got %d", len(schemas))
	}

	s := schemas[0]
	if s.Type != "turnio.provider" {
		t.Errorf("ModuleSchemas[0].Type = %q, want %q", s.Type, "turnio.provider")
	}
	if s.Label == "" {
		t.Error("ModuleSchemas[0].Label must not be empty")
	}
	if s.Category == "" {
		t.Error("ModuleSchemas[0].Category must not be empty")
	}
	if s.Description == "" {
		t.Error("ModuleSchemas[0].Description must not be empty")
	}

	// Verify required apiToken field is present.
	var foundAPIToken, foundBaseURL bool
	for _, cf := range s.ConfigFields {
		switch cf.Name {
		case "apiToken":
			foundAPIToken = true
			if !cf.Required {
				t.Error("apiToken config field must be marked Required")
			}
			if cf.Type != "string" {
				t.Errorf("apiToken type = %q, want %q", cf.Type, "string")
			}
		case "baseUrl":
			foundBaseURL = true
			if cf.Required {
				t.Error("baseUrl config field must not be Required")
			}
			if cf.DefaultValue == "" {
				t.Error("baseUrl must have a DefaultValue")
			}
		}
	}
	if !foundAPIToken {
		t.Error("ModuleSchemas[0].ConfigFields must include apiToken")
	}
	if !foundBaseURL {
		t.Error("ModuleSchemas[0].ConfigFields must include baseUrl")
	}
}

// TestPluginManifestStepSchemas verifies that plugin.json contains a stepSchema
// entry for every step type advertised by the plugin, ensuring no step is missing
// a contract descriptor.
func TestPluginManifestStepSchemas(t *testing.T) {
	// Locate plugin.json relative to the repo root.
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// internal/ → repo root
	repoRoot := filepath.Join(filepath.Dir(file), "..")
	manifestPath := filepath.Join(repoRoot, "plugin.json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read plugin.json: %v", err)
	}

	var manifest struct {
		StepTypes   []string `json:"stepTypes"`
		ModuleTypes []string `json:"moduleTypes"`
		StepSchemas []struct {
			Type string `json:"type"`
		} `json:"stepSchemas"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse plugin.json: %v", err)
	}

	// Build a set of schema types.
	schemaTypes := make(map[string]bool, len(manifest.StepSchemas))
	for _, s := range manifest.StepSchemas {
		schemaTypes[s.Type] = true
	}

	// Every stepType must have a matching stepSchema.
	for _, st := range manifest.StepTypes {
		if !schemaTypes[st] {
			t.Errorf("step type %q is advertised in stepTypes but has no entry in stepSchemas", st)
		}
	}

	// Every stepSchema must correspond to an advertised stepType.
	advertised := make(map[string]bool, len(manifest.StepTypes))
	for _, st := range manifest.StepTypes {
		advertised[st] = true
	}
	for _, s := range manifest.StepSchemas {
		if !advertised[s.Type] {
			t.Errorf("stepSchema type %q is not listed in stepTypes", s.Type)
		}
	}

	if len(manifest.ModuleTypes) == 0 {
		t.Error("plugin.json must list at least one moduleType")
	}
}

// TestPluginManifestSchemaFields verifies that each step schema in plugin.json
// has at least one configField and at least one output, which are required for
// a valid strict contract descriptor.
func TestPluginManifestSchemaFields(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Join(filepath.Dir(file), "..")
	manifestPath := filepath.Join(repoRoot, "plugin.json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read plugin.json: %v", err)
	}

	var manifest struct {
		StepSchemas []struct {
			Type         string `json:"type"`
			Description  string `json:"description"`
			ConfigFields []struct {
				Key  string `json:"key"`
				Type string `json:"type"`
			} `json:"configFields"`
			Outputs []struct {
				Key  string `json:"key"`
				Type string `json:"type"`
			} `json:"outputs"`
		} `json:"stepSchemas"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("parse plugin.json: %v", err)
	}

	for _, s := range manifest.StepSchemas {
		if s.Description == "" {
			t.Errorf("step schema %q must have a non-empty description", s.Type)
		}
		if len(s.ConfigFields) == 0 {
			t.Errorf("step schema %q must declare at least one configField", s.Type)
		}
		if len(s.Outputs) == 0 {
			t.Errorf("step schema %q must declare at least one output", s.Type)
		}
		for _, cf := range s.ConfigFields {
			if cf.Key == "" {
				t.Errorf("step schema %q has a configField with an empty key", s.Type)
			}
			if cf.Type == "" {
				t.Errorf("step schema %q configField %q has an empty type", s.Type, cf.Key)
			}
		}
		for _, o := range s.Outputs {
			if o.Key == "" {
				t.Errorf("step schema %q has an output with an empty key", s.Type)
			}
			if o.Type == "" {
				t.Errorf("step schema %q output %q has an empty type", s.Type, o.Key)
			}
		}
	}
}
