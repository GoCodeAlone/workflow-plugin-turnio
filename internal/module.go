package internal

import (
	"context"
	"fmt"
)

// turnIOModule creates the TurnClient and registers it.
type turnIOModule struct {
	name     string
	config   map[string]any
	apiToken string
	baseURL  string
}

func newTurnIOModule(name string, config map[string]any) (*turnIOModule, error) {
	apiToken, _ := config["apiToken"].(string)
	if apiToken == "" {
		return nil, fmt.Errorf("turnio.provider %q: config.apiToken is required", name)
	}
	baseURL, _ := config["baseUrl"].(string)
	return &turnIOModule{
		name:     name,
		config:   config,
		apiToken: apiToken,
		baseURL:  baseURL,
	}, nil
}

// Init creates the TurnClient and registers it in the global registry.
func (m *turnIOModule) Init() error {
	client := NewTurnClient(m.apiToken, m.baseURL)
	RegisterClient(m.name, client)
	return nil
}

// Start is a no-op for this module.
func (m *turnIOModule) Start(_ context.Context) error { return nil }

// Stop unregisters the client.
func (m *turnIOModule) Stop(_ context.Context) error {
	UnregisterClient(m.name)
	return nil
}
