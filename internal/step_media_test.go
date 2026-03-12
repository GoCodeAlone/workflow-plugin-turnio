package internal

import (
	"context"
	"net/http"
	"testing"
)

func TestGetMediaStep_MissingID(t *testing.T) {
	step, _ := newGetMediaStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing media_id")
	}
}

func TestGetMediaStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"abc","url":"https://example.com/media/abc"}`))
	})
	defer cleanup()

	step, _ := newGetMediaStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"media_id": "abc123",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestDeleteMediaStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer cleanup()

	step, _ := newDeleteMediaStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"media_id": "abc123",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["deleted"] != true {
		t.Error("expected deleted=true")
	}
}
