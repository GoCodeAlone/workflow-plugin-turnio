package internal

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// createTemplateStep implements step.turnio_create_template.
type createTemplateStep struct {
	name       string
	moduleName string
}

func newCreateTemplateStep(name string, config map[string]any) (*createTemplateStep, error) {
	return &createTemplateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *createTemplateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	templateName := resolveValue("template_name", current, config)
	if templateName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "template_name is required"}}, nil
	}

	payload := map[string]any{"name": templateName}
	if category := resolveValue("category", current, config); category != "" {
		payload["category"] = category
	}
	if langCode := resolveValue("language_code", current, config); langCode != "" {
		payload["language"] = langCode
	}
	if components := resolveMap("components", current, config); components != nil {
		payload["components"] = components
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/configs/templates", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"template": string(result)}}, nil
}

// listTemplatesStep implements step.turnio_list_templates.
type listTemplatesStep struct {
	name       string
	moduleName string
}

func newListTemplatesStep(name string, config map[string]any) (*listTemplatesStep, error) {
	return &listTemplatesStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listTemplatesStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "GET", "/v1/configs/templates", nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"templates": string(result)}}, nil
}

// fetchTemplateStep implements step.turnio_fetch_template.
type fetchTemplateStep struct {
	name       string
	moduleName string
}

func newFetchTemplateStep(name string, config map[string]any) (*fetchTemplateStep, error) {
	return &fetchTemplateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *fetchTemplateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	templateName := resolveValue("template_name", current, config)
	if templateName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "template_name is required"}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "GET", fmt.Sprintf("/v1/configs/templates/%s", templateName), nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"template": string(result)}}, nil
}

// updateTemplateStep implements step.turnio_update_template.
type updateTemplateStep struct {
	name       string
	moduleName string
}

func newUpdateTemplateStep(name string, config map[string]any) (*updateTemplateStep, error) {
	return &updateTemplateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateTemplateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	templateName := resolveValue("template_name", current, config)
	if templateName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "template_name is required"}}, nil
	}

	payload := resolveMap("updates", current, config)
	if payload == nil {
		payload = map[string]any{}
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "PATCH", fmt.Sprintf("/v1/configs/templates/%s", templateName), payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"template": string(result)}}, nil
}

// deleteTemplateStep implements step.turnio_delete_template.
type deleteTemplateStep struct {
	name       string
	moduleName string
}

func newDeleteTemplateStep(name string, config map[string]any) (*deleteTemplateStep, error) {
	return &deleteTemplateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteTemplateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	templateName := resolveValue("template_name", current, config)
	if templateName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "template_name is required"}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "DELETE", fmt.Sprintf("/v1/configs/templates/%s", templateName), nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "template_name": templateName}}, nil
}
