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

	fmt.Println("Solanum is running on ", addr)
	server.Engine.Run(addr)
}

func (server *runner) appendHelloWorld() *Module {
	//* Init Vars
	// helloWorldController, err := NewHelloWorldController()

	// if err != nil {
	// 	panic(err.Error())
	// }

	helloUri := "/"
	helloWorldModule, err := NewHelloWorldModule(
		server.Engine.Group(helloUri),
		helloUri,
	)
	if err != nil {
		panic(err.Error())
	}

	return &helloWorldModule
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

func NewSolanum(port int) *Runner {
	if SolanumRunner == nil {
		SolanumRunner = &runner{
			Engine: gin.New(),
			port:   port,
		}
	}

	SolanumRunner.InitGlobalMiddlewares()

	helloUri := "/"
	helloWorldModule, _ := NewHelloWorldModule(
		SolanumRunner.(*runner).Engine.Group(helloUri),
		helloUri,
	)
	SolanumRunner.AddModule(&helloWorldModule)

	SolanumRunner.InitModules()

	return &SolanumRunner
}

func Run() {
	if SolanumRunner == nil {
		SolanumRunner = *NewSolanum(5050)
	}

	SolanumRunner.Run()
}
