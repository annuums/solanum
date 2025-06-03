package user

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(u *User) error
	FindAll() ([]User, error)
}

type UserRepoImpl struct {
	DB *sql.DB
}

func (r *UserRepoImpl) Create(u *User) error {
	return r.DB.QueryRow(
		"INSERT INTO users(name,email,created_at) VALUES($1,$2,$3) RETURNING id",
		u.Name, u.Email, u.CreatedAt,
	).Scan(&u.ID)
}

func (r *UserRepoImpl) FindAll() ([]User, error) {
	rows, err := r.DB.Query("SELECT id,name,email,created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}
