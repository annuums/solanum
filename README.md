# Solanum - Web Server Framework Based on Gin

- This project provides Modulability to Gin Project.
- You can implement `Module`, `Controller`, `service` to routes, handles, and intercept requests.

## Annuum, Potato Also Can Change The World!

- dev.whoan(싹난 감자) in Annuums
  - [Github](https://github.com/dev-whoan)

### Run Solanum

#### Install Go Module

```shell
$ go get github.com/annuums/solanum
```

- Fast Example

```go
package main

import "github.com/annuums/solanum"

func main() {
	server := *solanum.NewSolanum(5050)

	var helloWorldModule solanum.Module
	helloUri := "/"
	helloWorldModule, _ = solanum.NewHelloWorldModule(
		server.GetGinEngine().Group(helloUri),
		helloUri,
	)

	server.AddModule(&helloWorldModule)

	server.Run()
}
```

#### Implements Modules, Controllers, Handlers

- You should develop `Module, Controller, Handler` which are specified in `module.interface.go`. This example let you know how to implement `Module, Controller, Handler`.

##### `Module`

```go
var helloWorldModule *module

func NewHelloWorldModule(router *gin.RouterGroup, uri string) (*module, error) {
	if helloWorldModule == nil {
		helloWorldModule, _ = NewModule(router, uri)
		attachControllers()
	}

	return helloWorldModule, nil
}

func attachControllers() {
	//* Attatching Controller Directly
	var (
		ctr Controller
		err error
	)
	ctr, err = NewHelloWorldController()

	if err != nil {
		log.Fatal(err)
	}
	// ctr2, _ := NewAnotherController()
	//	...

	helloWorldModule.SetControllers(&ctr)
}
```

##### `Controller`

```go
var helloWorldController *controller

func NewHelloWorldController() (*controller, error) {
	if helloWorldController == nil {
		helloWorldController, _ = NewController()
		addHandlers()
	}

	return helloWorldController, nil
}

func addHandlers() {
	helloHandler := NewHelloWorldHandler()
	// anotherHandler := NewHelloWorldHandler()
	//* ...

	helloWorldController.AddHandler(helloHandler)
}
```

##### `Handler`

- For now, you should implement multiple handlers for multiple routings.
- It means that if you want to routes `/` and `/healthz`, should implement two `*service` for each of those.

```go
func NewHelloWorldHandler() *Service {
	return &Service{
		Uri:        "/",
		Method:     http.MethodGet,
		Handler:    indexHandler,
		Middleware: indexMiddleware,
	}
}

func indexHandler(ctx *gin.Context) {
	ctx.JSON(200, "Hello, World! From HelloWorld Index Handler.")
}

func indexMiddleware(ctx *gin.Context) {
	log.Println("Hello Index Middleware")
	ctx.Next()
}
```

- Finally, you should composite `Module, Controller, Handler`. After you composited the modules, you can add it with calling `SolanumRunner.Addmodule(module_name)`.

  - If you explicitly declare the `Module, Controller, Handler`, then you should attach `Handler` to `Controller`, and `Controller` to `Module` using functions.
  - `helloWorldController.AddHandler(helloHandler)`, `helloWorldModule.SetControllers(ctr)`

```go
func main() {
	server := *solanum.NewSolanum(5050)

	helloUri := "/"
	helloWorldModule, _ := solanum.NewHelloWorldModule(
		server.GetGinEngine().Group(helloUri),
		helloUri,
	)

	server.AddModule(helloWorldModule)

	server.Run()
}
```

- You can connect to `http://localhost:5050`. There should be a message: "Hello, World! From HelloWorld Index Handler"
