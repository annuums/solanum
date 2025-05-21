package solanum_test

import (
	"github.com/annuums/solanum/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	solanum "github.com/annuums/solanum"
)

// TestNewSolanumSingleton ensures NewSolanum returns the same Runner instance.
func TestNewSolanumSingleton(t *testing.T) {
	first := solanum.NewSolanum(1234)
	second := solanum.NewSolanum(5678)
	assert.Equal(t, first, second)
}

// TestGinEngineAccess verifies that GinEngine returns a non-nil *gin.Engine.
func TestGinEngineAccess(t *testing.T) {
	app := *solanum.NewSolanum(5050)
	eng := app.GinEngine()
	assert.NotNil(t, eng)
}

// TestCorsIntegration checks that calling Cors does not panic.
func TestCorsIntegration(t *testing.T) {
	// Set Gin mode to test to avoid logs
	gin.SetMode(gin.TestMode)
	app := *solanum.NewSolanum(5050)
	// Should not panic
	assert.NotPanics(t, func() {
		app.Cors(
			util.WithUrls([]string{"*"}),
			util.WithMethods([]string{"GET"}),
		)
	})
}
