# Dependency Injection (MVC) Example

Demonstrates registering and injecting services (e.g., database and repository) in an MVC style.

## Overview
This example demonstrates how to use Solanum's dependency injection capabilities to create a simple MVC application with a PostgreSQL database.

1. RegisterDependencies (dependencies.go)

Registers a PostgreSQL *sql.DB singleton under key "db".

Registers a UserRepository factory (transient) bound to its interface.

2. Bootstrap (main.go)

Calls RegisterDependencies().

Creates a Solanum server on port 5050.

Mounts the user module.

3. User Module (user/):

model.go: User struct, UserRepository interface, Postgres implementation.

controller.go: Defines routes GET /users and POST /users.

handler.go: Handlers use solanum.GetDependency[UserRepository].

module.go: Wires controllers and DI into a SolaModule.

## Run
```bash
# If you need PostgreSQL:
docker run -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres
psql -h localhost -U postgres -c "CREATE DATABASE annuums;"

cd docs/examples/dependency_injection/mvc
go run dependencies.go main.go
```

## Endpoints:
- GET 
  - /users → List all users (JSON).

- POST
  - /users → Create a user.
    ```json
    {
      "name": "Alice",
      "email": "alice@example.com"
    }
    ```

## Test
```bash
curl http://localhost:5050/users
curl -X POST -H "Content-Type: application/json" \
     -d '{"name":"Alice","email":"alice@example.com"}' \
     http://localhost:5050/users
```
