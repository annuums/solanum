package modules

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	Controller interface {
		AddHandler(handler *service)
		GetHandlers() []service
	}

	Module interface {
		//* Middlewares
		GetGlobalMiddlewares() []*gin.HandlerFunc
		SetGlobalMiddleware(middlewares ...*gin.HandlerFunc)

		//* Controllers
		GetControllers() []Controller
		SetControllers(c ...Controller)

		//* Controllers -> Handlers
		SetRoutes()
	}

	module struct {
		uri         string
		controllers []Controller
		middlewares []*gin.HandlerFunc
		router      *gin.RouterGroup
	}

	controller struct {
		handlers []service
	}

	service struct {
		uri        string
		method     string
		handler    gin.HandlerFunc
		middleware gin.HandlerFunc
	}
)

/*
새로운 모듈을 만듭니다. 이 때, 요청받은 router의 uri가 이미 등록되어 있다면 error를 반환합니다.
*/
func NewModule(router *gin.RouterGroup, uri string) (Module, error) {
	//* Mapper Check
	//* Check Duplicated URIs

	return &module{
		uri:         uri,
		router:      router,
		controllers: []Controller{},
		middlewares: []*gin.HandlerFunc{},
	}, nil
}

func (m *module) GetGlobalMiddlewares() []*gin.HandlerFunc {
	return m.middlewares
}
func (m *module) SetGlobalMiddleware(middlewares ...*gin.HandlerFunc) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *module) GetControllers() []Controller {
	return m.controllers
}
func (m *module) SetControllers(c ...Controller) {
	m.controllers = append(m.controllers, c...)
}

func (m *module) SetRoutes() {
	for _, c := range m.controllers {
		services := c.GetHandlers()
		for _, svc := range services {
			switch svc.method {
			case http.MethodGet:
				m.router.GET(svc.uri, svc.handler)
			case http.MethodPost:
				m.router.POST(svc.uri, svc.handler)
			case http.MethodPut:
				m.router.PUT(svc.uri, svc.handler)
			case http.MethodPatch:
				m.router.PATCH(svc.uri, svc.handler)
			case http.MethodDelete:
				m.router.DELETE(svc.uri, svc.handler)
			case http.MethodHead:
				m.router.HEAD(svc.uri, svc.handler)
			case http.MethodOptions:
				m.router.OPTIONS(svc.uri, svc.handler)
			default:
				log.Fatalf("Unknown method registered: %v", svc)
			}
		}
	}
}

// //* Services
// GetServices() *[]Service
// SetServices(...*Service) (*Module, error)

//* Middlewares

/*
새로운 컨트롤러를 만듭니다.
*/
func NewController() (Controller, error) {
	var ctr Controller = &controller{
		handlers: nil,
	}
	return ctr, nil
}

func (ctr *controller) AddHandler(svc *service) {
	ctr.handlers = append(ctr.handlers, *svc)
}
func (ctr *controller) GetHandlers() []service {
	return ctr.handlers
}
