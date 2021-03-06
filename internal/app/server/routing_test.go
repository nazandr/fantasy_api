package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/stretchr/testify/assert"
)

func TestRouter_authenticateUser(t *testing.T) {
	s := TestServer(t)
	defer s.Store.DropDb()

	u := models.TestUser(t)
	s.Store.User().Create(u)

	conf := NewConfig()
	newToken := NewToken()
	newToken.Auth(u.ID, conf)
	s.Store.User().UpdateRefreshToken(u.ID, newToken.RefreshToken, conf.RefreshTokenExp)

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
		w.Header().Add("Content-Type", "application/json")
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

func TestRouter_admin(t *testing.T) {
	s := TestServer(t)
	defer s.Store.DropDb()

	u := models.TestUser(t)
	s.Store.User().Create(u)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]interface{}{
				"user_id": u.ID.Hex(),
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
			json.NewEncoder(b).Encode(tc.payload)
			req, err := http.NewRequest(http.MethodGet, "/", b)
			assert.NoError(t, err)
			s.admin(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestRouter_handleSingUp(t *testing.T) {
	s := TestServer(t)
	defer s.Store.DropDb()
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
	defer s.Store.DropDb()

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
	s.Store.User().Create(u)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, err := http.NewRequest(http.MethodPost, "/singin", b)
			assert.NoError(t, err)

			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestRouter_addCardsPacks(t *testing.T) {
	s := TestServer(t)
	u := models.TestUser(t)
	u.FantacyCoins = 10000
	defer s.Store.DropDb()

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: models.PacksCount{
				Common:  1,
				Special: 1,
			},
			expectedCode: http.StatusOK,
		},
		// {
		// 	name:         "invalid payload",
		// 	payload:      "invalid",
		// 	expectedCode: http.StatusBadRequest,
		// },
		// {
		// 	name: "invalid type",
		// 	payload: map[string]interface{}{
		// 		"common":  "invalid",
		// 		"special": 1,
		// 	},
		// 	expectedCode: http.StatusBadRequest,
		// },
		{
			name: "multiple",
			payload: map[string]interface{}{
				"acsses_token":  "token",
				"refresh_token": "token",
				"common":        1,
				"special":       1,
			},
			expectedCode: http.StatusOK,
		},
	}
	s.Store.User().Create(u)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, err := http.NewRequest(http.MethodPost, "/addCardsPack", b)
			assert.NoError(t, err)

			s.addCardsPacks().ServeHTTP(rec, req.WithContext(context.WithValue(req.Context(), cxtKeyUser, u)))
			assert.Equal(t, http.StatusText(tc.expectedCode), http.StatusText(rec.Code))
		})
	}
}
