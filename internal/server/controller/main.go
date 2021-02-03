package controller

import (
	"godmin/internal/server/response"
	"net/http"
)

type MainController struct {
	responseHandler response.Handler
}

func (c *MainController) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.responseHandler.Respond(w, r, http.StatusOK, "hello")
	}
}

func NewMainController(r response.Handler) *MainController {
	return &MainController{responseHandler: r}
}
