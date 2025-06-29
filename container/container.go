package container

import (
	"fmt"
	"reflect"
	"sync"
)

type (

	// DependencyConfig defines a key and Go type for a dependency to be injected into handlers.
	DependencyConfig struct {
		// Key identifier used when registering and retrieving the dependency
		// registration key for the dependency
		Key string

		// Type reflect.Type of the interface or concrete type to resolve
		// expected reflect.Type for resolution
		Type reflect.Type
	}

	// providerEntry represents a registration record for a service provider.
	// It stores the factory function, lifecycle (singleton or transient),
	// any initialized instance, initialization hook, and type metadata.
	providerEntry struct {
		// factory constructs a new instance of the provider.
		factory func(...interface{}) interface{}

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

		// deps holds the dependencies of the provider, if any.
		deps []DependencyConfig
	}

	// container is the internal DI container managing provider registrations
	// and type-to-key mappings. It is safe for concurrent use.
	container struct {
		mu           sync.RWMutex              // protects all maps below
		providers    map[string]*providerEntry // key -> providerEntry
		interfaceMap map[reflect.Type]string   // interface type -> key
		typeMap      map[reflect.Type][]string // concrete type -> key
	}
)

// globalContainer is the shared, package-level DI container instance.
var globalContainer = &container{
	providers:    make(map[string]*providerEntry),
	interfaceMap: make(map[reflect.Type]string),
	typeMap:      make(map[reflect.Type][]string),
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

// WithDep lets you specify a key and type for the dependency.
// declares dep, before calling your provider func,
// the container should Resolve(key) and inject it as the T-typed argument.
func WithDep[T any](key string) RegisterOption {

	return func(pe *providerEntry) {

		var ptr *T
		pe.deps = append(pe.deps, DependencyConfig{
			Key:  key,
			Type: reflect.TypeOf(ptr).Elem(),
		})
	}
}

// As binds this provider under a specified interface type, allowing ResolveByType to work.
// Example: solanum.As((*MyInterface)(nil)) binds to MyInterface.
func As(ifacePtr interface{}) RegisterOption {

	// reflect.TypeOf((*MyIface)(nil)).Elem()
	t := reflect.TypeOf(ifacePtr).Elem()
	return func(pe *providerEntry) {

		pe.interfaceType = t
	}
}

// Register adds a new provider under the given key. The provider can be either:
//   - a factory function: func(...) (T, error)
//   - a concrete value: T
//
// Options control scope, init hooks, and interface binding.
func Register(key string, provider interface{}, opts ...RegisterOption) {

	// Default to singleton scope
	pe := &providerEntry{singleton: true}

	for _, opt := range opts {
		opt(pe)
	}

	pv := reflect.ValueOf(provider)
	pt := pv.Type()

	// inference dependency
	deps := pe.deps

	if pt.Kind() == reflect.Func && len(deps) == 0 {

		for i := 0; i < pt.NumIn(); i++ {

			deps = append(deps, DependencyConfig{
				Key:  "", // deps are not yet registered
				Type: pt.In(i),
			})
		}
	}

	// If provider is a function, wrap it to perform nested dependency resolution
	if pt.Kind() == reflect.Func {

		pe.factory = func(...interface{}) interface{} {

			// Resolve each function parameter by type
			args := make([]reflect.Value, len(deps))

			for i, d := range deps {

				var inst interface{}
				var err error

				// pick ResolveByType if interface, else Resolve by key
				if d.Key != "" {

					if d.Type.Kind() == reflect.Interface {

						inst, err = ResolveByType(d.Key, d.Type)
					} else {

						inst, err = Resolve(d.Key)
					}
				} else {

					// inference dependency from type only
					inst, err = globalContainer.resolveByReflectType(d.Type)
				}

				if err != nil {

					panic(fmt.Errorf("failed to resolve %q [%v]: %w", d.Key, d.Type, err))
				}

				args[i] = reflect.ValueOf(inst)
			}

			// call original provider with resolved args
			out := pv.Call(args)
			if len(out) == 2 && !out[1].IsNil() {

				panic(out[1].Interface())
			}

			return out[0].Interface()
		}
		// Provider returns the first result type
		pe.providerType = pt.Out(0)
	} else {

		// Static instance provider
		pe.factory = func(_ ...interface{}) interface{} {
			return provider
		}
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

	// Store the key under the concrete type
	globalContainer.typeMap[pe.providerType] = append(globalContainer.typeMap[pe.providerType], key)
	globalContainer.mu.Unlock()
}

// resolveByReflectType finds a registration key by interface or concrete type,
// then calls Resolve(key). Returns an error if no matching provider found.
func (c *container) resolveByReflectType(t reflect.Type) (interface{}, error) {

	c.mu.RLock()

	// If the type is an interface, check the interfaceMap first
	if key, ok := c.interfaceMap[t]; ok {

		c.mu.RUnlock()
		return Resolve(key)
	}

	// If the type is concrete, check the typeMap
	if keys, ok := c.typeMap[t]; ok {

		if len(keys) == 1 {

			c.mu.RUnlock()
			return Resolve(keys[0])
		}

		// If multiple keys exist for this type, return an error
		// e.g., ResolveByType("key", t)
		c.mu.RUnlock()
		return nil, fmt.Errorf("ambiguous providers for type %v: keys=%v", t, keys)
	}

	c.mu.RUnlock()
	return nil, fmt.Errorf("no provider for type %v", t)
}

// Resolve retrieves an instance registered under the given key.
//   - If the provider is a singleton and has already been constructed, it returns the stored instance.
//   - Otherwise, it invokes the factory (outside any locks), stores the instance if singleton,
//     and calls initHook exactly once (also outside the read-lock).
func Resolve(key string) (interface{}, error) {

	globalContainer.mu.RLock()
	pe, exists := globalContainer.providers[key]
	if !exists {

		globalContainer.mu.RUnlock()
		return nil, fmt.Errorf("no provider registered for key %q", key)
	}

	isSingleton := pe.singleton
	existing := pe.instance
	globalContainer.mu.RUnlock()

	// Return existing instance if already created.
	if isSingleton && existing != nil {

		return existing, nil
	}

	// Build the instance outside of any locks to avoid deadlocks.
	inst := pe.factory()
	if isSingleton {

		globalContainer.mu.Lock()
		if pe.instance == nil {

			pe.instance = inst
		}

		globalContainer.mu.Unlock()
	}

	var doHook bool
	if pe.initHook != nil {

		globalContainer.mu.Lock()
		if !pe.hookCalled {

			pe.hookCalled = true
			doHook = true
		}

		globalContainer.mu.Unlock()

		// Call the hook without holding any locks.
		if doHook {

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

	instType := reflect.TypeOf(inst)

	switch ifaceType.Kind() {
	case reflect.Interface:

		if !instType.Implements(ifaceType) {

			return nil, fmt.Errorf("provider %q does not implement %v", key, ifaceType)
		}

	default:

		if !instType.AssignableTo(ifaceType) {

			return nil, fmt.Errorf("provider %q: type %v not assignable to %v", key, instType, ifaceType)
		}
	}

	return inst, nil
}
