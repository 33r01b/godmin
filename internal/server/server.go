package server

import (
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/store"
	"net/http"
)

func Run(context context.Context, config *config.Config) error {
	log.SetLevel(config.LogLevel)

	db, err := store.NewDB(&config.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	return http.ListenAndServe(config.BindAddr, newServer(store.New(db)))
}

func newServer(store *store.Store) *Server {
	s := &Server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

type Server struct {
	router *mux.Router
	store  *store.Store
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
