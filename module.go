package solanum

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// SolaModule encapsulates a self-contained HTTP module with its own URI prefix,
// controllers, and middleware stacks.
type SolaModule struct {
	uri             string            // base URI path for the module (e.g., "/users")
	controllers     []Controller      // registered controllers for this module
	preMiddlewares  []gin.HandlerFunc // middleware to run before each handler
	postMiddlewares []gin.HandlerFunc // middleware to run after each handler
}

type moduleOption func(*SolaModule) error

func WithUri(uri string) moduleOption {

	return func(m *SolaModule) error {

		m.uri = uri
		return nil
	}
}

// NewModule creates a new SolaModule with the given URI prefix.
// The module starts with empty controller and middleware lists.
func NewModule(opts ...moduleOption) *SolaModule {

	module := &SolaModule{
		controllers:     []Controller{},
		preMiddlewares:  []gin.HandlerFunc{},
		postMiddlewares: []gin.HandlerFunc{},
	}

	for _, opt := range opts {

		if err := opt(module); err != nil {

			panic(fmt.Sprintf("failed to apply module option :: %v", err))
		}
	}

	return module
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

	m.preMiddlewares = middlewares
}

// SetPostMiddlewares replaces the post-handler middleware chain with the provided list.
func (m *SolaModule) SetPostMiddlewares(middlewares ...gin.HandlerFunc) {

	m.postMiddlewares = middlewares
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

	m.controllers = c
}

// AddControllers appends one or more Controller implementations to this module.
func (m *SolaModule) AddControllers(c ...Controller) {

	m.controllers = append(m.controllers, c...)
}

// SetRoutes registers the module's routes and middleware on the given RouterGroup.
// The final handler chain per route is: pre-middlewares → handler → post-middlewares.
func (m *SolaModule) SetRoutes(router *gin.RouterGroup) {

	for _, c := range m.controllers {
		for _, svc := range c.Handlers() {

			chain := make([]gin.HandlerFunc, 0, len(m.preMiddlewares)+len(svc.PreHandlers)+1+len(svc.PostHandlers)+len(m.postMiddlewares))
			chain = append(chain, m.preMiddlewares...)
			chain = append(chain, svc.PreHandlers...)
			chain = append(chain, svc.Handler)
			chain = append(chain, svc.PostHandlers...)
			chain = append(chain, m.postMiddlewares...)
			router.Handle(svc.Method, svc.Uri, chain...)
		}
	}
}

// Uri returns the base URI prefix for this module.
func (m *SolaModule) Uri() string {

	return m.uri
}
