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

import (
	"github.com/annuums/solanum"
)

func main() {
	server := *solanum.NewSolanum(5050)

	helloUri := "/"
	helloWorldModule, _ := solanum.NewHelloWorldModule(
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
var helloWorldModule Module

//* Creating New Module
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

##### `Controller`

```go
var helloWorldController Controller

//* Creating New Controller

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

##### `Handler`

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

- Finally, you should composite `Module, Controller, Handler`. After you composited the modules, you can add it with calling `SolanumRunner.Addmodule(module_name)`.

  - If you explicitly declare the `Module, Controller, Handler`, then you should attach `Handler` to `Controller`, and `Controller` to `Module` using functions.
  - `helloWorldController.AddHandler(helloHandler)`, `helloWorldModule.SetControllers(ctr)`

```go
func main() {
  server := *solanum.NewSolanum(5050)

	helloUri := "/"
	helloWorldModule, _ := NewHelloWorldModule(
		server.GetGinEngine().Group(helloUri),
		helloUri,
	)

	server.AddModule(&helloWorldModule)
}
```

- You can connect to `http://localhost:5050`. There should be a message: "Hello, World! From HelloWorld Index Handler"
