package main

import (
	"database/sql"
	"fmt"
	"github.com/annuums/solanum/container"
	"github.com/annuums/solanum/docs/examples/dependency_injection/cli/user"
	_ "github.com/lib/pq"
)

func RegisterDependencies() {

	container.Register("db", func() *sql.DB {

		dsn := "postgres://postgres:postgres@localhost:5432/annuums?sslmode=disable"
		db, err := sql.Open("postgres", dsn)

		if err != nil {
			panic(fmt.Errorf("db open: %w", err))
		}

		return db
	}, container.WithSingleton())

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
