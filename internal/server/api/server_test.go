package api

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"godmin/config"
	"godmin/internal/model"
	"godmin/internal/server/request"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("GODMIN_ENV", "test"); err != nil {
		os.Exit(1)
	}

	config.Bootstrap()
	os.Exit(m.Run())
}

func TestServer_Login(t *testing.T) {
	conf := config.NewConfig()

	conn, err := NewConnections(conf)
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer(conn, conf)

	u := model.TestUser(t)
	if err := server.SqlStore().User().Create(u); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		if err := server.SqlStore().User().Delete(u); err != nil {
			log.Fatal(err)
		}
		conn.Redis.FlushAll()

		defer conn.Close()
	})

	r := request.Login{
		Email:    u.Email,
		Password: u.Password,
	}

	testCases := []struct {
		name         string
		user         func() request.Login
		expectedCode int
	}{
		{
			name: "logged",
			user: func() request.Login {
				return r
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "wrong password",
			user: func() request.Login {
				r.Password = "wrong_password"
				return r
			},
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}

			if err := json.NewEncoder(b).Encode(tc.user()); err != nil {
				t.Fatal(err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/login", b)
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
