package router

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"godmin/internal/server"
	"godmin/internal/server/controller"
	"godmin/internal/server/middleware"
	"godmin/internal/server/response"
	"net/http"
	"time"
)

func Configure(s server.ServiceContainer) {
	s.Router().Use(setRequestID)
	s.Router().Use(logRequest)

	responseHandler := response.NewResponse()

	// main
	mainController := controller.NewMainController(responseHandler)
	s.Router().HandleFunc("/", mainController.Handle()).Methods(http.MethodGet)

	// users
	userController := controller.NewUserController(responseHandler, s.SqlStore())
	user := s.Router().PathPrefix("/users").Subrouter()
	user.HandleFunc("/", userController.UserCreateHandle()).Methods(http.MethodPost)

	// login
	authController := controller.NewAuthController(s.JwtService(), responseHandler)
	s.Router().HandleFunc("/login", authController.HandleLogin()).Methods(http.MethodPost)

	// admin
	admin := s.Router().PathPrefix("/admin").Subrouter()
	jwtAuthMiddleware := middleware.NewJwtAuth(s.JwtService(), responseHandler)
	admin.Use(jwtAuthMiddleware.JwtAuthentication)
	admin.HandleFunc("/logout", authController.HandleLogout()).Methods(http.MethodGet)
	admin.HandleFunc("/whoami", userController.HandleWhoami()).Methods(http.MethodGet)
}

func setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), server.CtxKeyRequestID, id)))
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithFields(log.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(server.CtxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &response.Writer{ResponseWriter: w, Code: http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.Code,
			http.StatusText(rw.Code),
			time.Now().Sub(start),
		)
	})
}
