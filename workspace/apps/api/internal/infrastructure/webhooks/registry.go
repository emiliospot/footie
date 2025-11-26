package webhooks

import (
	"fmt"
	"strings"
)

// Registry manages provider instances and routes webhooks to the correct provider.
type Registry struct {
	providers map[string]Provider
}

// NewRegistry creates a new provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(provider Provider) {
	r.providers[strings.ToLower(provider.Name())] = provider
}

// GetProvider retrieves a provider by name (case-insensitive).
func (r *Registry) GetProvider(name string) (Provider, error) {
	provider, exists := r.providers[strings.ToLower(name)]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

// ListProviders returns all registered provider names.
func (r *Registry) ListProviders() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}
