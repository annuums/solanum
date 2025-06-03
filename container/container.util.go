package container

import (
	"context"
	"github.com/gin-gonic/gin"
	"reflect"
)

type ContextKey struct {
	solanumDepKey string
}

func NewContextKey(key string) ContextKey {
	return ContextKey{solanumDepKey: key}
}

// GetDependency retrieves a previously injected dependency from the Gin context.
// It looks up the value under the key composed of DependencyPrefix + key.
// If the value exists, it casts it to the requested generic type T and returns it.
// Otherwise, it returns the zero value for type T.
// Deprecated: This function is deprecated and will be removed in future versions. Use DepFromContext instead.
func GetDependency[T any](c *gin.Context, key string) T {
	if v, ok := c.Get(key); ok {
		return v.(T)
	}
	var zero T
	return zero
}

// DepFromContext retrieves a previously injected dependency from the Gin context.
// It looks up the value under the key composed of DependencyPrefix + key.
// If the value exists, it casts it to the requested generic type T and returns it.
// Otherwise, it returns the zero value for type T.
func DepFromContext[T any](ctx context.Context, key string) T {

	val := ctx.Value(NewContextKey(key))

	if inst, ok := val.(T); ok {

		return inst
	}

	var zero T
	return zero
}

// DepFromGinContext retrieves a previously injected dependency from the Gin context.
func DepFromGinContext[T any](c *gin.Context, key string) T {
	return DepFromContext[T](c.Request.Context(), key)
}

// Dep retrieves a previously injected dependency from the global container.
func Dep[T any](key string) T {

	config := DepConfig[T](key)

	// Use ResolveByType to enforce that the instance implements or is assignable to T
	inst, err := ResolveByType(key, config.Type)
	if err != nil {
		panic(err)
	}

	// Perform a type assertion to T
	if v, ok := inst.(T); ok {
		return v
	}

	// If we reach here, the instance was found but does not match T
	panic("dependency type mismatch: could not cast resolved instance to " + config.Type.String())
}

// DepConfig creates a DependencyConfig for type T against the specified key.
// It uses a pointer-to-T to obtain the reflect.Type of T and returns
// a DependencyConfig that can be passed to Module.SetDependencies.
func DepConfig[T any](key string) *DependencyConfig {
	// Create a nil pointer of type *T so that reflect.TypeOf(ptr).Elem()
	// yields the reflect.Type representing T.
	var ptr *T
	return &DependencyConfig{
		Key:  key,
		Type: reflect.TypeOf(ptr).Elem(),
	}
}
