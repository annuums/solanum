package solanum_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	solanum "github.com/annuums/solanum"
)

// FooService defines an interface for demonstration of DI in tests.
type FooService interface {
	Foo() string
}

// impl is a concrete implementation of FooService for testing.
type impl struct{}

// Foo returns a static string for testing dependency injection.
func (impl) Foo() string {
	return "bar"
}

func init() {
	// Set Gin to test mode to suppress logging
	gin.SetMode(gin.TestMode)
}

// TestNewModuleImplementsModule ensures NewModule returns a Module.
func TestNewModuleImplementsModule(t *testing.T) {
	var _ solanum.Module = solanum.NewModule("/test")
}

// TestMiddlewareChains validates pre- and post-middleware management.
func TestMiddlewareChains(t *testing.T) {
	m := solanum.NewModule("/")

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

// TestControllersAndDependencies ensures controllers and dependencies can be registered.
func TestControllersAndDependencies(t *testing.T) {
	m := solanum.NewModule("/")

	// Controller registration
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{Uri: "/a", Method: http.MethodGet, Handler: func(c *gin.Context) {}})
	m.SetControllers(ctrl)
	assert.Len(t, m.Controllers(), 1)

	// Dependency configuration
	dc := solanum.Dep[string]("key")
	m.SetDependencies(dc)
	assert.Len(t, m.Dependencies(), 1)
}

// TestSetRoutesWithoutDependencies ensures routes work with no DI configured.
func TestSetRoutesWithoutDependencies(t *testing.T) {
	r := gin.New()
	// Standalone route to verify Gin works
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	m := solanum.NewModule("/api")
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{Uri: "/ping", Method: "GET", Handler: func(c *gin.Context) { c.String(http.StatusOK, "ok") }})
	m.SetControllers(ctrl)
	// Register module routes without dependencies
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/ping", nil)
	r.ServeHTTP(rec, req)
	assert.Equal(t, "ok", rec.Body.String())
}

// TestSetRoutesWithDependencies verifies DI middleware injects FooService.
func TestSetRoutesWithDependencies(t *testing.T) {
	// Register FooService provider
	solanum.Register("foo", func() impl { return impl{} },
		solanum.WithSingleton(),
		solanum.As((*FooService)(nil)),
	)

	r := gin.New()
	m := solanum.NewModule("/api")
	m.SetDependencies(solanum.Dep[FooService]("foo"))

	// Handler uses injected service
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:    "/dep",
		Method: "GET",
		Handler: func(c *gin.Context) {
			d := solanum.GetDependency[FooService](c, "foo")
			c.String(http.StatusOK, d.Foo())
		},
	})
	m.SetControllers(ctrl)
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/dep", nil)
	r.ServeHTTP(rec, req)
	assert.Equal(t, "bar", rec.Body.String())
}
