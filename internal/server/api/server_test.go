package api

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"godmin/config"
	"godmin/internal/model"
	"godmin/internal/server"
	"godmin/internal/server/request"
	"godmin/internal/server/response"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	config.BootstrapTest()
	os.Exit(m.Run())
}

func TestServer_Login(t *testing.T) {
	conf := config.NewConfig()

	conn, err := server.NewConnections(conf)
	if err != nil {
		log.Fatal(err)
	}

	services := NewServices(conn, conf)

	u := model.TestUser(t)
	if err := services.SqlStore().User().Create(u); err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		if err := services.SqlStore().User().Delete(u); err != nil {
			log.Fatal(err)
		}
		conn.Redis.FlushAll()

		conn.Close()
	})

	r := request.Login{
		Email:    u.Email,
		Password: u.Password,
	}

	testCases := []struct {
		name         string
		user         func() request.Login
		expectedCode int
		testBody     func(rec *httptest.ResponseRecorder)
	}{
		{
			name: "logged",
			user: func() request.Login {
				return r
			},
			expectedCode: http.StatusOK,
			testBody: func(rec *httptest.ResponseRecorder) {
				token := &response.Token{}
				if err := json.NewDecoder(rec.Body).Decode(token); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, token)
			},
		},
		{
			name: "wrong password",
			user: func() request.Login {
				r.Password = "wrong_password"
				return r
			},
			expectedCode: http.StatusUnauthorized,
			testBody: func(rec *httptest.ResponseRecorder) {
				body := make(map[string]string)
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				assert.NotEmpty(t, body)
				assert.Equal(t, "incorrect email or password", body["error"])
			},
		},
	}

	api := NewApi(conf, NewServices(conn, conf))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}

			if err := json.NewEncoder(b).Encode(tc.user()); err != nil {
				t.Fatal(err)
			}

			req, _ := http.NewRequest(http.MethodPost, "/login", b)
			api.server.Handler.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
			tc.testBody(rec)
		})
	}
}
