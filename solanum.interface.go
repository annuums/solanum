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
		PreMiddlewares() []gin.HandlerFunc
		SetPreMiddleware(middlewares ...gin.HandlerFunc)
		PostMiddlewares() []gin.HandlerFunc
		SetPostMiddleware(middlewares ...gin.HandlerFunc)

		//* Controllers
		Controllers() []Controller
		SetControllers(c ...Controller)

		//* Controllers -> Handlers
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
		Cors(url, headers, methods []string, allowCredentials bool, originFunc func(origin string) bool, maxAge int)

		Run()
	}
)
