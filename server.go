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
	engine  *gin.Engine // underlying Gin engine
	port    int         // TCP port to listen on
	modules []Module    // registered modules
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

// InitGlobalMiddlewares is a placeholder for registering application-wide middlewares
// such as logging, authentication, and authorization. Implement as needed.
func (server *Server) InitGlobalMiddlewares() {
	//* 1. Logger, ...

	//* 2. Authentication, ...

	//* 3. Authorization, ...
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

// NewSolanum creates and returns a new Server configured with the given options.
func NewSolanum(opts ...Option) *Server {

	s := &Server{
		engine: gin.New(),
		port:   0,
	}

	for _, opt := range opts {
		if opt != nil {

			opt(s)
		}
	}

	s.InitGlobalMiddlewares()
	return s
}
