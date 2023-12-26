package solanum

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type (
	Runner interface {
		InitModules()
		InitGlobalMiddlewares()
		AddModule(m ...*Module)
		GetModules() []*Module
		GetGinEngine() *gin.Engine

		Run()
	}
	runner struct {
		Engine  *gin.Engine
		port    int
		modules []*Module
	}
)

var SolanumRunner Runner

var (
	CorsDefaultMethods     = []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"}
	CorsDefaultHeaders     = []string{"Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers"}
	CorsDefaultCredentials = false
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
		var _m Module = *m
		_m.SetRoutes()
	}
}

func (server *runner) AddModule(m ...*Module) {
	if server.modules == nil {
		server.modules = make([]*Module, 0)
	}

	server.modules = append(server.modules, m...)
}

func (server *runner) GetModules() []*Module {
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

func (server *runner) GetGinEngine() *gin.Engine {
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
