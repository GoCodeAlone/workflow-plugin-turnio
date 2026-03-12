package internal

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// listJourneysStep implements step.turnio_list_journeys.
type listJourneysStep struct {
	name       string
	moduleName string
}

func newListJourneysStep(name string, config map[string]any) (*listJourneysStep, error) {
	return &listJourneysStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listJourneysStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "GET", "/v1/journeys", nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"journeys": string(result)}}, nil
}

// triggerJourneyStep implements step.turnio_trigger_journey.
type triggerJourneyStep struct {
	name       string
	moduleName string
}

func newTriggerJourneyStep(name string, config map[string]any) (*triggerJourneyStep, error) {
	return &triggerJourneyStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *triggerJourneyStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	journeyID := resolveValue("journey_id", current, config)
	if journeyID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "journey_id is required"}}, nil
	}
	phone := resolveValue("phone", current, config)
	if phone == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "phone is required"}}, nil
	}

	payload := map[string]any{
		"contact": map[string]any{"phone_number": phone},
	}
	if data := resolveMap("data", current, config); data != nil {
		payload["data"] = data
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", fmt.Sprintf("/v1/journeys/%s/trigger", journeyID), payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": string(result), "journey_id": journeyID}}, nil
}
