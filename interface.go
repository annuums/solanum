package solanum

import (
	"github.com/annuums/solanum/middleware/cors"
	"github.com/gin-gonic/gin"
)

type (

	// Controller declares a set of service handlers for a logical grouping of routes.
	Controller interface {
		// Handlers returns the slice of service definitions this controller manages.
		Handlers() []*SolaService

		// SetHandlers replaces the controller's service handlers with the provided list.
		SetHandlers(handler ...*SolaService)
	}

	// Module represents a self-contained HTTP module with its own URI prefix,
	// middleware layers, controllers, and dependencies.
	Module interface {
		// PreMiddlewares returns middleware to run before each handler.
		PreMiddlewares() []gin.HandlerFunc

		// AddPreMiddleware appends a middleware to the pre-handler chain.
		AddPreMiddleware(middleware gin.HandlerFunc)

		// SetPreMiddlewares replaces the entire pre-handler middleware chain.
		SetPreMiddlewares(middlewares ...gin.HandlerFunc)

		// PostMiddlewares returns middleware to run after each handler.
		PostMiddlewares() []gin.HandlerFunc

		// AddPostMiddleware appends a middleware to the post-handler chain.
		AddPostMiddleware(middleware gin.HandlerFunc)

		// SetPostMiddlewares replaces the entire post-handler middleware chain.
		SetPostMiddlewares(middlewares ...gin.HandlerFunc)

		// Controllers returns the slice of Controller implementations in this module.
		Controllers() []Controller

		// SetControllers replaces all Controller implementations for this module.
		SetControllers(c ...Controller)

		// AddControllers appends one or more Controller implementations to this module.
		AddControllers(c ...Controller)

		// SetRoutes binds the module's controllers, middleware, and routes onto a RouterGroup.
		SetRoutes(router *gin.RouterGroup)

		// Uri returns the base URI path for this module (e.g., "/users").
		Uri() string
	}

	// Runner is the application entrypoint interface for Solanum.
	// It manages module initialization, global middlewares, CORS, and server start.
	Runner interface {
		// InitModules sets up all registered modules and their routes.
		InitModules()

		// InitGlobalMiddlewares registers any application-wide middleware.
		InitGlobalMiddlewares()

		// Modules returns a slice of pointers to registered Modules.
		Modules() []Module

		// SetModules registers one or more Modules with the Runner.
		SetModules(m ...Module)

		// GinEngine exposes the underlying *gin.Engine for custom setup.
		GinEngine() *gin.Engine

		// Port exposes the configured port for the HTTP server.
		Port() int

		// Cors applies CORS configuration to the Gin engine using functional options.
		Cors(opts ...func(*cors.Option))

		// Run boots the HTTP server, initializing modules and listening on the configured port.
		Run()
	}
)
