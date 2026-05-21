package solanum

import (
	"fmt"
	"time"

	cors2 "github.com/annuums/solanum/middleware/cors"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run initializes all modules and starts the Gin HTTP server on the configured port.
func (server *runner) Run() {
	if server.port <= 0 {

		panic("Server port is not configured. Please set a port before running.")
	}

	server.InitModules()
	addr := fmt.Sprintf(":%d", server.port)

	fmt.Printf("Solanum is running on %s\n", addr)
	if err := server.Engine.Run(addr); err != nil {

		panic("fail to run server on addr :: " + addr + " :: " + err.Error())
	}
}

// InitModules sets up routing groups for each Module and applies their routes.
func (server *runner) InitModules() {

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
func (server *runner) SetModules(m ...Module) {
	server.modules = append(server.modules, m...)
}

// Modules returns the slice of all registered Module pointers.
func (server *runner) Modules() []Module {

	return server.modules
}

// InitGlobalMiddlewares is a placeholder for registering application-wide middlewares
// such as logging, authentication, and authorization. Implement as needed.
func (server *runner) InitGlobalMiddlewares() {
	//* 1. Logger, ...

	//* 2. Authentication, ...

	//* 3. Authorization, ...
}

// Cors applies configured CORS settings to the Gin engine using the cors middleware.
// Accepts functional options for customizing allowed origins, methods, headers, etc.
func (server *runner) Cors(opts ...func(*cors2.Option)) {

	options := cors2.Options(opts...)

	server.Engine.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     options.Urls,
				AllowMethods:     options.Methods,
				AllowHeaders:     options.Headers,
				AllowCredentials: options.AllowCredentials,
				AllowOriginFunc:  options.OriginFunc,
				MaxAge:           time.Duration(options.MaxAge) * time.Hour,
			},
		),
	)
}

// GinEngine returns the underlying *gin.Engine for direct access and customization.
func (server *runner) GinEngine() *gin.Engine {

	return server.Engine
}

// Port returns the configured port for the HTTP server.
func (server *runner) Port() int {
	return server.port
}

type Option func(*runner)

func WithPort(port int) Option {
	return func(r *runner) {
		r.port = port
	}
}

// NewSolanum creates and returns a new Runner configured with the given options.
func NewSolanum(opts ...Option) Runner {

	s := &runner{
		Engine: gin.New(),
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
