package solanum_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	solanum "github.com/annuums/solanum"
)

func init() {
	// Set Gin to test mode to suppress logging
	gin.SetMode(gin.TestMode)
}

// TestNewModuleImplementsModule ensures NewModule returns a Module.
func TestNewModuleImplementsModule(t *testing.T) {
	var _ solanum.Module = solanum.NewModule(
		solanum.WithUri("/test"),
	)
}

// TestMiddlewareChains validates pre- and post-middleware management.
func TestMiddlewareChains(t *testing.T) {
	m := solanum.NewModule(
		solanum.WithUri("/"),
	)

	// Replace and count pre-middlewares
	m.SetPreMiddlewares(func(c *gin.Context) {}, func(c *gin.Context) {})
	assert.Len(t, m.PreMiddlewares(), 2)

	// Append another pre-middleware
	m.AddPreMiddleware(func(c *gin.Context) {})
	assert.Len(t, m.PreMiddlewares(), 3)

	// Replace and count post-middlewares
	m.SetPostMiddlewares(func(c *gin.Context) {})
	assert.Len(t, m.PostMiddlewares(), 1)

	// Append another post-middleware
	m.AddPostMiddleware(func(c *gin.Context) {})
	assert.Len(t, m.PostMiddlewares(), 2)
}

// TestSetRoutes ensures routes work with no DI configured.
func TestSetRoutes(t *testing.T) {
	r := gin.New()
	// Standalone route to verify Gin works
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	m := solanum.NewModule(
		solanum.WithUri("/api"),
	)
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{Uri: "/ping", Method: "GET", Handler: func(c *gin.Context) { c.String(http.StatusOK, "ok") }})
	m.AddControllers(ctrl)
	// Register module routes without dependencies
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/ping", nil)
	r.ServeHTTP(rec, req)
	assert.Equal(t, "ok", rec.Body.String())
}
