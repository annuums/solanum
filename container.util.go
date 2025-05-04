package solanum

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

// GetDependency retrieves a previously injected dependency from the Gin context.
// It looks up the value under the key composed of DependencyPrefix + key.
// If the value exists, it casts it to the requested generic type T and returns it.
// Otherwise, it returns the zero value for type T.
func GetDependency[T any](c *gin.Context, key string) T {
	dependencyKey := DependencyPrefix + key
	if v, ok := c.Get(dependencyKey); ok {
		return v.(T)
	}
	var zero T
	return zero
}

// Dep creates a DependencyConfig for type T against the specified key.
// It uses a pointer-to-T to obtain the reflect.Type of T and returns
// a DependencyConfig that can be passed to Module.SetDependencies.
func Dep[T any](key string) DependencyConfig {
	// Create a nil pointer of type *T so that reflect.TypeOf(ptr).Elem()
	// yields the reflect.Type representing T.
	var ptr *T
	return DependencyConfig{
		Key:  key,
		Type: reflect.TypeOf(ptr).Elem(),
	}
}
