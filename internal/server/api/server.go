package api

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"godmin/config"
	"godmin/internal/server/router"
	"godmin/internal/store/memorystore"
	"godmin/internal/store/sqlstore"
	"net/http"
)

func Run(context context.Context, config *config.Config) error {
	conn, err := NewConnections(config)
	if err != nil {
		return err
	}
	defer conn.Close()

	return http.ListenAndServe(config.BindAddr, NewServer(NewServices(conn, config)))
}

func NewConnections(config *config.Config) (*Connections, error) {
	db, err := sqlstore.NewDB(config.Database)
	if err != nil {
		return nil, err
	}

	memory, err := memorystore.NewClient(config.RedisUrl)
	if err != nil {
		return nil, err
	}

	return &Connections{
		Db:    db,
		Redis: memory,
	}, nil
}

type Connections struct {
	Db    *sqlx.DB
	Redis *redis.Client
}

func (c *Connections) Close() error {
	if err := c.Db.Close(); err != nil {
		return err
	}

	if err := c.Redis.Close(); err != nil {
		return err
	}

	return nil
}

type Server struct {
	router *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer(services *Services) *Server {
	return &Server{
		router: router.NewRouter(services),
	}
}
