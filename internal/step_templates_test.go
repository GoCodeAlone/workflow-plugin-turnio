package internal

import (
	"context"
	"net/http"
	"testing"
)

func TestListTemplatesStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"templates":[]}`))
	})
	defer cleanup()

	step, _ := newListTemplatesStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestCreateTemplateStep_MissingName(t *testing.T) {
	step, _ := newCreateTemplateStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing template_name")
	}
}

func TestFetchTemplateStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"hello_world"}`))
	})
	defer cleanup()

	step, _ := newFetchTemplateStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"template_name": "hello_world",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}
