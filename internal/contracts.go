package internal

import (
	turniov1 "github.com/GoCodeAlone/workflow-plugin-turnio/gen"
	pb "github.com/GoCodeAlone/workflow/plugin/external/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/structpb"
)

// ContractRegistry returns the typed contract descriptors for the turnio.provider
// module and all 25 step types. The workflow engine calls this via the
// sdk.ContractProvider interface to resolve proto message types for strict validation.
func (p *turnIOPlugin) ContractRegistry() *pb.ContractRegistry {
	return turnioContractRegistry
}

// turnioContractRegistry declares STRICT_PROTO contracts for the turnio module
// and all step types. The FileDescriptorSet includes google.protobuf.Struct
// (used in Input/Config fields) so the engine can resolve all message types.
var turnioContractRegistry = &pb.ContractRegistry{
	FileDescriptorSet: &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			protodesc.ToFileDescriptorProto(structpb.File_google_protobuf_struct_proto),
			protodesc.ToFileDescriptorProto(turniov1.File_turnio_proto),
		},
	},
	Contracts: []*pb.ContractDescriptor{
		// ── module ────────────────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_MODULE,
			ModuleType:    "turnio.provider",
			ConfigMessage: turnioProtoPkg + "ProviderConfig",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: messages ───────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_send_text",
			ConfigMessage: turnioProtoPkg + "SendTextConfig",
			InputMessage:  turnioProtoPkg + "SendTextInput",
			OutputMessage: turnioProtoPkg + "SendTextOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_send_media",
			ConfigMessage: turnioProtoPkg + "SendMediaConfig",
			InputMessage:  turnioProtoPkg + "SendMediaInput",
			OutputMessage: turnioProtoPkg + "SendMediaOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_send_template",
			ConfigMessage: turnioProtoPkg + "SendTemplateConfig",
			InputMessage:  turnioProtoPkg + "SendTemplateInput",
			OutputMessage: turnioProtoPkg + "SendTemplateOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_send_interactive",
			ConfigMessage: turnioProtoPkg + "SendInteractiveConfig",
			InputMessage:  turnioProtoPkg + "SendInteractiveInput",
			OutputMessage: turnioProtoPkg + "SendInteractiveOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_send_location",
			ConfigMessage: turnioProtoPkg + "SendLocationConfig",
			InputMessage:  turnioProtoPkg + "SendLocationInput",
			OutputMessage: turnioProtoPkg + "SendLocationOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_list_messages",
			ConfigMessage: turnioProtoPkg + "ListMessagesConfig",
			InputMessage:  turnioProtoPkg + "ListMessagesInput",
			OutputMessage: turnioProtoPkg + "ListMessagesOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: contacts ───────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_check_contact",
			ConfigMessage: turnioProtoPkg + "CheckContactConfig",
			InputMessage:  turnioProtoPkg + "CheckContactInput",
			OutputMessage: turnioProtoPkg + "CheckContactOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_upload_contacts",
			ConfigMessage: turnioProtoPkg + "UploadContactsConfig",
			InputMessage:  turnioProtoPkg + "UploadContactsInput",
			OutputMessage: turnioProtoPkg + "UploadContactsOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_update_profile",
			ConfigMessage: turnioProtoPkg + "UpdateProfileConfig",
			InputMessage:  turnioProtoPkg + "UpdateProfileInput",
			OutputMessage: turnioProtoPkg + "UpdateProfileOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: media ─────────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_upload_media",
			ConfigMessage: turnioProtoPkg + "UploadMediaConfig",
			InputMessage:  turnioProtoPkg + "UploadMediaInput",
			OutputMessage: turnioProtoPkg + "UploadMediaOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_get_media",
			ConfigMessage: turnioProtoPkg + "GetMediaConfig",
			InputMessage:  turnioProtoPkg + "GetMediaInput",
			OutputMessage: turnioProtoPkg + "GetMediaOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_delete_media",
			ConfigMessage: turnioProtoPkg + "DeleteMediaConfig",
			InputMessage:  turnioProtoPkg + "DeleteMediaInput",
			OutputMessage: turnioProtoPkg + "DeleteMediaOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: templates ─────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_create_template",
			ConfigMessage: turnioProtoPkg + "CreateTemplateConfig",
			InputMessage:  turnioProtoPkg + "CreateTemplateInput",
			OutputMessage: turnioProtoPkg + "CreateTemplateOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_list_templates",
			ConfigMessage: turnioProtoPkg + "ListTemplatesConfig",
			InputMessage:  turnioProtoPkg + "ListTemplatesInput",
			OutputMessage: turnioProtoPkg + "ListTemplatesOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_fetch_template",
			ConfigMessage: turnioProtoPkg + "FetchTemplateConfig",
			InputMessage:  turnioProtoPkg + "FetchTemplateInput",
			OutputMessage: turnioProtoPkg + "FetchTemplateOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_update_template",
			ConfigMessage: turnioProtoPkg + "UpdateTemplateConfig",
			InputMessage:  turnioProtoPkg + "UpdateTemplateInput",
			OutputMessage: turnioProtoPkg + "UpdateTemplateOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_delete_template",
			ConfigMessage: turnioProtoPkg + "DeleteTemplateConfig",
			InputMessage:  turnioProtoPkg + "DeleteTemplateInput",
			OutputMessage: turnioProtoPkg + "DeleteTemplateOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: webhooks ──────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_configure_webhook",
			ConfigMessage: turnioProtoPkg + "ConfigureWebhookConfig",
			InputMessage:  turnioProtoPkg + "ConfigureWebhookInput",
			OutputMessage: turnioProtoPkg + "ConfigureWebhookOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: flows ─────────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_create_flow",
			ConfigMessage: turnioProtoPkg + "CreateFlowConfig",
			InputMessage:  turnioProtoPkg + "CreateFlowInput",
			OutputMessage: turnioProtoPkg + "CreateFlowOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_list_flows",
			ConfigMessage: turnioProtoPkg + "ListFlowsConfig",
			InputMessage:  turnioProtoPkg + "ListFlowsInput",
			OutputMessage: turnioProtoPkg + "ListFlowsOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_send_flow",
			ConfigMessage: turnioProtoPkg + "SendFlowConfig",
			InputMessage:  turnioProtoPkg + "SendFlowInput",
			OutputMessage: turnioProtoPkg + "SendFlowOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: journeys ──────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_list_journeys",
			ConfigMessage: turnioProtoPkg + "ListJourneysConfig",
			InputMessage:  turnioProtoPkg + "ListJourneysInput",
			OutputMessage: turnioProtoPkg + "ListJourneysOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_trigger_journey",
			ConfigMessage: turnioProtoPkg + "TriggerJourneyConfig",
			InputMessage:  turnioProtoPkg + "TriggerJourneyInput",
			OutputMessage: turnioProtoPkg + "TriggerJourneyOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		// ── steps: context ───────────────────────────────────────────────────────
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_get_context",
			ConfigMessage: turnioProtoPkg + "GetContextConfig",
			InputMessage:  turnioProtoPkg + "GetContextInput",
			OutputMessage: turnioProtoPkg + "GetContextOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
		{
			Kind:          pb.ContractKind_CONTRACT_KIND_STEP,
			StepType:      "step.turnio_set_context",
			ConfigMessage: turnioProtoPkg + "SetContextConfig",
			InputMessage:  turnioProtoPkg + "SetContextInput",
			OutputMessage: turnioProtoPkg + "SetContextOutput",
			Mode:          pb.ContractMode_CONTRACT_MODE_STRICT_PROTO,
		},
	},
}

// turnioProtoPkg is the proto package prefix for all turnio typed messages.
const turnioProtoPkg = "workflow.plugin.turnio.v1."

// Compile-time assertion: turnIOPlugin implements sdk.ContractProvider.
var _ interface{ ContractRegistry() *pb.ContractRegistry } = (*turnIOPlugin)(nil)
