package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"godmin/internal/model"
)

type User struct {
	db *sqlx.DB
}

func (ur *User) Create(u *model.User) error {
	emailExists, err := ur.EmailExists(u)
	if err != nil {
		return err
	}

	if emailExists {
		return fmt.Errorf("already used email: %s", u.Email)
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return ur.db.QueryRow(
		"INSERT INTO users (name, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id",
		u.Name,
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (ur *User) EmailExists(u *model.User) (bool, error) {
	var count int

	err := ur.db.QueryRow(
		"SELECT count(1) FROM users WHERE email = $1",
		u.Email,
	).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func NewUser(db *sqlx.DB) *User {
	return &User{
		db: db,
	}
}
