package controller

import (
	"net/http"
)

type MainController struct {
	responseHandler ResponseHandler
}

func (c *MainController) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.responseHandler.Respond(w, r, http.StatusOK, "hello")
	}
}

func NewMainController(r ResponseHandler) *MainController {
	return &MainController{responseHandler: r}
}
