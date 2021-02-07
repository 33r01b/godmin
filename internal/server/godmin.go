package server

import (
	"github.com/gorilla/mux"
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
	Router() *mux.Router
	SqlStore() *sqlstore.Store
	MemoryStore() *memorystore.Store
	JwtService() *service.JWTService
}
