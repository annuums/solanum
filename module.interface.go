package solanum

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	Controller interface {
		AddHandler(handler ...*Service)
		GetHandlers() []*Service
	}

	Module interface {
		//* Middlewares
		GetGlobalMiddlewares() []*gin.HandlerFunc
		SetGlobalMiddleware(middlewares ...*gin.HandlerFunc)

		//* Controllers
		GetControllers() []*Controller
		SetControllers(c ...*Controller)

		//* Controllers -> Handlers
		SetRoutes()
	}

	module struct {
		uri         string
		controllers []*Controller
		middlewares []*gin.HandlerFunc
		router      *gin.RouterGroup
	}

	controller struct {
		handlers []*Service
	}

	Service struct {
		Uri        string
		Method     string
		Handler    gin.HandlerFunc
		Middleware gin.HandlerFunc
	}
)

/*
새로운 모듈을 만듭니다. 이 때, 요청받은 router의 uri가 이미 등록되어 있다면 error를 반환합니다.

	//* Middlewares
	GetGlobalMiddlewares() []*gin.HandlerFunc
	SetGlobalMiddleware(middlewares ...*gin.HandlerFunc)

	//* Controllers
	GetControllers() []*Controller
	SetControllers(c ...*Controller)

	//* Controllers -> Handlers
	SetRoutes()
*/
func NewModule(router *gin.RouterGroup, uri string) (*Module, error) {
	var m Module = &module{
		uri:         uri,
		router:      router,
		controllers: []*Controller{},
		middlewares: []*gin.HandlerFunc{},
	}
	return &m, nil
}

func (m *module) GetGlobalMiddlewares() []*gin.HandlerFunc {
	return m.middlewares
}
func (m *module) SetGlobalMiddleware(middlewares ...*gin.HandlerFunc) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *module) GetControllers() []*Controller {
	return m.controllers
}
func (m *module) SetControllers(c ...*Controller) {
	m.controllers = append(m.controllers, c...)
}

func (m *module) SetRoutes() {
	for _, c := range m.controllers {
		ctr, ok := (*c).(*controller)

		if !ok {
			log.Fatal("Fail to set routes for controller::", c)
		}

		services := ctr.GetHandlers()

		for _, svc := range services {
			switch svc.Method {
			case http.MethodGet:
				m.router.GET(svc.Uri, svc.Handler)
			case http.MethodPost:
				m.router.POST(svc.Uri, svc.Handler)
			case http.MethodPut:
				m.router.PUT(svc.Uri, svc.Handler)
			case http.MethodPatch:
				m.router.PATCH(svc.Uri, svc.Handler)
			case http.MethodDelete:
				m.router.DELETE(svc.Uri, svc.Handler)
			case http.MethodHead:
				m.router.HEAD(svc.Uri, svc.Handler)
			case http.MethodOptions:
				m.router.OPTIONS(svc.Uri, svc.Handler)
			default:
				log.Fatalf("Unknown method registered: %v", svc)
			}
		}
	}
}

//* Middlewares

/*
새로운 컨트롤러를 만듭니다.
*/
func NewController() (*Controller, error) {
	var ctr Controller = &controller{
		handlers: nil,
	}
	return &ctr, nil
}

func (ctr *controller) AddHandler(svc ...*Service) {
	if ctr.handlers == nil {
		ctr.handlers = make([]*Service, 0)
	}

	// ctr.handlers = append(ctr.handlers, *svc)
	ctr.handlers = append(ctr.handlers, svc...)
}
func (ctr *controller) GetHandlers() []*Service {
	return ctr.handlers
}
