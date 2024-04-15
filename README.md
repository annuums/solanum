# Solanum - Web Server Framework Based on Gin

- This project provides Modulability to Gin Project.
- You can implement `Module`, `Controller`, `service` to routes, handles, and intercept requests.

## Annuum, Potato Can Change The World!

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
  
  var healthCheckModule solanum.Module
  healthCheckUri := "/ping"
  healthCheckModule = solanum.NewHealthCheckModule(
    healthCheckUri,
  )
  
  server.SetModules(healthCheckModule)
  
  server.Run()
}
```

#### Implements Modules, Controllers, Handlers

- You should develop `Module, Controller, Handler` which are specified in `module.interface.go`. This example let you know how to implement `Module, Controller, Handler`.

##### `Module`

```go
var hzOnce sync.Once

func NewHealthCheckModule(uri string) *SolaModule {
  hzOnce.Do(func() {
    if helathCheckModule == nil {
      helathCheckModule = NewModule(uri)
      attachControllers()
    }
  })
  
  return helathCheckModule
}

func attachControllers() {
  //* Attatching Controller Directly
  ctr := NewHealthCheckController()
  // ctr2, _ := NewAnotherController()
  //	...
  
  helathCheckModule.SetControllers(ctr)
}

```

##### `Controller`

```go
func NewHealthCheckController() *SolaController {
  healthCheckController := NewController()
  
  healthCheckController.SetHandlers(SolaService{
    Uri:        "",
    Method:     http.MethodGet,
    Handler:    hzHandler,
    Middleware: hzMiddleware,
  })
  
  return healthCheckController
}

func hzMiddleware(ctx *gin.Context) {
  log.Println("Health Checking...")
  ctx.Next()
}

```

##### `Handler`

- For now, you should implement multiple handlers for multiple routings.
- It means that if you want to routes `/` and `/healthz`, should implement two `*service` for each of those.

```go
func hzHandler(ctx *gin.Context) {
    ctx.String(200, "pong")
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
		helloUri,
	)

	server.AddModule(helloWorldModule)

	server.Run()
}
```

- You can connect to `http://localhost:5050`. There should be a message: "Hello, World! From HelloWorld Index Handler"
