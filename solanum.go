package solanum

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var SolanumRunner Runner

var (
	CorsDefaultMethods      = []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"}
	CorsDefaultHeaders      = []string{"Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers"}
	CorsDefaultCredentials  = false
	CorsDefaultOriginalFunc = func(origin string) bool {
		return strings.Contains(origin, ":://localhost")
	}
)

func (server *runner) Run() {
	addr := fmt.Sprintf(":%v", server.port)

	SolanumRunner.InitModules()

	fmt.Println("Solanum is running on ", addr)
	server.Engine.Run(addr)
}

func (server *runner) InitModules() {
	//* setRoutes
	fmt.Println("Initialize Modules...")
	for _, m := range server.modules {
		(*m).SetRoutes(
			server.GinEngine().Group(
				(*m).Uri(),
			),
		)
	}
}

func (server *runner) SetModules(m ...Module) {
	if server.modules == nil {
		server.modules = make([]*Module, 0)
	}

	for _, module := range m {
		server.modules = append(server.modules, &module)
	}

}

func (server *runner) Modules() []*Module {
	return server.modules
}

func (server *runner) InitGlobalMiddlewares() {
	//* 1. Logger, ...

	//* 2. Authentication, ...

	//* 3. Authorization, ...
}

func (server *runner) Cors(url, headers, methods []string, allowCredentials bool, originFunc func(origin string) bool, maxAge int) {
	server.Engine.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     url,
				AllowMethods:     methods,
				AllowHeaders:     headers,
				AllowCredentials: allowCredentials,
				AllowOriginFunc:  originFunc,
				MaxAge:           time.Duration(maxAge) * time.Hour,
			},
		),
	)
}

func (server *runner) GinEngine() *gin.Engine {
	return server.Engine
}

func NewSolanum(port int) *Runner {
	if SolanumRunner == nil {
		SolanumRunner = &runner{
			Engine: gin.New(),
			port:   port,
		}
	}

	SolanumRunner.InitGlobalMiddlewares()

	return &SolanumRunner
}
