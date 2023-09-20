package main

import (
	"fmt"

	"github.com/dev-whoan/go-study-net-http/modules"
	"github.com/gin-gonic/gin"
)

type (
	Solanum interface {
		InitModules()
		GetModules() *[]modules.Module
		InitGlobalMiddlewares()
		Run()
	}
	solanum struct {
		Engine  *gin.Engine
		port    int
		modules *[]modules.Module
	}
)

var SolanumRunner Solanum

func (server *solanum) Run() {
	addr := fmt.Sprintf(":%v", server.port)

	fmt.Println("Solanum is running on ", addr)
	server.Engine.Run(addr)
}

func (server *solanum) appendHelloWorld() {
	//* Init Vars
	helloWorldController, err := modules.NewController()

	if err != nil {
		panic(err.Error())
	}

	helloWorldModule, err := modules.NewModule(server.Engine.Group("/"), "/")
	if err != nil {
		panic(err.Error())
	}

	//* Set Handlers
	helloHandler := modules.NewHelloHandler()

	helloWorldController.AddHandler(helloHandler)
	helloWorldModule.SetControllers(helloWorldController)

	//* Set Router
	helloWorldModule.SetRoutes()
}

func (server *solanum) InitModules() {
	//* initialize Modules
	//* Modules are Singleton Designed
	// myModules := make([]modules.Module, 0)
	myModules := []modules.Module{}

	if length := len(myModules); length == 0 {
		server.appendHelloWorld()
	}

	server.modules = &myModules
}

func (server *solanum) GetModules() *[]modules.Module {
	return server.modules
}

func (server *solanum) InitGlobalMiddlewares() {
	//* 1. Logger, ...

	//* 2. Authentication, ...

	//* 3. Authorization, ...
}

func NewSolanum(port int) *Solanum {
	if SolanumRunner == nil {
		SolanumRunner = &solanum{
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
