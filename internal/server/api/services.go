package api

import (
	"godmin/config"
	"godmin/internal/server/service"
	"godmin/internal/store/memorystore"
	"godmin/internal/store/sqlstore"
)

type Services struct {
	sqlStore    *sqlstore.Store
	memoryStore *memorystore.Store
	jwtService  *service.JWTService
}

func (s *Services) SqlStore() *sqlstore.Store {
	return s.sqlStore
}

func (s *Services) MemoryStore() *memorystore.Store {
	return s.memoryStore
}

func (s *Services) JwtService() *service.JWTService {
	return s.jwtService
}

func NewServices(conn *Connections, config *config.Config) *Services {
	sqlStore := sqlstore.New(conn.Db)
	memoryStore := memorystore.New(conn.Redis)

	return &Services{
		sqlStore:    sqlStore,
		memoryStore: memoryStore,
		jwtService:  service.NewJwtService(sqlStore, memoryStore, config.Jwt),
	}
}
