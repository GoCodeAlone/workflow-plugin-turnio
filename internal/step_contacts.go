package internal

import (
	"context"
	"encoding/json"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// checkContactStep implements step.turnio_check_contact.
type checkContactStep struct {
	name       string
	moduleName string
}

func newCheckContactStep(name string, config map[string]any) (*checkContactStep, error) {
	return &checkContactStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *checkContactStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	contacts := resolveStringSlice("contacts", current, config)
	if len(contacts) == 0 {
		phone := resolveValue("phone", current, config)
		if phone == "" {
			return &sdk.StepResult{Output: map[string]any{"error": "contacts or phone is required"}}, nil
		}
		contacts = []string{phone}
	}

	blocking := resolveValue("blocking", current, config)
	if blocking == "" {
		blocking = "wait"
	}

	payload := map[string]any{
		"blocking": blocking,
		"contacts": contacts,
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/contacts", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"contacts": string(result)}}, nil
}

// uploadContactsStep implements step.turnio_upload_contacts.
type uploadContactsStep struct {
	name       string
	moduleName string
}

func newUploadContactsStep(name string, config map[string]any) (*uploadContactsStep, error) {
	return &uploadContactsStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *uploadContactsStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	contacts := resolveStringSlice("contacts", current, config)
	if len(contacts) == 0 {
		return &sdk.StepResult{Output: map[string]any{"error": "contacts is required"}}, nil
	}

	payload := map[string]any{
		"blocking": "wait",
		"contacts": contacts,
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "POST", "/v1/contacts", payload, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": string(result)}}, nil
}

// updateProfileStep implements step.turnio_update_profile.
type updateProfileStep struct {
	name       string
	moduleName string
}

func newUpdateProfileStep(name string, config map[string]any) (*updateProfileStep, error) {
	return &updateProfileStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *updateProfileStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	profile := resolveMap("profile", current, config)
	if profile == nil {
		profile = map[string]any{}
		if desc := resolveValue("description", current, config); desc != "" {
			profile["description"] = desc
		}
		if email := resolveValue("email", current, config); email != "" {
			profile["email"] = email
		}
		if address := resolveValue("address", current, config); address != "" {
			profile["address"] = address
		}
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "PATCH", "/v1/settings/profile/business", profile, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": string(result)}}, nil
}
