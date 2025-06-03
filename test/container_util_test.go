package solanum_test

import (
	"github.com/annuums/solanum/container"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Disable Gin logger for tests
	gin.SetMode(gin.TestMode)
}

// TestDep ensures that Dep constructs the correct DependencyConfig for type T.
func TestDep(t *testing.T) {
	dc := container.Dep[int]("intKey")
	assert.Equal(t, "intKey", dc.Key)
	assert.Equal(t, reflect.TypeOf(0), dc.Type)
}

// TestGetDependency retrieves an injected value or returns zero value if absent.
func TestGetDependency(t *testing.T) {
	// Create test context
	tc := &gin.Context{}

	// Manually set a dependency in context
	tc.Set(container.DependencyPrefix+"foo", "bar")
	val := container.DepFromContext[string](tc, "foo")
	assert.Equal(t, "bar", val)

	// Missing key returns zero of type
	zero := container.DepFromContext[bool](tc, "missing")
	assert.False(t, zero)
}
