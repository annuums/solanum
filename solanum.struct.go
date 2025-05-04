package solanum

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type (
	SolaModule struct {
		uri             string             // base URI path for the module (e.g., "/users")
		controllers     []Controller       // registered controllers for this module
		preMiddlewares  []gin.HandlerFunc  // middleware to run before each handler
		postMiddlewares []gin.HandlerFunc  // middleware to run after each handler
		dependencies    []DependencyConfig // dependencies to inject via DI middleware
	}

	SolaController struct {
		handlers []SolaService
	}

	SolaService struct {
		Uri     string
		Method  string
		Handler gin.HandlerFunc
	}

	runner struct {
		Engine  *gin.Engine
		port    int
		modules []*Module
	}
)

// NewModule 새로운 모듈을 만듭니다. 이 때, 요청받은 router의 uri가 이미 등록되어 있다면 panic
func NewModule(uri string) *SolaModule {
	return &SolaModule{
		uri:             uri,
		controllers:     []Controller{},
		preMiddlewares:  []gin.HandlerFunc{},
		postMiddlewares: []gin.HandlerFunc{},
		dependencies:    []DependencyConfig{},
	}
}

func (m *SolaModule) PreMiddlewares() []gin.HandlerFunc {
	return m.preMiddlewares
}

func (m *SolaModule) PostMiddlewares() []gin.HandlerFunc {
	return m.postMiddlewares
}

func (m *SolaModule) SetPreMiddleware(middlewares ...gin.HandlerFunc) {
	m.preMiddlewares = make([]gin.HandlerFunc, len(middlewares))
	m.preMiddlewares = append(m.preMiddlewares, middlewares...)
}

func (m *SolaModule) SetPostMiddleware(middlewares ...gin.HandlerFunc) {
	m.postMiddlewares = make([]gin.HandlerFunc, len(middlewares))
	m.postMiddlewares = append(m.postMiddlewares, middlewares...)
}

func (m *SolaModule) AddPreMiddleware(middleware gin.HandlerFunc) {
	m.preMiddlewares = append(m.preMiddlewares, middleware)
}

func (m *SolaModule) AddPostMiddleware(middleware gin.HandlerFunc) {
	m.postMiddlewares = append(m.postMiddlewares, middleware)
}

func (m *SolaModule) Controllers() []Controller {
	return m.controllers
}
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

func (m *SolaModule) Uri() string {
	return m.uri
}

// NewController 새로운 Controller를 생성합니다.
func NewController() *SolaController {
	return &SolaController{
		handlers: nil,
	}
}

func (ctr *SolaController) SetHandlers(handlers ...SolaService) {
	if ctr.handlers == nil {
		ctr.handlers = make([]SolaService, 0)
	}

	// ctr.handlers = append(ctr.handlers, *svc)
	ctr.handlers = append(ctr.handlers, handlers...)
}
func (ctr *SolaController) Handlers() []SolaService {
	return ctr.handlers
}
