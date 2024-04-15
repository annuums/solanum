package solanum

import (
	"github.com/gin-gonic/gin"
)

type (
	Controller interface {
		Handlers() []SolaService
		SetHandlers(handler ...SolaService)
	}

	Module interface {
		//* Middlewares
		GlobalMiddlewares() []gin.HandlerFunc
		SetGlobalMiddleware(middlewares ...gin.HandlerFunc)

		//* Controllers
		Controllers() []Controller
		SetControllers(c ...Controller)

		//* Controllers -> Handlers
		SetRoutes()
	}

	Runner interface {
		InitModules()
		InitGlobalMiddlewares()
		Modules() []*Module
		SetModules(m ...*Module)
		GinEngine() *gin.Engine
		Cors(url, headers, methods []string, allowCredentials bool, originFunc func(origin string) bool, maxAge int)

		Run()
	}
)
