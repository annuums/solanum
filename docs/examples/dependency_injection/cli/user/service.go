package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type UserService struct {
	repo UserRepository
	db   *sql.DB
}

func NewUserService(repo UserRepository, db *sql.DB) *UserService {
	return &UserService{repo: repo, db: db}
}

func (svc *UserService) ListUsersCLI() error {

	if svc.repo == nil {
		return fmt.Errorf("userRepository is not found in context")
	}

	users, err := svc.repo.FindAll()
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

func (svc *UserService) AddUserCLI(name, email string) error {

	newUser := &User{
		Name:  name,
		Email: email,
	}

	if err := svc.repo.Create(newUser); err != nil {
		return err
	}

	out, err := json.MarshalIndent(newUser, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
