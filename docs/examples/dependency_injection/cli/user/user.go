package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/annuums/solanum/container"
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

// TODO ctx 어떻게 처리할건지..
// TODO Module에서 등록해버리기 때문에,, container.Resolve가 아니라 Module에서 획득할 수 있지 않을까...
// 그런데 Module에 등록된 의존성에서 Module을 참조할 수는 없으니...
// TODO 이는 CLI, Gin 둘다 공통인듯... (ctx 말고 module 획득 방식 고민해보기)
func ListUsersCLI() error {

	repo := container.Dep[UserRepository]("userRepository")

	//repo := container.DepFromContext[UserRepository](ctx, "userRepository")
	if repo == nil {
		return fmt.Errorf("userRepository is not found in context")
	}

	users, err := repo.FindAll()
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func AddUserCLI(name, email string) error {

	repo := container.Dep[UserRepository]("userRepository")
	//repo := container.DepFromContext[UserRepository](ctx, "userRepository")
	if repo == nil {
		return fmt.Errorf("userRepository is not found in context")
	}

	newUser := &User{
		Name:  name,
		Email: email,
	}

	if err := repo.Create(newUser); err != nil {
		return err
	}

	out, err := json.MarshalIndent(newUser, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
