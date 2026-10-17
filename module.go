package solanum

import (
	"fmt"
	"strings"

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

// ModuleOption configures a SolaModule. NewModule panics when an option returns an error.
type ModuleOption func(*SolaModule) error

// WithUri sets the module's base URI prefix. It must be empty or start with '/'.
func WithUri(uri string) ModuleOption {

	return func(m *SolaModule) error {
		if uri != "" && !strings.HasPrefix(uri, "/") {

			return fmt.Errorf("uri %q must start with '/'", uri)
		}

		m.uri = uri
		return nil
	}
}

// NewModule creates a new SolaModule with the given URI prefix.
// The module starts with empty controller and middleware lists.
func NewModule(opts ...ModuleOption) *SolaModule {

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
// The final handler chain per route is: module pre-middlewares → service PreHandlers → handler → service PostHandlers → module post-middlewares.
// Post handlers run after the rest of the chain completes, even when it was aborted
// (check c.IsAborted() to tell); they are skipped only when the chain panics.
// It panics when a service is nil, has no Handler, or has no Method.
func (m *SolaModule) SetRoutes(router *gin.RouterGroup) {
	for _, c := range m.controllers {
		for _, svc := range c.Handlers() {
			if err := m.validateService(svc); err != nil {

				panic(err.Error())
			}

			posts := make([]gin.HandlerFunc, 0, len(svc.PostHandlers)+len(m.postMiddlewares))
			posts = append(posts, svc.PostHandlers...)
			posts = append(posts, m.postMiddlewares...)

			chain := make([]gin.HandlerFunc, 0, 1+len(m.preMiddlewares)+len(svc.PreHandlers)+1)
			if len(posts) > 0 {

				chain = append(chain, runAfter(posts))
			}

			chain = append(chain, m.preMiddlewares...)
			chain = append(chain, svc.PreHandlers...)
			chain = append(chain, svc.Handler)
			router.Handle(svc.Method, svc.Uri, chain...)
		}
	}
}

// runAfter returns a handler that runs the rest of the chain first and then the given
// post handlers, so they execute even when a pre-handler or the handler aborts.
func runAfter(posts []gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Next()
		for _, post := range posts {
			if post != nil {

				post(c)
			}
		}
	}
}

// validateService reports why a service cannot be registered as a route.
func (m *SolaModule) validateService(svc *SolaService) error {
	if svc == nil {

		return fmt.Errorf("solanum: module %q has a nil SolaService", m.uri)
	}

	if svc.Method == "" {

		return fmt.Errorf("solanum: SolaService %q in module %q has no Method", svc.Uri, m.uri)
	}

	if svc.Handler == nil {

		return fmt.Errorf("solanum: SolaService %s %q in module %q has no Handler", svc.Method, svc.Uri, m.uri)
	}

	return nil
}

// Uri returns the base URI prefix for this module.
func (m *SolaModule) Uri() string {

	return m.uri
}
