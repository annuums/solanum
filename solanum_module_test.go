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

	m.SetPreMiddlewares(func(c *gin.Context) {}, func(c *gin.Context) {})
	assert.Len(t, m.PreMiddlewares(), 2)

	m.AddPreMiddleware(func(c *gin.Context) {})
	assert.Len(t, m.PreMiddlewares(), 3)

	m.SetPostMiddlewares(func(c *gin.Context) {})
	assert.Len(t, m.PostMiddlewares(), 1)

	m.AddPostMiddleware(func(c *gin.Context) {})
	assert.Len(t, m.PostMiddlewares(), 2)
}

// TestServiceHandlerExecutionOrder asserts the full middleware execution order:
// module pre → service pre → handler → service post → module post.
func TestServiceHandlerExecutionOrder(t *testing.T) {

	var order []string
	r := gin.New()

	m := solanum.NewModule(solanum.WithUri("/api"))
	m.SetPreMiddlewares(func(c *gin.Context) { order = append(order, "module-pre") })
	m.SetPostMiddlewares(func(c *gin.Context) { order = append(order, "module-post") })

	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:    "/order",
		Method: http.MethodGet,
		PreHandlers: []gin.HandlerFunc{
			func(c *gin.Context) { order = append(order, "service-pre") },
		},
		Handler: func(c *gin.Context) {
			order = append(order, "handler")
			c.String(http.StatusOK, "ok")
		},
		PostHandlers: []gin.HandlerFunc{
			func(c *gin.Context) { order = append(order, "service-post") },
		},
	})
	m.AddControllers(ctrl)
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/order", nil)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, []string{"module-pre", "service-pre", "handler", "service-post", "module-post"}, order)
}

// TestServiceHandlerNilPrePostHandlers ensures a SolaService with no PreHandlers/PostHandlers
// behaves identically to the original single-handler registration.
func TestServiceHandlerNilPrePostHandlers(t *testing.T) {

	r := gin.New()

	m := solanum.NewModule(solanum.WithUri("/api"))
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:    "/nil",
		Method: http.MethodGet,
		Handler: func(c *gin.Context) {
			c.String(http.StatusOK, "ok")
		},
	})
	m.AddControllers(ctrl)
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/nil", nil)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

// TestSetRoutes ensures routes are correctly registered and respond to requests.
func TestSetRoutes(t *testing.T) {
	r := gin.New()

	m := solanum.NewModule(
		solanum.WithUri("/api"),
	)
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:     "/ping",
		Method:  "GET",
		Handler: func(c *gin.Context) { c.String(http.StatusOK, "ok") },
	})
	m.AddControllers(ctrl)
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/ping", nil)
	r.ServeHTTP(rec, req)
	assert.Equal(t, "ok", rec.Body.String())
}
