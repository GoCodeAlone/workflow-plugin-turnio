package internal

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// sendTextStep implements step.turnio_send_text.
type sendTextStep struct {
	name       string
	moduleName string
}

func newSendTextStep(name string, config map[string]any) (*sendTextStep, error) {
	return &sendTextStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sendTextStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}
	to := resolveValue("to", current, config)
	if to == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "to is required"}}, nil
	}
	body := resolveValue("body", current, config)
	if body == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "body is required"}}, nil
	}

	payload := map[string]any{
		"to":   to,
		"type": "text",
		"text": map[string]any{"body": body},
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/messages", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"response": string(result), "to": to}}, nil
}

// sendMediaStep implements step.turnio_send_media.
type sendMediaStep struct {
	name       string
	moduleName string
}

func newSendMediaStep(name string, config map[string]any) (*sendMediaStep, error) {
	return &sendMediaStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sendMediaStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}
	to := resolveValue("to", current, config)
	if to == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "to is required"}}, nil
	}
	mediaType := resolveValue("media_type", current, config)
	if mediaType == "" {
		mediaType = "image"
	}
	link := resolveValue("link", current, config)
	if link == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "link is required"}}, nil
	}

	mediaObj := map[string]any{"link": link}
	if caption := resolveValue("caption", current, config); caption != "" {
		mediaObj["caption"] = caption
	}

	payload := map[string]any{
		"to":      to,
		"type":    mediaType,
		mediaType: mediaObj,
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/messages", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"response": string(result), "to": to}}, nil
}

// sendTemplateStep implements step.turnio_send_template.
type sendTemplateStep struct {
	name       string
	moduleName string
}

func newSendTemplateStep(name string, config map[string]any) (*sendTemplateStep, error) {
	return &sendTemplateStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sendTemplateStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}
	to := resolveValue("to", current, config)
	if to == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "to is required"}}, nil
	}
	templateName := resolveValue("template_name", current, config)
	if templateName == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "template_name is required"}}, nil
	}
	namespace := resolveValue("namespace", current, config)
	langCode := resolveValue("language_code", current, config)
	if langCode == "" {
		langCode = "en"
	}

	tmpl := map[string]any{
		"name":     templateName,
		"language": map[string]any{"code": langCode},
	}
	if namespace != "" {
		tmpl["namespace"] = namespace
	}
	if components := resolveStringSlice("components", current, config); len(components) > 0 {
		tmpl["components"] = components
	} else if comp := resolveMap("components", current, config); comp != nil {
		tmpl["components"] = comp
	}

	payload := map[string]any{
		"to":       to,
		"type":     "template",
		"template": tmpl,
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/messages", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"response": string(result), "to": to}}, nil
}

// sendInteractiveStep implements step.turnio_send_interactive.
type sendInteractiveStep struct {
	name       string
	moduleName string
}

func newSendInteractiveStep(name string, config map[string]any) (*sendInteractiveStep, error) {
	return &sendInteractiveStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sendInteractiveStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}
	to := resolveValue("to", current, config)
	if to == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "to is required"}}, nil
	}
	interactive := resolveMap("interactive", current, config)
	if interactive == nil {
		return &sdk.StepResult{Output: map[string]any{"error": "interactive is required"}}, nil
	}

	payload := map[string]any{
		"to":          to,
		"type":        "interactive",
		"interactive": interactive,
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/messages", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"response": string(result), "to": to}}, nil
}

// sendLocationStep implements step.turnio_send_location.
type sendLocationStep struct {
	name       string
	moduleName string
}

func newSendLocationStep(name string, config map[string]any) (*sendLocationStep, error) {
	return &sendLocationStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *sendLocationStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}
	to := resolveValue("to", current, config)
	if to == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "to is required"}}, nil
	}
	longitude := resolveFloat64("longitude", current, config)
	latitude := resolveFloat64("latitude", current, config)

	loc := map[string]any{
		"longitude": longitude,
		"latitude":  latitude,
	}
	if name := resolveValue("location_name", current, config); name != "" {
		loc["name"] = name
	}
	if address := resolveValue("address", current, config); address != "" {
		loc["address"] = address
	}

	payload := map[string]any{
		"to":       to,
		"type":     "location",
		"location": loc,
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/messages", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"response": string(result), "to": to}}, nil
}

// listMessagesStep implements step.turnio_list_messages.
type listMessagesStep struct {
	name       string
	moduleName string
}

func newListMessagesStep(name string, config map[string]any) (*listMessagesStep, error) {
	return &listMessagesStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *listMessagesStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	path := "/v1/messages"
	contactID := resolveValue("contact_id", current, config)
	if contactID != "" {
		path = fmt.Sprintf("/v1/contacts/%s/messages", contactID)
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "GET", path, nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"messages": string(result)}}, nil
}
