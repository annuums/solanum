# Examples for Solanum
This directory contains two self-contained example applications demonstrating core Solanum features.

```
docs/examples
├── README.md                 ← this file
├── simple
│   └── main.go               ← minimal “health check” server
└── dependency_injection
    ├── README.md             ← DI/MVC example details
    └── mvc
        ├── dependencies.go   ← register DI providers
        ├── main.go           ← application bootstrap
        └── user              ← user domain module
            ├── controller.go
            ├── handler.go
            ├── model.go
            └── module.go
```

## Simple Examples
- [Simple Healthcheck Example](simple/README.md): A minimal example demonstrating how to create a simple health check server using Solanum.

## Dependency Injection Examples
- [Postgres DI MVC Example](dependency_injection/mvc/README.md): A more complex example demonstrating how to use Solanum's dependency injection capabilities to create a simple MVC application with a PostgreSQL database.
