package solanum

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type (
	// DependencyConfig defines a key and Go type for a dependency to be injected into handlers.
	DependencyConfig struct {
		// Key identifier used when registering and retrieving the dependency
		// registration key for the dependency
		Key string

		// Type reflect.Type of the interface or concrete type to resolve
		// expected reflect.Type for resolution
		Type reflect.Type
	}

	// Controller declares a set of service handlers for a logical grouping of routes.
	Controller interface {
		Handlers() []SolaService
		SetHandlers(handler ...SolaService)
	}

	Module interface {
		//* Middlewares
		PreMiddlewares() []gin.HandlerFunc
		AddPreMiddleware(middleware gin.HandlerFunc)
		SetPreMiddleware(middlewares ...gin.HandlerFunc)
		PostMiddlewares() []gin.HandlerFunc
		AddPostMiddleware(middleware gin.HandlerFunc)
		SetPostMiddleware(middlewares ...gin.HandlerFunc)

		//* Controllers
		Controllers() []Controller
		SetControllers(c ...Controller)

		// Dependencies returns the list of dependencies that this module injects.
		Dependencies() []DependencyConfig

		// SetDependencies defines which dependencies to inject via middleware.
		SetDependencies(deps ...DependencyConfig)

		// SetRoutes binds the module's controllers, middleware, and routes onto a RouterGroup.
		SetRoutes(router *gin.RouterGroup)

		// Properties
		Uri() string
	}

	Runner interface {
		InitModules()
		InitGlobalMiddlewares()
		Modules() []*Module
		SetModules(m ...Module)
		GinEngine() *gin.Engine
		Cors(opts ...func(*CorsOption))

		Run()
	}
)
