package middleware

import (
	"context"
	"godmin/internal/server"
	"godmin/internal/server/response"
	"godmin/internal/server/service"
	"net/http"
)

type JwtAuth struct {
	jwtService      *service.JWTService
	responseHandler response.Handler
}

func (j *JwtAuth) JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := j.jwtService.Authenticate(r)
		if err != nil {
			j.responseHandler.Error(w, r, err.GetStatusCode(), err.GetError())
			return
		}

		user := response.NewUser(u)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), server.CtxKeyUser, user)))
	})
}

func NewJwtAuth(jwtService *service.JWTService, responseHandler response.Handler) *JwtAuth {
	return &JwtAuth{
		jwtService:      jwtService,
		responseHandler: responseHandler,
	}
}
