# Solanum - Web Server Framework Based on Gin

- This project provides Modulability to Gin Project.
- You can implement `Module`, `Controller`, `service` to routes, handles, and intercept requests.

## Annuum, Potato Also Can Change The World!

- dev.whoan(싹난 감자) in Annuums
  - [Github](https://github.com/dev-whoan)

### Run Example Server

```shell
$ cd example
$ go run .
# OR
$ cd example
$ go build
$ ./example
```

### Run Solanum

- You should implement `Runner` interface, and `runner` struct in `solanum.go`.

  - Must override methods: `InitModules(), InitGlobalMiddlewares()`

- You should develop `Modules, Controllers, Handlers` which are written in `module.interface.go`

#### Examples

- This example let you know how to implement `Module, Controller, Handler`.

- `Module`

```go
var helloWorldModule Module

//* Creating New Module
//* You can just initialize Module
func NewHelloWorldModule(router *gin.RouterGroup, uri string) (Module, error) {
	if helloWorldModule == nil {
		helloWorldModule, _ = NewModule(router, uri)
	}

	return helloWorldModule, nil
}

//* Or add controller directly here
func NewHelloWorldModule(router *gin.RouterGroup, uri string) (Module, error) {
	if helloWorldModule == nil {
		helloWorldModule, _ = NewModule(router, uri)
	}

  ctr, _ := NewHelloWorldController()
  ctr2, _ := NewAnotherController()
  ...

  helloWorldModule.SetControllers(ctr, ctr2, ...)

	return helloWorldModule, nil
}
```

- `Controller`

```go
var helloWorldController Controller

//* Creating New Controller
//* You can just initialize Controller
func NewHelloWorldController() (Controller, error) {
	if helloWorldController == nil {
		helloWorldController, _ = NewController()
	}

	return helloWorldController, nil
}

//* Or add handlers directly here

func NewHelloWorldController() (Controller, error) {
	if helloWorldController == nil {
		helloWorldController, _ = NewController()
	}

  helloHandler := NewHelloWorldHandler()
 	anotherHandler := NewHelloWorldHandler()
  ...

	helloWorldController.AddHandler(helloHandler, anotherHandler, ...)

	return helloWorldController, nil
}
```

- `Handler`
  - For now, you should implement multiple handlers for multiple routings.
  - It means that if you want to routes `/` and `/healthz`, should implement two `*service` for each of those.

```go
func NewHelloWorldHandler() *Service {
	return &Service{
		uri:        "/",
		method:     http.MethodGet,
		handler:    indexHandler,
		middleware: indexMiddleware,
	}
}

func indexHandler(ctx *gin.Context) {
	ctx.JSON(200, "Hello, World! From HelloWorld Index Handler. Greeting!")
}

func indexMiddleware(ctx *gin.Context) {
	log.Println("Hello Index Middleware")
	ctx.Next()
}
```

- Finally, you should integrate `Module, Controller, Handler`. It is done with overriding `func (server *server) InitModules()` which defined in `solanum.go`

  - If you explicitly declare the `Module, Controller, Handler`, then you should attach `Handler` to `Controller`, and `Controller` to `Module` using functions.
  - `helloWorldController.AddHandler(helloHandler)`, `helloWorldModule.SetControllers(ctr)`

```go
func (server *runner) InitModules() {
	//* Initialize Modules
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

	//* Setting Routes for Modules
	for _, m := range myModules {
		m.SetRoutes()
	}
}
```

- You can connect to `http://localhost:5050`. There should be a message: "Hello, World! From HelloWorld Index Handler"
