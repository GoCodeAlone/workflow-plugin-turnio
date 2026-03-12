package internal

import "sync"

var (
	clientMu       sync.RWMutex
	clientRegistry = make(map[string]*TurnClient)
)

// RegisterClient adds a TurnClient to the global registry under the given name.
func RegisterClient(name string, c *TurnClient) {
	clientMu.Lock()
	defer clientMu.Unlock()
	clientRegistry[name] = c
}

// GetClient looks up a TurnClient by name.
func GetClient(name string) (*TurnClient, bool) {
	clientMu.RLock()
	defer clientMu.RUnlock()
	c, ok := clientRegistry[name]
	return c, ok
}

// UnregisterClient removes a TurnClient from the registry.
func UnregisterClient(name string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	delete(clientRegistry, name)
}
