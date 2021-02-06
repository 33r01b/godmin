package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		Name:     "user",
		Email:    "user@example.org",
		Password: "password",
	}

}
