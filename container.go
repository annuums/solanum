package solanum

import (
	"fmt"
	"reflect"
	"sync"
)

// providerEntry represents a registration record for a service provider.
// It stores the factory function, lifecycle (singleton or transient),
// any initialized instance, initialization hook, and type metadata.
type providerEntry struct {
	// factory constructs a new instance of the provider.
	factory func() interface{}

	// singleton indicates whether to reuse the same instance across resolves.
	singleton bool

	// instance holds the created singleton instance after first resolution.
	instance interface{}

	// hookCalled tracks if the initHook has already been executed.
	hookCalled bool

	// initHook is an optional callback that runs once after creating the instance.
	initHook func(interface{})

	// interfaceType, if non-nil, registers this provider under a Go interface type.
	interfaceType reflect.Type

	// providerType is the concrete type returned by the factory or provided directly.
	providerType reflect.Type
}

// container is the internal DI container managing provider registrations
// and type-to-key mappings. It is safe for concurrent use.
type container struct {
	mu           sync.RWMutex              // protects all maps below
	providers    map[string]*providerEntry // key -> providerEntry
	interfaceMap map[reflect.Type]string   // interface type -> key
	typeMap      map[reflect.Type]string   // concrete type -> key
}

// globalContainer is the shared, package-level DI container instance.
var globalContainer = &container{
	providers:    make(map[string]*providerEntry),
	interfaceMap: make(map[reflect.Type]string),
	typeMap:      make(map[reflect.Type]string),
}

// RegisterOption configures a providerEntry (e.g., scope, init hook, interface binding).
type RegisterOption func(*providerEntry)

// WithSingleton marks the provider as a singleton (default behavior).
func WithSingleton() RegisterOption {
	return func(pe *providerEntry) { pe.singleton = true }
}

// WithTransient marks the provider as transient, creating a new instance on each resolve.
func WithTransient() RegisterOption {
	return func(pe *providerEntry) { pe.singleton = false }
}

// WithInit sets a hook function to be invoked once after the instance is created.
func WithInit(hook func(interface{})) RegisterOption {
	return func(pe *providerEntry) { pe.initHook = hook }
}

// As binds this provider under a specified interface type, allowing ResolveByType to work.
// Example: solanum.As((*MyInterface)(nil)) binds to MyInterface.
func As(ifacePtr interface{}) RegisterOption {
	// reflect.TypeOf((*MyIface)(nil)).Elem()
	t := reflect.TypeOf(ifacePtr).Elem()
	return func(pe *providerEntry) { pe.interfaceType = t }
}

// Register adds a new provider under the given key. The provider can be either:
//   - a factory function: func(...) (T, error)
//   - a concrete value: T
//
// Options control scope, init hooks, and interface binding.
func Register(key string, provider interface{}, opts ...RegisterOption) {
	// Default to singleton scope
	pe := &providerEntry{singleton: true}

	pv := reflect.ValueOf(provider)
	pt := pv.Type()

	// If provider is a function, wrap it to perform nested dependency resolution
	if pt.Kind() == reflect.Func {
		pe.factory = func() interface{} {
			// Resolve each function parameter by type
			in := make([]reflect.Value, pt.NumIn())
			for i := 0; i < pt.NumIn(); i++ {
				argType := pt.In(i)
				dep, err := globalContainer.resolveByReflectType(argType)
				if err != nil {
					panic(fmt.Errorf("cannot resolve dependency %v: %w", argType, err))
				}
				in[i] = reflect.ValueOf(dep)
			}
			// Call factory and handle optional error return
			out := pv.Call(in)
			if len(out) == 2 && !out[1].IsNil() {
				panic(out[1].Interface())
			}
			return out[0].Interface()
		}
		// Provider returns the first result type
		pe.providerType = pt.Out(0)
	} else {
		// Static instance provider
		pe.factory = func() interface{} { return provider }
		pe.providerType = reflect.TypeOf(provider)
	}

	// Apply all registration options
	for _, o := range opts {
		o(pe)
	}

	// Store in global container with thread safety
	globalContainer.mu.Lock()
	globalContainer.providers[key] = pe
	if pe.interfaceType != nil {
		globalContainer.interfaceMap[pe.interfaceType] = key
	}
	globalContainer.typeMap[pe.providerType] = key
	globalContainer.mu.Unlock()
}

// resolveByReflectType finds a registration key by interface or concrete type,
// then calls Resolve(key). Returns an error if no matching provider found.
func (c *container) resolveByReflectType(t reflect.Type) (interface{}, error) {
	c.mu.RLock()
	if key, ok := c.interfaceMap[t]; ok {
		c.mu.RUnlock()
		return Resolve(key)
	}
	if key, ok := c.typeMap[t]; ok {
		c.mu.RUnlock()
		return Resolve(key)
	}
	c.mu.RUnlock()
	return nil, fmt.Errorf("no provider for type %v", t)
}

// Resolve retrieves an instance by key. For singletons, it reuses the instance.
// If an init hook is defined, it will be called exactly once.
func Resolve(key string) (interface{}, error) {
	globalContainer.mu.RLock()
	pe, ok := globalContainer.providers[key]
	if !ok {
		globalContainer.mu.RUnlock()
		return nil, fmt.Errorf("no provider registered for key %s", key)
	}
	// Return existing singleton if already initialized and hook called
	if pe.singleton && pe.instance != nil && pe.hookCalled {
		inst := pe.instance
		globalContainer.mu.RUnlock()
		return inst, nil
	}
	globalContainer.mu.RUnlock()

	// Create a new instance via factory
	inst := pe.factory()
	// Store singleton instance if applicable
	if pe.singleton {
		globalContainer.mu.Lock()
		if pe.instance == nil {
			pe.instance = inst
		}
		globalContainer.mu.Unlock()
	}

	// Invoke init hook once
	if pe.initHook != nil {
		globalContainer.mu.Lock()
		already := pe.hookCalled
		if !already {
			pe.hookCalled = true
		}
		globalContainer.mu.Unlock()
		if !already {
			pe.initHook(inst)
		}
	}
	return inst, nil
}

// ResolveByType resolves by key and additionally asserts that the instance
// implements the specified interface type. Passing nil for ifaceType skips the check.
func ResolveByType(key string, ifaceType reflect.Type) (interface{}, error) {
	inst, err := Resolve(key)
	if err != nil {
		return nil, err
	}
	if ifaceType == nil {
		return inst, nil
	}
	if !reflect.TypeOf(inst).Implements(ifaceType) {
		return nil, fmt.Errorf("provider %q does not implement %v", key, ifaceType)
	}
	return inst, nil
}
