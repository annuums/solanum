package solanum

import (
	"fmt"
	"github.com/annuums/solanum/container"
	"github.com/annuums/solanum/util"
	"reflect"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SolanumRunner holds the global Runner instance used to configure and start the server.
var SolanumRunner Runner

// ValidateDependencies checks all registered modules for their dependencies.
func (server *runner) ValidateDependencies() error {

	for _, mPtr := range server.modules {

		for _, dep := range *(*mPtr).Dependencies() {

			inst, err := container.Resolve(dep.Key)
			if err != nil {

				return fmt.Errorf(
					"dependency validation failed for key=%q :: %w",
					dep.Key,
					err,
				)
			}

			if dep.Type != nil {

				instType := reflect.TypeOf(inst)

				switch dep.Type.Kind() {
				case reflect.Interface:

					if !instType.Implements(dep.Type) {

						return fmt.Errorf(
							"dependency %q: instance type %v does not implement %v",
							dep.Key,
							instType,
							dep.Type,
						)
					}

				default:

					if !instType.AssignableTo(dep.Type) {

						// If the instance type is not assignable to the dependency type,
						return fmt.Errorf(
							"dependency %q: instance type %v not assignable to %v",
							dep.Key,
							instType,
							dep.Type,
						)
					}
				}
			}
		}
	}

	return nil
}

// Run initializes all modules and starts the Gin HTTP server on the configured port.
func (server *runner) Run() {

	if err := server.ValidateDependencies(); err != nil {

		panic("Dependency check failed :: " + err.Error())
	}

	if server.port == nil {

		panic("Server port is not configured. Please set a port before running.")
	}

	// Start Gin Server
	if server.Port() != 0 {

		SolanumRunner.InitModules()

		addr := fmt.Sprintf(":%d", *server.port)
		fmt.Printf("Solanum is running on %s\n", addr)

		if err := server.Engine.Run(addr); err != nil {

			panic("fail to run server on addr :: " + addr + " :: " + err.Error())
		}
	}
}

// InitModules sets up routing groups for each Module and applies their routes.
func (server *runner) InitModules() {

	fmt.Println("Initialize Modules...")

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
func (server *runner) Cors(opts ...func(*util.CorsOption)) {

	options := util.CorsOptions(opts...)

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

	if server.port == nil {

		return 0
	}

	return *server.port
}

type option func(Runner)

func WithPort(port int) option {

	return func(r Runner) {

		if runner, ok := r.(*runner); ok {

			runner.port = &port
		} else {

			fmt.Println("⚠️ Unable to set port: Runner is not of type *runner")
		}
	}
}

// NewSolanum creates (once) and returns the global Runner configured for the given port.
// It ensures global middlewares are initialized. Subsequent calls return the same Runner.
func NewSolanum(opts ...option) Runner {

	if SolanumRunner == nil {

		SolanumRunner = &runner{}
	}

	for _, opt := range opts {

		if opt != nil {
			opt(SolanumRunner)
		}
	}

	port := SolanumRunner.Port()

	if port == 0 {

		fmt.Println("⚠️ No port specified, using default port 0 (random port).")
		port = 0
	}

	if port == 0 {

		SolanumRunner = &runner{}
	} else {

		SolanumRunner = &runner{
			Engine: gin.New(),
			port:   &port,
		}

		SolanumRunner.InitGlobalMiddlewares()
	}

	return SolanumRunner
}
