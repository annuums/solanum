package main

import (
	"database/sql"
	"fmt"
	"github.com/annuums/solanum"
	"github.com/annuums/solanum/docs/examples/dependency_injection/mvc/user"
	_ "github.com/lib/pq"
)

func RegisterDependencies() {
	solanum.Register("db", func() *sql.DB {
		dsn := "postgres://postgres:postgres@localhost:5432/annuums?sslmode=disable"
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			panic(fmt.Errorf("db open: %w", err))
		}
		return db
	}, solanum.WithSingleton())

	solanum.Register("userRepository", func(db *sql.DB) user.UserRepository {
		return &user.UserRepoImpl{DB: db}
	}, solanum.WithTransient(), solanum.As((*user.UserRepository)(nil)))
}
