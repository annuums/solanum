package solanum

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type (
	Runner interface {
		InitModules()
		GetModules() *[]Module
		InitGlobalMiddlewares()
		Run()
	}
	runner struct {
		Engine  *gin.Engine
		port    int
		modules *[]Module
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
	//* initialize Modules
	//* Modules are Singleton Designed
	// myModules := make([]Module, 0)

	helloUri := "/"
	helloWorldModule, _ := NewHelloWorldModule(
		server.Engine.Group(helloUri),
		helloUri,
	)

	myModules := []Module{
		helloWorldModule,
	}

	//* If you don't have any modules, then helloWorld will be added as a default module.
	if length := len(myModules); length == 0 {
		fmt.Println(`Appending HelloWorld... for "/"`)
		myModules = append(myModules, *server.appendHelloWorld())
	}

	server.modules = &myModules

	//* setRoutes
	for _, m := range myModules {
		m.SetRoutes()
	}
}

func (server *runner) GetModules() *[]Module {
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
