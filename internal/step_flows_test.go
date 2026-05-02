package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListFlowsStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/flows" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"flows":[]}`))
	})
	defer cleanup()

	step, _ := newListFlowsStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestListFlowsStep_MissingClient(t *testing.T) {
	step, _ := newListFlowsStep("s", map[string]any{"module": "nonexistent"})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing client")
	}
}

func TestCreateFlowStep_MissingName(t *testing.T) {
	step, _ := newCreateFlowStep("s", map[string]any{})
	RegisterClient("turnio", NewTurnClient("tok", "http://localhost"))
	defer UnregisterClient("turnio")
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for missing flow_name")
	}
}

func TestCreateFlowStep_Success(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"flow1","name":"my-flow"}`))
	})
	defer cleanup()

	step, _ := newCreateFlowStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"flow_name": "my-flow",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if errVal, ok := res.Output["error"]; ok {
		t.Fatalf("unexpected error: %v", errVal)
	}
}

func TestSendFlowStep_MissingTo(t *testing.T) {
	step, _ := newSendFlowStep("s", map[string]any{})
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

func TestSendFlowStep_HTTPError(t *testing.T) {
	_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid flow_id"}`))
	})
	defer cleanup()

	step, _ := newSendFlowStep("s", map[string]any{})
	res, err := step.Execute(context.Background(), nil, nil, map[string]any{
		"to":      "+27123456789",
		"flow_id": "bad-id",
	}, nil, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Output["error"] == nil {
		t.Fatal("expected error for 4xx response")
	}
}

func TestSendFlowStep_FlowMessageVersionDefaultAndOverride(t *testing.T) {
	tests := []struct {
		name    string
		current map[string]any
		config  map[string]any
		want    string
	}{
		{
			name:    "default",
			current: map[string]any{"to": "+27123456789", "flow_id": "flow-1"},
			config:  map[string]any{},
			want:    "3",
		},
		{
			name:    "current zero overrides default",
			current: map[string]any{"to": "+27123456789", "flow_id": "flow-1", "flow_message_version": 0},
			config:  map[string]any{"flow_message_version": 7},
			want:    "0",
		},
		{
			name:    "config override",
			current: map[string]any{"to": "+27123456789", "flow_id": "flow-1"},
			config:  map[string]any{"flow_message_version": 2},
			want:    "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cleanup := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.URL.Path != "/v1/messages" {
					t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("decode payload: %v", err)
				}
				got := payload["interactive"].(map[string]any)["action"].(map[string]any)["parameters"].(map[string]any)["flow_message_version"]
				if got != tt.want {
					t.Fatalf("flow_message_version = %v, want %q", got, tt.want)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"messages":[{"id":"msg1"}]}`))
			})
			defer cleanup()

			step, _ := newSendFlowStep("s", map[string]any{})
			res, err := step.Execute(context.Background(), nil, nil, tt.current, nil, tt.config)
			if err != nil {
				t.Fatal(err)
			}
			if errVal, ok := res.Output["error"]; ok {
				t.Fatalf("unexpected error: %v", errVal)
			}
		})
	}
}
