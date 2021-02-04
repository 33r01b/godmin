package api

import (
	"context"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/server/service"
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

	store := sqlstore.New(db)
	memoryStore, err := memorystore.New(config.RedisUrl)
	if err != nil {
		return err
	}

	jwtService := service.NewJwtService(store, memoryStore, config.Jwt)

	server := newServer(store, memoryStore, jwtService)

	return http.ListenAndServe(config.BindAddr, server)
}

type Server struct {
	router      *mux.Router
	store       *sqlstore.Store
	memoryStore *memorystore.Store
	jwtService  *service.JWTService
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store *sqlstore.Store, memoryStore *memorystore.Store, jwtService *service.JWTService) *Server {
	s := &Server{
		router:      mux.NewRouter(),
		store:       store,
		memoryStore: memoryStore,
		jwtService:  jwtService,
	}

	s.configureRouter()

	return s
}
