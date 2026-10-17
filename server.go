package solanum

import (
	"fmt"

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
}

// Run initializes all modules and starts the Gin HTTP server on the configured port.
func (server *Server) Run() {
	if server.port <= 0 {

		panic("Server port is not configured. Please set a port before running.")
	}

	server.InitModules()
	addr := fmt.Sprintf(":%d", server.port)

	fmt.Printf("Solanum is running on %s\n", addr)
	if err := server.engine.Run(addr); err != nil {

		panic("fail to run server on addr :: " + addr + " :: " + err.Error())
	}
}

// InitModules sets up routing groups for each Module and applies their routes.
func (server *Server) InitModules() {

	fmt.Printf("Initialize Modules...\n")
	for _, m := range server.modules {

		m.SetRoutes(
			server.GinEngine().Group(
				m.Uri(),
			),
		)
	}
}

// SetModules registers one or more Module implementations with the Runner.
func (server *Server) SetModules(m ...Module) {
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
// It panics when the options are insecure (e.g. credentials with all origins).
func (server *Server) Cors(opts ...func(*cors2.Option)) {

	options := cors2.Options(opts...)
	if err := options.Validate(); err != nil {

		panic(err.Error())
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

// NewSolanum creates and returns a new Server configured with the given options.
// gin.Recovery() is registered by default; see WithRecovery.
func NewSolanum(opts ...Option) *Server {

	s := &Server{
		port:     0,
		recovery: gin.Recovery(),
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
