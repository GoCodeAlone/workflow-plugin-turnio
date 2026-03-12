package internal

import (
	"context"
	"testing"
)

func TestTurnIOModuleInit(t *testing.T) {
	m, err := newTurnIOModule("test", map[string]any{
		"apiToken": "test-token",
	})
	if err != nil {
		t.Fatalf("newTurnIOModule: %v", err)
	}
	if err := m.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	_, ok := GetClient("test")
	if !ok {
		t.Fatal("expected client to be registered")
	}
	if err := m.Stop(context.Background()); err != nil {
		t.Fatalf("Stop: %v", err)
	}
	_, ok = GetClient("test")
	if ok {
		t.Fatal("expected client to be unregistered after stop")
	}
}

func TestTurnIOModuleMissingToken(t *testing.T) {
	_, err := newTurnIOModule("test", map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing apiToken")
	}
}

func TestTurnIOModuleCustomBaseURL(t *testing.T) {
	m, err := newTurnIOModule("custom", map[string]any{
		"apiToken": "tok",
		"baseUrl":  "https://custom.example.com",
	})
	if err != nil {
		t.Fatalf("newTurnIOModule: %v", err)
	}
	if err := m.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	c, ok := GetClient("custom")
	if !ok {
		t.Fatal("expected client registered")
	}
	if c.baseURL != "https://custom.example.com" {
		t.Errorf("expected custom baseURL, got %q", c.baseURL)
	}
	_ = m.Stop(context.Background())
}
