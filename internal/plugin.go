// Package internal implements the workflow-plugin-turnio plugin.
package internal

import (
	"fmt"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// turnIOPlugin implements sdk.PluginProvider, sdk.ModuleProvider, and sdk.StepProvider.
type turnIOPlugin struct{}

// NewTurnIOPlugin returns a new turnIOPlugin instance.
func NewTurnIOPlugin() sdk.PluginProvider {
	return &turnIOPlugin{}
}

// Manifest returns plugin metadata.
func (p *turnIOPlugin) Manifest() sdk.PluginManifest {
	return sdk.PluginManifest{
		Name:        "workflow-plugin-turnio",
		Version:     "0.1.0",
		Author:      "GoCodeAlone",
		Description: "turn.io WhatsApp API integration plugin",
	}
}

// ModuleTypes returns the module type names this plugin provides.
func (p *turnIOPlugin) ModuleTypes() []string {
	return []string{"turnio.provider"}
}

// CreateModule creates a module instance of the given type.
func (p *turnIOPlugin) CreateModule(typeName, name string, config map[string]any) (sdk.ModuleInstance, error) {
	switch typeName {
	case "turnio.provider":
		m, err := newTurnIOModule(name, config)
		if err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, fmt.Errorf("turnio plugin: unknown module type %q", typeName)
	}
}

// StepTypes returns all step type names this plugin provides.
func (p *turnIOPlugin) StepTypes() []string {
	return []string{
		"step.turnio_send_text",
		"step.turnio_send_media",
		"step.turnio_send_template",
		"step.turnio_send_interactive",
		"step.turnio_send_location",
		"step.turnio_list_messages",
		"step.turnio_check_contact",
		"step.turnio_upload_contacts",
		"step.turnio_update_profile",
		"step.turnio_upload_media",
		"step.turnio_get_media",
		"step.turnio_delete_media",
		"step.turnio_create_template",
		"step.turnio_list_templates",
		"step.turnio_fetch_template",
		"step.turnio_update_template",
		"step.turnio_delete_template",
		"step.turnio_configure_webhook",
		"step.turnio_create_flow",
		"step.turnio_list_flows",
		"step.turnio_send_flow",
		"step.turnio_list_journeys",
		"step.turnio_trigger_journey",
		"step.turnio_get_context",
		"step.turnio_set_context",
	}
}

// CreateStep creates a step instance of the given type.
func (p *turnIOPlugin) CreateStep(typeName, name string, config map[string]any) (sdk.StepInstance, error) {
	switch typeName {
	// Messages
	case "step.turnio_send_text":
		return newSendTextStep(name, config)
	case "step.turnio_send_media":
		return newSendMediaStep(name, config)
	case "step.turnio_send_template":
		return newSendTemplateStep(name, config)
	case "step.turnio_send_interactive":
		return newSendInteractiveStep(name, config)
	case "step.turnio_send_location":
		return newSendLocationStep(name, config)
	case "step.turnio_list_messages":
		return newListMessagesStep(name, config)
	// Contacts
	case "step.turnio_check_contact":
		return newCheckContactStep(name, config)
	case "step.turnio_upload_contacts":
		return newUploadContactsStep(name, config)
	case "step.turnio_update_profile":
		return newUpdateProfileStep(name, config)
	// Media
	case "step.turnio_upload_media":
		return newUploadMediaStep(name, config)
	case "step.turnio_get_media":
		return newGetMediaStep(name, config)
	case "step.turnio_delete_media":
		return newDeleteMediaStep(name, config)
	// Templates
	case "step.turnio_create_template":
		return newCreateTemplateStep(name, config)
	case "step.turnio_list_templates":
		return newListTemplatesStep(name, config)
	case "step.turnio_fetch_template":
		return newFetchTemplateStep(name, config)
	case "step.turnio_update_template":
		return newUpdateTemplateStep(name, config)
	case "step.turnio_delete_template":
		return newDeleteTemplateStep(name, config)
	// Webhooks
	case "step.turnio_configure_webhook":
		return newConfigureWebhookStep(name, config)
	// Flows
	case "step.turnio_create_flow":
		return newCreateFlowStep(name, config)
	case "step.turnio_list_flows":
		return newListFlowsStep(name, config)
	case "step.turnio_send_flow":
		return newSendFlowStep(name, config)
	// Journeys
	case "step.turnio_list_journeys":
		return newListJourneysStep(name, config)
	case "step.turnio_trigger_journey":
		return newTriggerJourneyStep(name, config)
	// Context
	case "step.turnio_get_context":
		return newGetContextStep(name, config)
	case "step.turnio_set_context":
		return newSetContextStep(name, config)
	default:
		return nil, fmt.Errorf("turnio plugin: unknown step type %q", typeName)
	}
}
