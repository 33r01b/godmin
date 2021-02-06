package controller

import (
	"encoding/json"
	"godmin/internal/server/request"
	"godmin/internal/server/response"
	"godmin/internal/server/service"
	"net/http"
)

type AuthController struct {
	jwtService      *service.JWTService
	responseHandler response.Handler
}

func (c *AuthController) HandleLogin() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		login := &request.Login{}
		if err := json.NewDecoder(r.Body).Decode(login); err != nil {
			c.responseHandler.Error(w, r, http.StatusBadRequest, err)
			return
		}

		token, err := c.jwtService.CreateToken(login)
		if err != nil {
			c.responseHandler.Error(w, r, err.GetStatusCode(), err.GetError())
			return
		}

		c.responseHandler.Respond(w, r, http.StatusOK, token)
	}
}

func (c *AuthController) HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := c.jwtService.Logout(r); err != nil {
			c.responseHandler.Error(w, r, err.GetStatusCode(), err.GetError())
			return
		}

		c.responseHandler.Respond(w, r, http.StatusOK, "Successfully logged out")
	}
}

func NewAuthController(jwtService *service.JWTService, responseHandler response.Handler) *AuthController {
	return &AuthController{
		jwtService:      jwtService,
		responseHandler: responseHandler,
	}
}
