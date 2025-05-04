package solanum_test

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	solanum "github.com/annuums/solanum"
)

func init() {
	// Disable Gin logger for tests
	gin.SetMode(gin.TestMode)
}

// TestDep ensures that Dep constructs the correct DependencyConfig for type T.
func TestDep(t *testing.T) {
	dc := solanum.Dep[int]("intKey")
	assert.Equal(t, "intKey", dc.Key)
	assert.Equal(t, reflect.TypeOf(0), dc.Type)
}

// TestGetDependency retrieves an injected value or returns zero value if absent.
func TestGetDependency(t *testing.T) {
	// Create test context
	tc := &gin.Context{}

	// Manually set a dependency in context
	tc.Set(solanum.DependencyPrefix+"foo", "bar")
	val := solanum.GetDependency[string](tc, "foo")
	assert.Equal(t, "bar", val)

	// Missing key returns zero of type
	zero := solanum.GetDependency[bool](tc, "missing")
	assert.False(t, zero)
}
