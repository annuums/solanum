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

// TestWithUriValidation rejects module prefixes that do not start with '/'.
func TestWithUriValidation(t *testing.T) {

	assert.Panics(t, func() { solanum.NewModule(solanum.WithUri("users")) })
	assert.NotPanics(t, func() { solanum.NewModule(solanum.WithUri("")) })

	opts := []solanum.ModuleOption{solanum.WithUri("/users")}
	assert.Equal(t, "/users", solanum.NewModule(opts...).Uri())
}

// TestSetRoutesValidatesServices rejects services that would register broken routes.
func TestSetRoutesValidatesServices(t *testing.T) {

	cases := map[string]*solanum.SolaService{
		"nil service": nil,
		"no handler":  {Uri: "/x", Method: http.MethodGet},
		"no method":   {Uri: "/x", Handler: func(c *gin.Context) {}},
	}

	for name, svc := range cases {
		t.Run(name, func(t *testing.T) {

			m := solanum.NewModule(solanum.WithUri("/api"))
			ctrl := solanum.NewController()
			ctrl.SetHandlers(svc)
			m.AddControllers(ctrl)
			assert.Panics(t, func() { m.SetRoutes(gin.New().Group("/api")) })
		})
	}
}

// TestPostHandlersRunWhenAborted ensures post handlers still run on error paths.
func TestPostHandlersRunWhenAborted(t *testing.T) {

	var order []string
	var abortedSeen bool
	r := gin.New()

	m := solanum.NewModule(solanum.WithUri("/api"))
	m.AddPostMiddleware(func(c *gin.Context) {
		order = append(order, "module-post")
		abortedSeen = c.IsAborted()
	})

	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:    "/fail",
		Method: http.MethodGet,
		Handler: func(c *gin.Context) {
			order = append(order, "handler")
			c.AbortWithStatus(http.StatusInternalServerError)
		},
		PostHandlers: []gin.HandlerFunc{
			func(c *gin.Context) { order = append(order, "service-post") },
		},
	})
	m.AddControllers(ctrl)
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/fail", nil))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, []string{"handler", "service-post", "module-post"}, order)
	assert.True(t, abortedSeen)
}

// TestPostHandlersRunWhenPreAborts ensures post handlers run when a pre-middleware aborts.
func TestPostHandlersRunWhenPreAborts(t *testing.T) {

	var order []string
	r := gin.New()

	m := solanum.NewModule(solanum.WithUri("/api"))
	m.AddPreMiddleware(func(c *gin.Context) { c.AbortWithStatus(http.StatusUnauthorized) })
	m.AddPostMiddleware(func(c *gin.Context) { order = append(order, "module-post") })

	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:     "/secret",
		Method:  http.MethodGet,
		Handler: func(c *gin.Context) { order = append(order, "handler") },
	})
	m.AddControllers(ctrl)
	m.SetRoutes(r.Group("/api"))

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/secret", nil))

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, []string{"module-post"}, order)
}
