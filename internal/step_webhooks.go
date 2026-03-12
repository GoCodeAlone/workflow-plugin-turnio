package internal

import (
	"context"
	"encoding/json"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// configureWebhookStep implements step.turnio_configure_webhook.
type configureWebhookStep struct {
	name       string
	moduleName string
}

func newConfigureWebhookStep(name string, config map[string]any) (*configureWebhookStep, error) {
	return &configureWebhookStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *configureWebhookStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	webhookURL := resolveValue("webhook_url", current, config)
	if webhookURL == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "webhook_url is required"}}, nil
	}

	payload := map[string]any{
		"webhooks": []map[string]any{
			{"url": webhookURL},
		},
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/configs/webhook", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": string(result), "webhook_url": webhookURL}}, nil
}
