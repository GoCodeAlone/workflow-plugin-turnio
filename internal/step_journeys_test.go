package internal

import (
	"context"
	"net/http"
	"testing"
)

func TestListJourneysStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/journeys" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"journeys":[]}`))
	})
	defer cleanup()

	step, _ := newListJourneysStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestListJourneysStep_HTTPError(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"unauthorized"}`))
	})
	defer cleanup()

	step, _ := newListJourneysStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for 4xx response")
	}
}

func TestTriggerJourneyStep_MissingJourneyID(t *testing.T) {
	step, _ := newTriggerJourneyStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing journey_id")
	}
}

func TestTriggerJourneyStep_MissingPhone(t *testing.T) {
	step, _ := newTriggerJourneyStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"journey_id": "j1",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing phone")
	}
}

func TestTriggerJourneyStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"triggered"}`))
	})
	defer cleanup()

	step, _ := newTriggerJourneyStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"journey_id": "j1",
		"phone":      "+27123456789",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
	if res.Output["journey_id"] != "j1" {
		t.Errorf("expected journey_id in output")
	}
}
