package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestClient(t *testing.T, handler http.HandlerFunc) (*TurnClient, func()) {
	t.Helper()
	srv := httptest.NewServer(handler)
	client := NewTurnClient("test-token", srv.URL)
	RegisterClient("turnio", client)
	return client, func() {
		UnregisterClient("turnio")
		srv.Close()
	}
}

func TestSendTextStep_MissingTo(t *testing.T) {
	step, _ := newSendTextStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing to")
	}
}

func TestSendTextStep_MissingBody(t *testing.T) {
	step, _ := newSendTextStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{"to": "+27123"}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing body")
	}
}

func TestSendTextStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("X-Ratelimit-Remaining", "99")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"messages": []map[string]any{{"id": "abc"}}})
	})
	defer cleanup()

	step, _ := newSendTextStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"to":   "+27123456789",
		"body": "Hello",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
	if res.Output["to"] != "+27123456789" {
		t.Errorf("expected to field, got %v", res.Output["to"])
	}
}

func TestSendTextStep_MissingClient(t *testing.T) {
	step, _ := newSendTextStep("s", map[string]any{"module": "nonexistent"})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{"to": "+27123", "body": "hi"}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing client")
	}
}

func TestSendTextStep_RateLimitTracking(t *testing.T) {
	client, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ratelimit-Remaining", "42")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer cleanup()

	step, _ := newSendTextStep("s", map[string]any{})
	_, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"to":   "+27123456789",
		"body": "Hello",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if client.RateLimitRemaining() != 42 {
		t.Errorf("expected rate limit 42, got %d", client.RateLimitRemaining())
	}
}

func TestSendMediaStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"messages":[{"id":"img1"}]}`))
	})
	defer cleanup()

	step, _ := newSendMediaStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"to":   "+27123456789",
		"link": "https://example.com/image.jpg",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestSendTemplateStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"messages":[{"id":"tmpl1"}]}`))
	})
	defer cleanup()

	step, _ := newSendTemplateStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"to":            "+27123456789",
		"template_name": "hello_world",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestSendLocationStep_MissingCoords(t *testing.T) {
	step, _ := newSendLocationStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"to": "+27123456789",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing longitude/latitude")
	}
}

func TestSendLocationStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	defer cleanup()

	step, _ := newSendLocationStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"to":        "+27123456789",
		"longitude": 18.4241,
		"latitude":  -33.9249,
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}
