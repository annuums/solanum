package solanum

import (
	"fmt"

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
