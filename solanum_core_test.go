package solanum_test

import (
	"testing"

	"github.com/annuums/solanum/middleware/cors"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/annuums/solanum"
)

// TestNewSolanumReturnsNewRunner ensures each NewSolanum call returns a distinct Runner.
func TestNewSolanumReturnsNewRunner(t *testing.T) {
	first := solanum.NewSolanum(
		solanum.WithPort(1234),
	)
	second := solanum.NewSolanum(
		solanum.WithPort(5678),
	)
	assert.NotEqual(t, first, second)
	assert.Equal(t, 1234, first.Port())
	assert.Equal(t, 5678, second.Port())
}

// TestGinEngineAccess verifies that GinEngine returns a non-nil *gin.Engine.
func TestGinEngineAccess(t *testing.T) {
	server := solanum.NewSolanum(
		solanum.WithPort(5050),
	)
	eng := server.GinEngine()
	assert.NotNil(t, eng)
}

// TestCorsIntegration checks that calling Cors does not panic.
func TestCorsIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	server := solanum.NewSolanum(
		solanum.WithPort(5050),
	)
	assert.NotPanics(t, func() {
		server.Cors(
			cors.WithUrls([]string{"*"}),
			cors.WithMethods([]string{"GET"}),
		)
	})
}
