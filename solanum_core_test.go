package solanum_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
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

// TestWithLoggerReceivesFrameworkLogs verifies framework logs go to the configured logger.
func TestWithLoggerReceivesFrameworkLogs(t *testing.T) {

	var buf bytes.Buffer
	server := solanum.NewSolanum(
		solanum.WithLogger(slog.New(slog.NewTextHandler(&buf, nil))),
	)
	server.Cors()
	assert.Contains(t, buf.String(), "cors allows all origins")

	assert.NotPanics(t, func() {
		solanum.NewSolanum(solanum.WithLogger(nil)).Cors()
	})
}

// TestInitModulesIsIdempotent verifies repeated initialization does not re-register routes.
func TestInitModulesIsIdempotent(t *testing.T) {

	server := newPingServer()
	assert.NotPanics(t, func() {
		server.InitModules()
		server.InitModules()
		server.Handler()
	})
}

// TestHandlerServesRoutes verifies Handler initializes routes and serves requests.
func TestHandlerServesRoutes(t *testing.T) {

	server := newPingServer()
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/ping", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}

// TestConfigurationAfterInitPanics verifies late configuration fails loudly instead of being ignored.
func TestConfigurationAfterInitPanics(t *testing.T) {

	server := newPingServer()
	server.InitModules()

	assert.Panics(t, func() { server.Cors() })
	assert.Panics(t, func() { server.SetModules(solanum.NewModule()) })
	assert.Panics(t, func() { server.AddModules(solanum.NewModule()) })
}

// freePort returns a TCP port that is currently free on localhost.
func freePort(t *testing.T) int {

	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// TestRunWithoutPortReturnsError verifies Run reports a missing port instead of panicking.
func TestRunWithoutPortReturnsError(t *testing.T) {

	server := solanum.NewSolanum()
	assert.ErrorIs(t, server.RunContext(context.Background()), solanum.ErrPortNotConfigured)
}

// TestRunContextGracefulShutdown serves requests and returns nil once the context is cancelled.
func TestRunContextGracefulShutdown(t *testing.T) {

	port := freePort(t)
	m := solanum.NewModule(solanum.WithUri("/api"))
	ctrl := solanum.NewController()
	ctrl.SetHandlers(&solanum.SolaService{
		Uri:     "/ping",
		Method:  http.MethodGet,
		Handler: func(c *gin.Context) { c.String(http.StatusOK, "pong") },
	})
	m.AddControllers(ctrl)
	server := solanum.NewSolanum(solanum.WithPort(port), solanum.WithLogger(nil))
	server.SetModules(m)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- server.RunContext(ctx) }()

	url := fmt.Sprintf("http://127.0.0.1:%d/api/ping", port)
	assert.Eventually(t, func() bool {
		res, err := http.Get(url)
		if err != nil {

			return false
		}

		defer res.Body.Close()
		return res.StatusCode == http.StatusOK
	}, 2*time.Second, 20*time.Millisecond)

	cancel()
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("server did not shut down")
	}
}

// TestRunContextPortInUseReturnsError verifies listen failures are returned as errors.
func TestRunContextPortInUseReturnsError(t *testing.T) {

	l, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	defer l.Close()

	server := solanum.NewSolanum(
		solanum.WithPort(l.Addr().(*net.TCPAddr).Port),
		solanum.WithLogger(nil),
	)
	assert.Error(t, server.RunContext(context.Background()))
}

// TestSetModulesReplacesAndAddModulesAppends verifies Set/Add semantics match Module and Controller.
func TestSetModulesReplacesAndAddModulesAppends(t *testing.T) {

	server := solanum.NewSolanum()
	server.SetModules(solanum.NewModule(), solanum.NewModule())
	server.SetModules(solanum.NewModule())
	assert.Len(t, server.Modules(), 1)

	server.AddModules(solanum.NewModule(), solanum.NewModule())
	assert.Len(t, server.Modules(), 3)
}
