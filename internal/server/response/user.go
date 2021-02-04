package response

import "godmin/internal/model"

type User struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUser(u *model.User) *User {
	return &User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
