package server

import (
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/store/memorystore"
	"godmin/internal/store/sqlstore"
	"net/http"
)

func Run(context context.Context, config *config.Config) error {
	log.SetLevel(config.LogLevel)

	db, err := sqlstore.NewDB(config.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	memoryStore, err := memorystore.New(config.RedisUrl)
	if err != nil {
		return err
	}

	server := newServer(sqlstore.New(db), memoryStore)

	return http.ListenAndServe(config.BindAddr, server)
}

func newServer(store *sqlstore.Store, memoryStore *memorystore.Store) *Server {
	s := &Server{
		router:      mux.NewRouter(),
		store:       store,
		memoryStore: memoryStore,
	}

	s.configureRouter()

	return s
}

type Server struct {
	router      *mux.Router
	store       *sqlstore.Store
	memoryStore *memorystore.Store
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
