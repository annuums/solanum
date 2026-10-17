# Simple Example

A minimal server that mounts a health-check endpoint using the built-in module.

```go
package main

import (
	"log"
	"net/http"

	"github.com/annuums/solanum"
	"github.com/gin-gonic/gin"
)

func main() {

	pingModule := solanum.NewModule(
		solanum.WithUri("/ping"),
	)

	ctrl := solanum.NewController()
	ctrl.SetHandlers(
		&solanum.SolaService{
			Uri:    "",
			Method: http.MethodGet,
			Handler: func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "pong")
			},
		},
	)
	pingModule.AddControllers(ctrl)

	server := solanum.NewSolanum(
		solanum.WithPort(5050),
		solanum.WithMiddlewares(gin.Logger()),
	)

	server.AddModules(pingModule)
	if err := server.Run(); err != nil {

		log.Fatal(err)
	}
}
```

## Run
```bash
cd docs/examples/simple
go run main.go
# then in another terminal:
curl http://localhost:5050/ping
```

Press `Ctrl+C` to stop: the server shuts down gracefully, letting in-flight requests finish.
