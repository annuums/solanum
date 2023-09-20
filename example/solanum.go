package main

import (
	"fmt"

	solanum "github.com/annuums/solanum"
	"github.com/gin-gonic/gin"
)

type (
	Runner interface {
		InitModules()
		GetModules() *[]solanum.Module
		InitGlobalMiddlewares()
		Run()
	}
	runner struct {
		Engine  *gin.Engine
		port    int
		modules *[]solanum.Module
	}
)

var SolanumRunner Runner

func (server *runner) Run() {
	addr := fmt.Sprintf(":%v", server.port)

	fmt.Println("Solanum is running on ", addr)
	server.Engine.Run(addr)
}

func (server *runner) appendHelloWorld() {
	//* Init Vars
	helloWorldController, err := solanum.NewController()

	if err != nil {
		panic(err.Error())
	}

	helloWorldModule, err := solanum.NewModule(server.Engine.Group("/"), "/")
	if err != nil {
		panic(err.Error())
	}

	//* Set Handlers
	helloHandler := solanum.NewHelloHandler()

	helloWorldController.AddHandler(helloHandler)
	helloWorldModule.SetControllers(helloWorldController)

	//* Set Router
	helloWorldModule.SetRoutes()
}

func (server *runner) InitModules() {
	//* initialize Modules
	//* Modules are Singleton Designed
	// myModules := make([]solanum.Module, 0)
	myModules := []solanum.Module{}

	if length := len(myModules); length == 0 {
		server.appendHelloWorld()
	}

	server.modules = &myModules
}

func (server *runner) GetModules() *[]solanum.Module {
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
	SolanumRunner.InitModules()

	return &SolanumRunner
}

func Run() {
	if SolanumRunner == nil {
		SolanumRunner = *NewSolanum(5050)
	}

	SolanumRunner.Run()
}
