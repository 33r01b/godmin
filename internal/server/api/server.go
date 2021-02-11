package api

import (
	"context"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/server/router"
	"net/http"
	"strconv"
	"time"
)

type Api struct {
	server *http.Server
	errors chan error
}

func (a *Api) Run() {
	go func() {
		log.Info("run api server")
		a.errors <- a.server.ListenAndServe()
		close(a.errors)
	}()
}

func (a *Api) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	log.Info("api server stopped")

	return nil
}

// Notify returns a channel to notify the caller about errors.
// If you receive an error from the channel you should stop the application.
func (a *Api) Notify() <-chan error {
	return a.errors
}

func NewApi(config *config.Config, services *Services) *Api {
	return &Api{
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(int(config.Port)),
			Handler: router.NewRouter(services),
		},
		errors: make(chan error, 1),
	}
}
