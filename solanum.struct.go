package solanum

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type (
	// SolaModule encapsulates a self-contained HTTP module with its own URI prefix,
	// controllers, middleware stacks, and dependency configurations.
	SolaModule struct {
		uri             string             // base URI path for the module (e.g., "/users")
		controllers     []Controller       // registered controllers for this module
		preMiddlewares  []gin.HandlerFunc  // middleware to run before each handler
		postMiddlewares []gin.HandlerFunc  // middleware to run after each handler
		dependencies    []DependencyConfig // dependencies to inject via DI middleware
	}

	// SolaController groups one or more SolaService handlers under a logical controller.
	// It implements the Controller interface, managing a list of SolaService entries.
	SolaController struct {
		handlers []*SolaService // handlers service handlers defined for this controller
	}

	// SolaService represents a single HTTP route handler configuration.
	SolaService struct {
		// Uri the relative path for the service (e.g., "/:id")
		// route URI relative to module prefix
		Uri string

		// Method the HTTP method to bind (GET, POST, etc.)
		Method string

		// Handler the Gin handler function to execute
		Handler gin.HandlerFunc
	}

	// runner implements the Runner interface and drives the application startup.
	// It holds the Gin engine, listening port, and registered modules.
	runner struct {
		Engine  *gin.Engine // underlying Gin engine
		port    int         // TCP port to listen on
		modules []*Module   // pointers to registered modules
	}
)

// NewModule creates a new SolaModule with the given URI prefix.
// The module starts with empty controller, middleware, and dependency lists.
func NewModule(uri string) *SolaModule {
	return &SolaModule{
		uri:             uri,
		controllers:     []Controller{},
		preMiddlewares:  []gin.HandlerFunc{},
		postMiddlewares: []gin.HandlerFunc{},
		dependencies:    []DependencyConfig{},
	}
}

// PreMiddlewares returns the list of middleware to execute before handlers.
func (m *SolaModule) PreMiddlewares() []gin.HandlerFunc {
	return m.preMiddlewares
}

// PostMiddlewares returns the list of middleware to execute after handlers.
func (m *SolaModule) PostMiddlewares() []gin.HandlerFunc {
	return m.postMiddlewares
}

// SetPreMiddlewares replaces the pre-handler middleware chain with the provided list.
func (m *SolaModule) SetPreMiddlewares(middlewares ...gin.HandlerFunc) {
	m.preMiddlewares = make([]gin.HandlerFunc, 0)
	m.preMiddlewares = append(m.preMiddlewares, middlewares...)
}

// SetPostMiddlewares replaces the post-handler middleware chain with the provided list.
func (m *SolaModule) SetPostMiddlewares(middlewares ...gin.HandlerFunc) {
	m.postMiddlewares = make([]gin.HandlerFunc, 0)
	m.postMiddlewares = append(m.postMiddlewares, middlewares...)
}

// AddPreMiddleware appends a single middleware to the pre-handler chain.
func (m *SolaModule) AddPreMiddleware(middleware gin.HandlerFunc) {
	m.preMiddlewares = append(m.preMiddlewares, middleware)
}

// AddPostMiddleware appends a single middleware to the post-handler chain.
func (m *SolaModule) AddPostMiddleware(middleware gin.HandlerFunc) {
	m.postMiddlewares = append(m.postMiddlewares, middleware)
}

// Controllers returns all controllers registered in this module.
func (m *SolaModule) Controllers() []Controller {
	return m.controllers
}

// SetControllers registers one or more Controller implementations for this module.
func (m *SolaModule) SetControllers(c ...Controller) {
	m.controllers = append(m.controllers, c...)
}

// Dependencies returns the list of DependencyConfig entries for this module.
func (m *SolaModule) Dependencies() []DependencyConfig { return m.dependencies }

// SetDependencies replaces the module's dependency list with the provided configs.
func (m *SolaModule) SetDependencies(deps ...DependencyConfig) {
	m.dependencies = make([]DependencyConfig, len(deps))
	copy(m.dependencies, deps)
}

// SetRoutes registers the module's routes, middleware, and DI middleware on the given RouterGroup.
// It applies DI if dependencies are defined, then mounts each SolaService handler with pre- and post-middleware.
func (m *SolaModule) SetRoutes(router *gin.RouterGroup) {
	// Apply DI middleware if dependencies exist
	if len(m.dependencies) > 0 {
		router.Use(diMiddleware(m.dependencies))
	}

	// Iterate controllers and their services
	for _, c := range m.controllers {
		ctr, ok := c.(*SolaController)
		if !ok {
			panic(fmt.Sprintf("controller is not *SolaController: %T", c))
		}
		for _, svc := range ctr.handlers {
			// pre → handler → post
			chain := append(append(m.preMiddlewares, svc.Handler), m.postMiddlewares...)
			router.Handle(svc.Method, svc.Uri, chain...)
		}
	}
}

// diMiddleware returns a Gin middleware that resolves and injects dependencies for each request.
// It sets each dependency instance in the context under DependencyPrefix+key.
func diMiddleware(deps []DependencyConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		seen := make(map[string]struct{}, len(deps))
		for _, d := range deps {
			log.Printf("key :: %v, interface :: %v", d.Key, d.Type)
			if _, dup := seen[d.Key]; dup {
				panic(fmt.Sprintf("duplicate dependency key: %q", d.Key))
			}
			seen[d.Key] = struct{}{}

			var inst interface{}
			var err error
			if d.Type != nil {
				inst, err = ResolveByType(d.Key, d.Type)
			} else {
				inst, err = Resolve(d.Key)
			}

			if err != nil {
				panic(fmt.Errorf("failed to resolve %q: %w", d.Key, err))
			}

			depKey := DependencyPrefix + d.Key
			c.Set(depKey, inst)
		}
		c.Next()
	}
}

// Uri returns the base URI prefix for this module.
func (m *SolaModule) Uri() string {
	return m.uri
}

// NewController constructs an empty SolaController ready to receive handlers.
func NewController() *SolaController {
	return &SolaController{
		handlers: make([]*SolaService, 0),
	}
}

// SetHandlers appends one or more SolaService entries to the controller's handler list.
func (ctr *SolaController) SetHandlers(handlers ...*SolaService) {
	if ctr.handlers == nil {
		ctr.handlers = make([]*SolaService, 0)
	}

	// ctr.handlers = append(ctr.handlers, *svc)
	ctr.handlers = append(ctr.handlers, handlers...)
}

// Handlers returns the slice of SolaService entries managed by this controller.
func (ctr *SolaController) Handlers() []*SolaService {
	return ctr.handlers
}
