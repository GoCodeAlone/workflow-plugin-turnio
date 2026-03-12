package internal

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// getContextStep implements step.turnio_get_context.
type getContextStep struct {
	name       string
	moduleName string
}

func newGetContextStep(name string, config map[string]any) (*getContextStep, error) {
	return &getContextStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *getContextStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	contactID := resolveValue("contact_id", current, config)
	if contactID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "contact_id is required"}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "GET", fmt.Sprintf("/v1/contacts/%s/context", contactID), nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"context": string(result), "contact_id": contactID}}, nil
}

// setContextStep implements step.turnio_set_context.
type setContextStep struct {
	name       string
	moduleName string
}

func newSetContextStep(name string, config map[string]any) (*setContextStep, error) {
	return &setContextStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *setContextStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	contactID := resolveValue("contact_id", current, config)
	if contactID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "contact_id is required"}}, nil
	}
	contextData := resolveMap("context_data", current, config)
	if contextData == nil {
		return &sdk.StepResult{Output: map[string]any{"error": "context_data is required"}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "PATCH", fmt.Sprintf("/v1/contacts/%s/context", contactID), contextData, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": string(result), "contact_id": contactID}}, nil
}
