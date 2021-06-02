package api_server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestRouter_authenticateUser(t *testing.T) {
	s := TestServer(t)
	defer s.store.DropDb()

	u := models.TestUser(t)
	s.store.User().Create(u)

	conf := NewConfig()
	newToken := NewToken()
	newToken.Auth(u.ID, conf)
	s.store.User().UpdateRefreshToken(u.ID, newToken.RefreshToken, conf.RefreshTokenExp)

	testCases := []struct {
		name         string
		tokenValue   map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			tokenValue: map[interface{}]interface{}{
				"token": newToken,
			},
			expectedCode: http.StatusOK,
		},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.tokenValue["token"])
			req, err := http.NewRequest(http.MethodGet, "/", b)
			assert.NoError(t, err)
			s.authenticateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestRouter_handleSingUp(t *testing.T) {
	s := TestServer(t)
	defer s.store.DropDb()
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "user allready exist",
			payload: map[string]string{
				"email":    "user@example.org",
				"password": "password",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email":    "invalid",
				"password": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, err := http.NewRequest(http.MethodPost, "/singup", b)
			if err != nil {
				t.Fatal(err)
			}
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestRouter_handleSingIn(t *testing.T) {
	s := TestServer(t)
	u := models.TestUser(t)
	defer s.store.DropDb()

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalid payload",
			payload:      "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": u.Password,
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "invalid",
			},
			expectedCode: http.StatusUnauthorized,
		},
	}
	s.store.User().Create(u)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, err := http.NewRequest(http.MethodPost, "/singin", b)
			if err != nil {
				assert.NoError(t, err)
			}
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
