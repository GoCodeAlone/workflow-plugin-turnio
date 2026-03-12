package internal

import (
	"context"
	"net/http"
	"testing"
)

func TestCheckContactStep_MissingContacts(t *testing.T) {
	step, _ := newCheckContactStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing contacts")
	}
}

func TestCheckContactStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"contacts":[{"wa_id":"27123","status":"valid"}]}`))
	})
	defer cleanup()

	step, _ := newCheckContactStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"phone": "+27123456789",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestUploadContactsStep_MissingContacts(t *testing.T) {
	step, _ := newUploadContactsStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing contacts")
	}
}
