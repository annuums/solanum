# Simple Example

A minimal server that mounts a health-check endpoint using the built-in module.

```go
package main

import "github.com/annuums/solanum"

func main() {

	pingModule := solanum.NewModule(
		solanum.WithUri("/ping"),
	)

	ctrl := solanum.NewController()
	ctrl.SetHandlers(
		&solanum.SolaService{
			Uri:    "",
			Method: http.MethodGet,
			Handler: func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			},
		},
	)
	pingModule.SetControllers(ctrl)

	server := solanum.NewSolanum(
		solanum.WithPort(5050),
	)

	server.SetModules(pingModule)
	server.Run()
}
```

## Run
```bash
cd docs/examples/simple
go run main.go
# then in another terminal:
curl http://localhost:5050/ping
```
