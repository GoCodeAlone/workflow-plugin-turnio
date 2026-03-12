package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	sdk "github.com/GoCodeAlone/workflow/plugin/external/sdk"
)

// uploadMediaStep implements step.turnio_upload_media.
type uploadMediaStep struct {
	name       string
	moduleName string
}

func newUploadMediaStep(name string, config map[string]any) (*uploadMediaStep, error) {
	return &uploadMediaStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *uploadMediaStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	mediaURL := resolveValue("media_url", current, config)
	if mediaURL != "" {
		// Upload via URL reference
		payload := map[string]any{"url": mediaURL}
		var result json.RawMessage
		if err := client.DoInto(ctx, "POST", "/v1/media", payload, &result); err != nil {
			return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
		}
		return &sdk.StepResult{Output: map[string]any{"result": string(result)}}, nil
	}

	// Upload via multipart form
	fileData := resolveValue("file_data", current, config)
	fileName := resolveValue("file_name", current, config)
	mimeType := resolveValue("mime_type", current, config)
	if fileData == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "media_url or file_data is required"}}, nil
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	if _, err := io.WriteString(fw, fileData); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", client.baseURL+"/v1/media", &buf)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	req.Header.Set("Authorization", "Bearer "+client.apiToken)
	if mimeType != "" {
		req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", w.Boundary()))
	} else {
		req.Header.Set("Content-Type", w.FormDataContentType())
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return &sdk.StepResult{Output: map[string]any{"error": fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body))}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"result": string(body)}}, nil
}

// getMediaStep implements step.turnio_get_media.
type getMediaStep struct {
	name       string
	moduleName string
}

func newGetMediaStep(name string, config map[string]any) (*getMediaStep, error) {
	return &getMediaStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *getMediaStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	mediaID := resolveValue("media_id", current, config)
	if mediaID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "media_id is required"}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "GET", fmt.Sprintf("/v1/media/%s", mediaID), nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"media": string(result)}}, nil
}

// deleteMediaStep implements step.turnio_delete_media.
type deleteMediaStep struct {
	name       string
	moduleName string
}

func newDeleteMediaStep(name string, config map[string]any) (*deleteMediaStep, error) {
	return &deleteMediaStep{name: name, moduleName: getModuleName(config)}, nil
}

func (s *deleteMediaStep) Execute(ctx context.Context, _ map[string]any, _ map[string]map[string]any, current map[string]any, _ map[string]any, config map[string]any) (*sdk.StepResult, error) {
	client, ok := GetClient(s.moduleName)
	if !ok {
		return &sdk.StepResult{Output: map[string]any{"error": "turnio client not found: " + s.moduleName}}, nil
	}

	mediaID := resolveValue("media_id", current, config)
	if mediaID == "" {
		return &sdk.StepResult{Output: map[string]any{"error": "media_id is required"}}, nil
	}

	var result json.RawMessage
	if err := client.DoInto(ctx, "DELETE", fmt.Sprintf("/v1/media/%s", mediaID), nil, &result); err != nil {
		return &sdk.StepResult{Output: map[string]any{"error": err.Error()}}, nil
	}
	return &sdk.StepResult{Output: map[string]any{"deleted": true, "media_id": mediaID}}, nil
}
