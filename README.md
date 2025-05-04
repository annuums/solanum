# Solanum - Web Server Framework Based on Gin

**Solanum** is a lightweight, Gin-based web framework for Go that brings a clean modular structure and built-in dependency injection (DI). It makes it easy to organize your routes, middleware, and service dependencies in a scalable, testable way.

## Annuums, Potato Can Change The World!

- Eugene
  - [Github](https://github.com/dev-whoan)


## Features

- **Modular Architecture**  
  Define self-contained `Module`s, each with its own URI prefix, middleware, controllers, and dependencies.

- **Controller & Service Layer**  
  Organize handlers as `SolaService` entries within `Controller`s, enabling clear separation of concerns.

- **Dependency Injection Container**  
  Register providers (factory functions or concrete instances) with singleton or transient scope, support init hooks, and resolve by key or Go interface type.

- **Flexible CORS Support**  
  Fluent API for configuring origins, methods, headers, credentials, and max age.

## Installation

```bash
go get github.com/annuums/solanum
```

## Quick Start
```go
package main

import "github.com/annuums/solanum"

func main() {
	server := *solanum.NewSolanum(5050)

	healthCheckUri := "/ping"
	healthCheckModule := solanum.NewHealthCheckModule(
		healthCheckUri,
	)

	server.SetModules(healthCheckModule)

	server.Run()
}
```

### 👉 [Examples](./docs/examples/README.md)
