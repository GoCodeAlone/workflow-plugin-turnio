package internal

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// createFlowStep implements step.turnio_create_flow.
type createFlowStep struct {
	name       string
	moduleName string
}

func newCreateFlowStep(name string, config map[string]any) (*createFlowStep, error) {
	return &createFlowStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createFlowStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	flowName := resolveValue("flow_name", current, config)
	if flowName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "flow_name is required"}}, nil
	}

	payload := map[string]any{"name": flowName}
	if desc := resolveValue("description", current, config); desc != "" {
		payload["description"] = desc
	}
	if steps := resolveMap("steps", current, config); steps != nil {
		payload["steps"] = steps
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/flows", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"flow": string(result)}}, nil
}

// listFlowsStep implements step.turnio_list_flows.
type listFlowsStep struct {
	name       string
	moduleName string
}

func newListFlowsStep(name string, config map[string]any) (*listFlowsStep, error) {
	return &listFlowsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listFlowsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "GET", "/v1/flows", nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"flows": string(result)}}, nil
}

// sendFlowStep implements step.turnio_send_flow.
type sendFlowStep struct {
	name       string
	moduleName string
}

func newSendFlowStep(name string, config map[string]any) (*sendFlowStep, error) {
	return &sendFlowStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sendFlowStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	to := resolveValue("to", current, config)
	if to == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "to is required"}}, nil
	}
	flowID := resolveValue("flow_id", current, config)
	if flowID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "flow_id is required"}}, nil
	}

	payload := map[string]any{
		"to":   to,
		"type": "interactive",
		"interactive": map[string]any{
			"type": "flow",
			"action": map[string]any{
				"name": "flow",
				"parameters": map[string]any{
					"flow_id":              flowID,
					"flow_cta":             resolveValue("flow_cta", current, config),
					"flow_action":          resolveValue("flow_action", current, config),
					"flow_message_version": fmt.Sprintf("%d", resolveInt("flow_message_version", current, config)),
				},
			},
		},
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/messages", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"response": string(result), "to": to, "flow_id": flowID}}, nil
}
