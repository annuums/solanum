# Solanum — A Modular Gin-Based Framework with Built-In Dependency Injection

**Solanum** is designed to help Go developers ship clean, maintainable HTTP services and micro-APIs by combining:

- **A clear module/controller structure**  
- **First-class dependency injection** (no global state, minimal reflection)  
- **Pluggable middleware** (CORS, logging, auth, etc.)  
- **Testable, decoupled components**

---

## Why Solanum?

When you build a typical Gin application, it’s easy to end up with:

- One giant `main.go` wiring routes, middleware, and handlers inline  
- Hard-to-test handlers that call global singletons or package-level variables  
- Tight coupling between controllers, services, and data stores  
- Boilerplate to register dependencies and then invoke them manually  

Solanum solves these pain points by:

1. **Inversion of Control & DI container**  
   Register your services, repositories, and clients once. Solanum will resolve and inject them for you.  
2. **Modular structure**  
   Group routes, middleware, and handlers into self-contained `Module`s. Each module declares its own dependencies.  
3. **Minimal reflection, maximal clarity**  
   Define constructor functions with explicit parameters; Solanum uses simple reflection only to wire them up.  
4. **Test-driven development**  
   Because your handlers request only the dependencies they need, you can swap in mocks or stubs without spinning up a server.

---

## What Can You Build?

- **RESTful APIs**  
- **Microservices**  
- **GraphQL or gRPC gateways**  
- **CLI tools or batch workers** (reuse DI container in non-web context)  
- **Web Server for Applications**; Health checks, admin dashboards, metrics endpoints, etc.

Whether you need a single `/ping` health check or a complex set of user, order, and payment services, Solanum keeps your code organized and your dependencies explicit.

---

## Key Features

### 1. Modular Architecture  
```go
// Define a /users module
userModule := solanum.NewModule("/users")
userModule.SetControllers(myUserController)
userModule.SetDependencies(solanum.Dep[UserService]("userSvc"))
```

### 2. Controller & Service Layer
```go
type SolaService struct {
    Uri     string
    Method  string
    Handler gin.HandlerFunc
}
```

### 3. Dependency Injection Container
```go
solanum.Register("db", ProvideDB, solanum.WithSingleton())
solanum.Register(
  "userSvc",
  ProvideUserService,
  solanum.WithTransient(),
  solanum.As((*UserService)(nil)),
)
```

### 4. Flexible CORS & Middleware
```go
app.Cors(
  solanum.WithOrigins("https://example.com"),
  solanum.WithMethods("GET", "POST"),
  solanum.WithAllowCredentials(true),
)
```

---

## Getting Started
### Get Solanum
```bash
go get github.com/annuums/solanum
```

### Easy Start
```go
package main

import "github.com/annuums/solanum"

func main() {
    // 1) Register dependencies
    solanum.Register("db", ProvideDB, solanum.WithSingleton())
    solanum.Register("userRepo", ProvideUserRepo, solanum.WithTransient(), solanum.As((*UserRepo)(nil)))
    solanum.Register("userSvc", ProvideUserService, solanum.WithTransient())

    // 2) Define module
    userModule := solanum.NewModule("/users")
    userModule.SetControllers(userController)
    userModule.SetDependencies(
        solanum.Dep[*sql.DB]("db"),
        solanum.Dep[UserRepo]("userRepo"),
        solanum.Dep[*UserService]("userSvc"),
    )

    // 3) Create & run server
    app := solanum.NewSolanum(8080)
    app.SetModules(userModule)
    app.Cors()   // default CORS
    app.Run()    // starts on :8080
}
```

### Examples & More
👉 Learn by [Examples](./docs/examples/README.md)

## Who Is This For?
- Go backend teams building REST or microservices
- Developers who want clear separation between controllers, services, and data layers
- Projects that require testable, maintainable code without manual wiring in every main.go
- Engineers seeking a lightweight alternative to heavier frameworks with reflection-heavy DI
