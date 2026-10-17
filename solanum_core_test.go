package solanum_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

// corsRequest sends a GET request with the given Origin to a server exposing GET /api/ping.
func corsRequest(t *testing.T, server solanum.Runner, origin string) *httptest.ResponseRecorder {

	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/ping", nil)
	req.Header.Set("Origin", origin)
	server.GinEngine().ServeHTTP(rec, req)
	return rec
}

// newPingServer builds a server exposing GET /api/ping.
func newPingServer() solanum.Runner {

	m := solanum.NewModule(solanum.WithUri("/api"))
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:     "/ping",
		Method:  http.MethodGet,
		Handler: func(c *gin.Context) { c.String(http.StatusOK, "pong") },
	})
	m.AddControllers(ctrl)

	server := solanum.NewSolanum(solanum.WithPort(5050))
	server.SetModules(m)
	return server
}

// TestCorsDefaultDoesNotReflectOrigin verifies the default CORS config answers with '*'
// instead of echoing back an arbitrary origin.
func TestCorsDefaultDoesNotReflectOrigin(t *testing.T) {

	server := newPingServer()
	server.Cors()
	server.InitModules()

	rec := corsRequest(t, server, "https://evil.example")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
}

// TestCorsCredentialsWithAllOriginsPanics rejects credentials without an explicit origin list.
func TestCorsCredentialsWithAllOriginsPanics(t *testing.T) {

	server := newPingServer()
	assert.Panics(t, func() {
		server.Cors(cors.WithAllowCredentials(true))
	})
}

// TestCorsCredentialsWithUrls allows only the configured origin with credentials.
func TestCorsCredentialsWithUrls(t *testing.T) {

	server := newPingServer()
	server.Cors(
		cors.WithUrls([]string{"https://ok.example"}),
		cors.WithAllowCredentials(true),
	)
	server.InitModules()

	rec := corsRequest(t, server, "https://ok.example")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "https://ok.example", rec.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))

	rec = corsRequest(t, server, "https://evil.example")
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

// TestCorsMaxAge verifies MaxAge is passed through as a duration on preflight responses.
func TestCorsMaxAge(t *testing.T) {

	server := newPingServer()
	server.Cors(cors.WithMaxAge(10 * time.Minute))
	server.InitModules()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/ping", nil)
	req.Header.Set("Origin", "https://a.example")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	server.GinEngine().ServeHTTP(rec, req)

	assert.Equal(t, "600", rec.Header().Get("Access-Control-Max-Age"))
}

// TestRecoveryEnabledByDefault verifies a panicking handler yields 500 instead of crashing.
func TestRecoveryEnabledByDefault(t *testing.T) {

	m := solanum.NewModule(solanum.WithUri("/api"))
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:     "/panic",
		Method:  http.MethodGet,
		Handler: func(c *gin.Context) { panic("boom") },
	})
	m.AddControllers(ctrl)

	server := solanum.NewSolanum(
		solanum.WithPort(5050),
		solanum.WithRecovery(gin.RecoveryWithWriter(io.Discard)),
	)
	server.SetModules(m)
	server.InitModules()

	rec := httptest.NewRecorder()
	assert.NotPanics(t, func() {
		server.GinEngine().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/panic", nil))
	})
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// TestGlobalMiddlewaresRunBeforeModules verifies WithMiddlewares runs ahead of module middlewares.
func TestGlobalMiddlewaresRunBeforeModules(t *testing.T) {

	var order []string
	m := solanum.NewModule(solanum.WithUri("/api"))
	m.AddPreMiddleware(func(c *gin.Context) { order = append(order, "module-pre") })
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:     "/ping",
		Method:  http.MethodGet,
		Handler: func(c *gin.Context) { order = append(order, "handler") },
	})
	m.AddControllers(ctrl)

	server := solanum.NewSolanum(
		solanum.WithPort(5050),
		solanum.WithMiddlewares(func(c *gin.Context) { order = append(order, "global") }),
	)
	server.SetModules(m)
	server.InitModules()
	server.GinEngine().ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/api/ping", nil))

	assert.Equal(t, []string{"global", "module-pre", "handler"}, order)
}

// TestWithGinMode verifies the gin mode option is applied.
func TestWithGinMode(t *testing.T) {

	defer gin.SetMode(gin.TestMode)
	solanum.NewSolanum(solanum.WithGinMode(gin.ReleaseMode))
	assert.Equal(t, gin.ReleaseMode, gin.Mode())
}
