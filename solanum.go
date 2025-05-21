package solanum

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SolanumRunner holds the global Runner instance used to configure and start the server.
var SolanumRunner Runner

// Default CORS settings
var (
	CorsDefaultMethods      = []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"}
	CorsDefaultHeaders      = []string{"Access-Control-Allow-Headers, Origin, Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers"}
	CorsDefaultCredentials  = false
	CorsDefaultOriginalFunc = func(origin string) bool {
		// Default origin function allows any localhost origin
		return strings.Contains(origin, "://localhost")
	}
)

const (
	// DependencyPrefix is the prefix used in Gin context keys for injected dependencies.
	DependencyPrefix = "__sol_dep__"
)

// ValidateDependencies checks all registered modules for their dependencies.
func (server *runner) ValidateDependencies() error {
	for _, mPtr := range server.modules {
		m := *mPtr
		for _, dep := range m.Dependencies() {
			var err error
			if dep.Type != nil {
				_, err = ResolveByType(dep.Key, dep.Type)
			} else {
				_, err = Resolve(dep.Key)
			}
			if err != nil {
				return fmt.Errorf("dependency validation failed for key=%q type=%v: %w",
					dep.Key, dep.Type, err)
			}
		}
	}
	return nil
}

// Run initializes all modules and starts the Gin HTTP server on the configured port.
func (server *runner) Run() {
	addr := fmt.Sprintf(":%v", server.port)

	SolanumRunner.InitModules()

	if err := server.ValidateDependencies(); err != nil {
		log.Fatalf("❌ Dependency check failed: %v", err)
	}

	log.Println("Solanum is running on ", addr)
	server.Engine.Run(addr)
}

// InitModules sets up routing groups for each Module and applies their routes.
func (server *runner) InitModules() {
	log.Println("Initialize Modules...")
	for _, m := range server.modules {
		(*m).SetRoutes(
			server.GinEngine().Group(
				(*m).Uri(),
			),
		)
	}
}

// SetModules registers one or more Module implementations with the Runner.
func (server *runner) SetModules(m ...Module) {
	if server.modules == nil {
		server.modules = make([]*Module, 0)
	}

	for i := range m {
		server.modules = append(server.modules, &m[i])
	}
}

// Modules returns the slice of all registered Module pointers.
func (server *runner) Modules() []*Module {
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
func (server *runner) Cors(opts ...func(*CorsOption)) {
	options := CorsOptions(opts...)

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

// NewSolanum creates (once) and returns the global Runner configured for the given port.
// It ensures global middlewares are initialized. Subsequent calls return the same Runner.
func NewSolanum(port int) *Runner {
	if SolanumRunner == nil {
		SolanumRunner = &runner{
			Engine: gin.New(),
			port:   port,
		}
	}

	SolanumRunner.InitGlobalMiddlewares()

	return &SolanumRunner
}
