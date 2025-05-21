package main

import (
	"database/sql"
	"fmt"
	"github.com/annuums/solanum/container"
	"github.com/annuums/solanum/docs/examples/dependency_injection/mvc/user"
	_ "github.com/lib/pq"
)

func RegisterDependencies() {

	// Register a singleton database connection
	container.Register("db", func() *sql.DB {
		dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			panic(fmt.Errorf("db open: %w", err))
		}
		return db
	}, container.WithSingleton())

	// Register a transient user repository
	container.Register(
		"userRepository",
		func(db *sql.DB) user.UserRepository {
			return &user.UserRepoImpl{DB: db}
		},
		container.WithTransient(),
		container.As((*user.UserRepository)(nil)),
		container.WithDep[*sql.DB]("db"), // Declare *sql.DB dependency to find it in the container with key "db"
	)
}
