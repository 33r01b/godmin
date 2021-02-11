package server

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"godmin/config"
	"godmin/internal/server/service"
	"godmin/internal/store/memorystore"
	"godmin/internal/store/sqlstore"
)

type ctxKey int8

const (
	CtxKeyUser      ctxKey = iota
	CtxKeyRequestID ctxKey = iota
)

type ServiceContainer interface {
	SqlStore() *sqlstore.Store
	MemoryStore() *memorystore.Store
	JwtService() *service.JWTService
}

type Connections struct {
	Db    *sqlx.DB
	Redis *redis.Client
}

func (c *Connections) Close() {
	if err := c.Db.Close(); err != nil {
		log.Error(fmt.Errorf("database connection close error: %w", err))
	}

	if err := c.Redis.Close(); err != nil {
		log.Error(fmt.Errorf("redis connection close error: %w", err))
	}

	log.Info("connections closed")
}

// NewConnections initialize connections
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
