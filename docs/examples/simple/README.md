# Simple Example

A minimal server that mounts a health-check endpoint using the built-in module.

```go
package main

import "github.com/annuums/solanum"

func main() {
    server := *solanum.NewSolanum(5050)

    // Mount the built-in HealthCheck module at /ping
    healthCheck := solanum.NewHealthCheckModule("/ping")
    server.SetModules(healthCheck)

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
