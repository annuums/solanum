package solanum

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cors2 "github.com/annuums/solanum/middleware/cors"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server implements the Runner interface and drives the application startup.
// It holds the Gin engine, listening port, and registered modules.
type Server struct {
	engine      *gin.Engine       // underlying Gin engine
	port        int               // TCP port to listen on
	modules     []Module          // registered modules
	ginMode     string            // gin mode applied before the engine is created
	recovery    gin.HandlerFunc   // panic recovery middleware (nil disables it)
	middlewares []gin.HandlerFunc // application-wide middlewares
	logger      *slog.Logger      // framework logger
	initOnce    sync.Once         // guards module route registration
	initialized bool              // true once module routes are registered

	readHeaderTimeout time.Duration // http.Server.ReadHeaderTimeout
	readTimeout       time.Duration // http.Server.ReadTimeout
	writeTimeout      time.Duration // http.Server.WriteTimeout
	idleTimeout       time.Duration // http.Server.IdleTimeout
	shutdownTimeout   time.Duration // grace period for in-flight requests on shutdown
}

// Default http.Server timeouts applied by NewSolanum.
// WriteTimeout is unlimited by default so streaming responses (SSE, downloads) keep working.
const (
	DefaultReadHeaderTimeout = 10 * time.Second
	DefaultReadTimeout       = 30 * time.Second
	DefaultWriteTimeout      = 0
	DefaultIdleTimeout       = 120 * time.Second
	DefaultShutdownTimeout   = 10 * time.Second
)

// ErrPortNotConfigured is returned by Run when no valid port is set.
var ErrPortNotConfigured = errors.New("solanum: server port is not configured; use WithPort")

// Run initializes all modules and serves HTTP on the configured port until
// SIGINT or SIGTERM is received, then shuts down gracefully.
func (server *Server) Run() error {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return server.RunContext(ctx)
}

// RunContext initializes all modules and serves HTTP on the configured port until ctx is done,
// then shuts down gracefully within the shutdown timeout.
// It returns nil after a graceful shutdown and an error when the server fails to start or stop.
func (server *Server) RunContext(ctx context.Context) error {
	if server.port <= 0 {

		return ErrPortNotConfigured
	}

	addr := fmt.Sprintf(":%d", server.port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: server.readHeaderTimeout,
		ReadTimeout:       server.readTimeout,
		WriteTimeout:      server.writeTimeout,
		IdleTimeout:       server.idleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {

		server.logger.Info("solanum is running", "addr", addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("solanum: failed to run server on %s: %w", addr, err)
	case <-ctx.Done():
	}

	server.logger.Info("solanum is shutting down", "addr", addr)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), server.shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {

		return fmt.Errorf("solanum: failed to shut down server on %s: %w", addr, err)
	}

	return nil
}

// InitModules sets up routing groups for each Module and applies their routes.
// It is idempotent: routes are registered only on the first call.
func (server *Server) InitModules() {

	server.initOnce.Do(func() {

		server.logger.Debug("initializing modules", "count", len(server.modules))
		for _, m := range server.modules {

			m.SetRoutes(
				server.GinEngine().Group(
					m.Uri(),
				),
			)
		}

		server.initialized = true
	})
}

// Handler initializes module routes (once) and returns the server as an http.Handler.
// Useful for tests with httptest or for embedding into a custom http.Server.
func (server *Server) Handler() http.Handler {

	server.InitModules()
	return server.engine
}

// mustNotBeInitialized panics when called after module routes are registered,
// since later changes would silently not apply to existing routes.
func (server *Server) mustNotBeInitialized(action string) {
	if server.initialized {

		panic(fmt.Sprintf("solanum: %s must be called before routes are initialized (Run, Handler, InitModules)", action))
	}
}

// SetModules replaces all registered Module implementations with the provided list.
func (server *Server) SetModules(m ...Module) {

	server.mustNotBeInitialized("SetModules")
	server.modules = m
}

// AddModules appends one or more Module implementations to the Runner.
func (server *Server) AddModules(m ...Module) {

	server.mustNotBeInitialized("AddModules")
	server.modules = append(server.modules, m...)
}

// Modules returns the slice of all registered Module pointers.
func (server *Server) Modules() []Module {

	return server.modules
}

// InitGlobalMiddlewares registers the recovery middleware followed by
// the application-wide middlewares configured with WithMiddlewares.
func (server *Server) InitGlobalMiddlewares() {
	if server.recovery != nil {

		server.engine.Use(server.recovery)
	}

	server.engine.Use(server.middlewares...)
}

// Cors applies configured CORS settings to the Gin engine using the cors middleware.
// Accepts functional options for customizing allowed origins, methods, headers, etc.
// It panics when the options are insecure (e.g. credentials with all origins)
// or when called after routes are initialized.
func (server *Server) Cors(opts ...func(*cors2.Option)) {

	server.mustNotBeInitialized("Cors")
	options := cors2.Options(opts...)
	if err := options.Validate(); err != nil {

		panic(err.Error())
	}

	if options.AllowsAllOrigins() {

		server.logger.Warn("cors allows all origins; set cors.WithUrls or cors.WithOriginFunc to restrict them")
	}

	server.engine.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     options.Urls,
				AllowMethods:     options.Methods,
				AllowHeaders:     options.Headers,
				AllowCredentials: options.AllowCredentials,
				AllowOriginFunc:  options.OriginFunc,
				MaxAge:           options.MaxAge,
			},
		),
	)
}

// GinEngine returns the underlying *gin.Engine for direct access and customization.
func (server *Server) GinEngine() *gin.Engine {

	return server.engine
}

// Port returns the configured port for the HTTP server.
func (server *Server) Port() int {
	return server.port
}

// Option configures a Server.
type Option func(*Server)

// WithPort sets the TCP port the server listens on.
func WithPort(port int) Option {
	return func(r *Server) {
		r.port = port
	}
}

// WithGinMode sets the global gin mode (gin.DebugMode, gin.ReleaseMode, gin.TestMode)
// before the engine is created. Note that gin mode is process-wide.
func WithGinMode(mode string) Option {
	return func(r *Server) {
		r.ginMode = mode
	}
}

// WithRecovery replaces the default gin.Recovery() middleware (e.g. with gin.CustomRecovery).
// Passing nil disables panic recovery.
func WithRecovery(recovery gin.HandlerFunc) Option {
	return func(r *Server) {
		r.recovery = recovery
	}
}

// WithMiddlewares appends application-wide middlewares (logger, authentication, ...)
// that run before every module's middlewares.
func WithMiddlewares(middlewares ...gin.HandlerFunc) Option {
	return func(r *Server) {
		r.middlewares = append(r.middlewares, middlewares...)
	}
}

// WithReadHeaderTimeout sets http.Server.ReadHeaderTimeout (zero means no timeout).
func WithReadHeaderTimeout(d time.Duration) Option {
	return func(r *Server) {
		r.readHeaderTimeout = d
	}
}

// WithReadTimeout sets http.Server.ReadTimeout (zero means no timeout).
func WithReadTimeout(d time.Duration) Option {
	return func(r *Server) {
		r.readTimeout = d
	}
}

// WithWriteTimeout sets http.Server.WriteTimeout (zero means no timeout).
func WithWriteTimeout(d time.Duration) Option {
	return func(r *Server) {
		r.writeTimeout = d
	}
}

// WithIdleTimeout sets http.Server.IdleTimeout (zero falls back to the read timeout).
func WithIdleTimeout(d time.Duration) Option {
	return func(r *Server) {
		r.idleTimeout = d
	}
}

// WithShutdownTimeout sets how long graceful shutdown waits for in-flight requests.
func WithShutdownTimeout(d time.Duration) Option {
	return func(r *Server) {
		r.shutdownTimeout = d
	}
}

// WithLogger sets the logger used by the framework itself (defaults to slog.Default()).
// Passing nil silences framework logs.
func WithLogger(logger *slog.Logger) Option {
	return func(r *Server) {
		if logger == nil {

			logger = slog.New(slog.DiscardHandler)
		}

		r.logger = logger
	}
}

// NewSolanum creates and returns a new Server configured with the given options.
// gin.Recovery() is registered by default; see WithRecovery.
func NewSolanum(opts ...Option) *Server {

	s := &Server{
		port:     0,
		recovery: gin.Recovery(),
		logger:   slog.Default(),

		readHeaderTimeout: DefaultReadHeaderTimeout,
		readTimeout:       DefaultReadTimeout,
		writeTimeout:      DefaultWriteTimeout,
		idleTimeout:       DefaultIdleTimeout,
		shutdownTimeout:   DefaultShutdownTimeout,
	}

	for _, opt := range opts {
		if opt != nil {

			opt(s)
		}
	}

	if s.ginMode != "" {

		gin.SetMode(s.ginMode)
	}

	s.engine = gin.New()
	s.InitGlobalMiddlewares()
	return s
}
