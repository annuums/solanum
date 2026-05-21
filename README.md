# Solanum — A Modular Gin-Based Web Framework

**Solanum** is a lightweight framework built on top of Gin that helps Go developers
build clean, maintainable HTTP services and micro-APIs with **explicit structure and flow**.

Instead of hiding control flow behind containers or implicit wiring,
Solanum encourages **clear module boundaries, predictable routing, and readable composition**.

---

## Why Solanum?

A typical Gin application often grows into:

- A massive `main.go` that wires routes, middleware, and handlers inline
- Implicit control flow that’s hard to follow once the project grows
- Business logic spread across handlers with no clear module boundaries
- Boilerplate duplicated across services and endpoints

Solanum addresses these problems by focusing on:

### 1. Explicit Flow Over Magic
All routes, handlers, and middleware are registered explicitly.
There is no hidden resolution or automatic injection —  
**what you write is exactly what runs**.

### 2. Modular Structure
Routes are grouped into self-contained `Module`s.
Each module owns its routing prefix, controllers, and middleware.

### 3. Predictable Composition
Applications are composed step by step:
modules → controllers → handlers → server.
You can trace request flow by reading the code top to bottom.

### 4. Test-Friendly Design
Because handlers and services are constructed explicitly,
you can test them in isolation without bootstrapping the entire server.

---

## What Can You Build?

- **RESTful APIs**
- **Microservices**
- **GraphQL or gRPC gateways**
- **Admin or internal tools**
- **Web servers for applications**
    - Health checks
    - Metrics
    - Admin dashboards

Solanum scales from a single `/ping` endpoint
to multi-module services without changing how you structure your code.

---

## Core Principles

- **Explicit is better than implicit**
- **Composition over configuration**
- **No hidden lifecycle or control flow**
- **Readable code beats clever abstractions**

Solanum does not try to replace Gin —
it provides a **clear structure around it**.

---

## Key Features

### 1. Modular Architecture
```go
// Define a /users module
userModule := solanum.NewModule(
    solanum.WithUri("/users"),
)
userModule.SetControllers(myUserController)
```

Each module:
- Owns its route prefix
- Registers its controllers explicitly
- Can be reasoned about independently


### 2. Controller-Centered Routing
```go
type SolaService struct {
    Uri     string
    Method  string
    Handler gin.HandlerFunc
}
```
Controllers define:
- HTTP method
- Relative URI
- Handler function

No reflection-based routing, no annotations — just plain Go structures.


### 3. Explicit Middleware Configuration
```go
import "github.com/annuums/solanum/middleware/cors"

server.Cors(
  cors.WithUrls([]string{"https://example.com"}),
  cors.WithMethods([]string{"GET", "POST"}),
  cors.WithAllowCredentials(true),
)
```
Middleware is applied explicitly at the application level, making request behavior easy to reason about.

---

## Getting Started
### Get Solanum
```bash
go get github.com/annuums/solanum
```

### Easy Start
```go
package main

import (
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

Every step is explicit:
- Module creation
- Controller registration
- Handler definition
- Server composition

No hidden wiring. No surprises.

### Examples & More
👉 Learn by [Examples](./docs/examples/README.md)

## Who Is This For?
- Go developers who value explicit control flow
- Teams building REST APIs or microservices with Gin
- Pjects that want structure without heavy frameworks
- ineers who prefer readable composition over implicit behavior

If you believe that clarity beats convenience, Solanum is for you.
