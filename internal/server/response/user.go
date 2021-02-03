package response

import "godmin/internal/model"

type UserCreated struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserCreated(u *model.User) *UserCreated {
	return &UserCreated{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
