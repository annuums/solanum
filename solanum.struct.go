package solanum

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type (
	SolaModule struct {
		uri         string
		controllers []Controller
		middlewares []gin.HandlerFunc
	}

	SolaController struct {
		handlers []SolaService
	}

	SolaService struct {
		Uri        string
		Method     string
		Handler    gin.HandlerFunc
		Middleware gin.HandlerFunc
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
		uri:         uri,
		controllers: []Controller{},
		middlewares: []gin.HandlerFunc{},
	}
}

func (m *SolaModule) GlobalMiddlewares() []gin.HandlerFunc {
	return m.middlewares
}
func (m *SolaModule) SetGlobalMiddleware(middlewares ...gin.HandlerFunc) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *SolaModule) Controllers() []Controller {
	return m.controllers
}
func (m *SolaModule) SetControllers(c ...Controller) {
	m.controllers = append(m.controllers, c...)
}

func (m *SolaModule) SetRoutes(router *gin.RouterGroup) {
	for _, c := range m.controllers {
		services := (c).Handlers()

		for _, svc := range services {
			switch svc.Method {
			case http.MethodGet:
				router.GET(svc.Uri, svc.Handler)
			case http.MethodPost:
				router.POST(svc.Uri, svc.Handler)
			case http.MethodPut:
				router.PUT(svc.Uri, svc.Handler)
			case http.MethodPatch:
				router.PATCH(svc.Uri, svc.Handler)
			case http.MethodDelete:
				router.DELETE(svc.Uri, svc.Handler)
			case http.MethodHead:
				router.HEAD(svc.Uri, svc.Handler)
			case http.MethodOptions:
				router.OPTIONS(svc.Uri, svc.Handler)
			default:
				log.Fatalf("Unknown method registered: %v", svc)
			}
		}
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
