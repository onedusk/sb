package converter

import (
	"fmt"
	"sort"
	"sync"
)

var (
	registry = &Registry{
		converters: make(map[string]Converter),
	}
)

// Registry manages available converters
type Registry struct {
	mu         sync.RWMutex
	converters map[string]Converter
}

// Register adds a converter to the registry
func Register(conv Converter) error {
	return registry.Register(conv)
}

// Get retrieves a converter by name
func Get(name string) (Converter, error) {
	return registry.Get(name)
}

// List returns all registered converter names
func List() []string {
	return registry.List()
}

// ListConverters returns all registered converters
func ListConverters() []Converter {
	return registry.ListConverters()
}

// Register adds a converter to this registry
func (r *Registry) Register(conv Converter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := conv.Name()
	if name == "" {
		return fmt.Errorf("converter name cannot be empty")
	}

	if _, exists := r.converters[name]; exists {
		return fmt.Errorf("converter %q already registered", name)
	}

	r.converters[name] = conv
	return nil
}

// Get retrieves a converter by name
func (r *Registry) Get(name string) (Converter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conv, exists := r.converters[name]
	if !exists {
		return nil, fmt.Errorf("converter %q not found", name)
	}

	return conv, nil
}

// List returns all registered converter names sorted alphabetically
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.converters))
	for name := range r.converters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// ListConverters returns all registered converters
func (r *Registry) ListConverters() []Converter {
	r.mu.RLock()
	defer r.mu.RUnlock()

	convs := make([]Converter, 0, len(r.converters))
	for _, conv := range r.converters {
		convs = append(convs, conv)
	}
	return convs
}
