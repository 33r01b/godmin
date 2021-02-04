package controller

import (
	"encoding/json"
	"godmin/internal/model"
	"godmin/internal/server/request"
	"godmin/internal/server/response"
	"godmin/internal/store/sqlstore"
	"net/http"
)

type UserController struct {
	responseHandler response.Handler
	store           *sqlstore.Store
}

func (c *UserController) UserCreateHandle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request.UserCreate{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			c.responseHandler.Error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := req.Validate(); err != nil {
			c.responseHandler.Error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}

		if err := c.store.User().Create(u); err != nil {
			c.responseHandler.Error(w, r, http.StatusUnprocessableEntity, err)
		}

		c.responseHandler.Respond(w, r, http.StatusCreated, response.NewUserCreated(u))
	}
}

func NewUserController(r response.Handler, s *sqlstore.Store) *UserController {
	return &UserController{responseHandler: r, store: s}
}
