package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/nazandr/fantasy_api/internal/app/store"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	cxtKeyUser cxtKey = iota
	cxtKeyRequestId
)

var (
	errIncorectEmailOrPassword = errors.New("incorect email or password")
	errExpiredRefreshToken     = errors.New("expired refresh token")
	errNoPacks                 = errors.New("user dosn`t have packs")
)

type cxtKey int

func (s *APIServer) configureRouter() {
	s.router.Use(s.setRequestId)
	s.router.Use(s.loggerReq)
	s.router.HandleFunc("/singup", s.handleSingUp()).Methods("POST")
	s.router.HandleFunc("/singin", s.handelSingIn()).Methods("POST")

	auth := s.router.PathPrefix("/auth").Subrouter()
	auth.Use(s.authenticateUser)
	auth.HandleFunc("/open-common-pack", s.openCommonPack()).Methods("POST")
	auth.HandleFunc("/collection", s.collection()).Methods("GET")
	auth.HandleFunc("/disenchant", s.disenchant())

	admin := s.router.PathPrefix("/admin").Subrouter()
	admin.Use(s.admin)
	admin.HandleFunc("/add-cards-pack", s.addCardsPacks()).Methods("POST")
}

func (s *APIServer) setRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), cxtKeyRequestId, id)))
	})
}

func (s *APIServer) loggerReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(cxtKeyRequestId),
		})

		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof("completed with %d %s at %v",
			rw.code, http.StatusText(rw.code),
			time.Since(start))
	})
}
func (s *APIServer) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		token := NewToken()
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(token); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, err := token.ParseJWT(s.config)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().Find(id)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		if u.Session.Refresh_token != token.RefreshToken {
			s.error(w, r, http.StatusUnauthorized, errExpiredRefreshToken)
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(b))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), cxtKeyUser, u)))
	})
}

func (s *APIServer) admin(next http.Handler) http.Handler {
	type request struct {
		UserId primitive.ObjectID `json:"user_id"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		req := &request{}

		if err := json.NewDecoder(bytes.NewReader(b)).Decode(req); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u, err := s.store.User().Find(req.UserId)

		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(b))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), cxtKeyUser, u)))
	})
}

func (s *APIServer) handleSingUp() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := models.NewUser()
		u.Email = req.Email
		u.Password = req.Password

		err := s.store.User().Create(u)
		if err == store.ErrUserAllreadyExist {
			s.error(w, r, http.StatusOK, err)
			return
		}
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitaze()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *APIServer) handelSingIn() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorectEmailOrPassword)
			return
		}
		if !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorectEmailOrPassword)
			return
		}
		token := NewToken()
		token.Auth(u.ID, s.config)
		if err := s.store.User().UpdateRefreshToken(u.ID, token.RefreshToken, s.config.RefreshTokenExp); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, token)
	}
}

func (s *APIServer) addCardsPacks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &models.PacksCount{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := r.Context().Value(cxtKeyUser).(*models.User)

		u.Packs.Common += req.Common
		u.Packs.Special += req.Special

		if err := s.store.User().ReplaseUser(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *APIServer) openCommonPack() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(cxtKeyUser).(*models.User)
		if u.Packs.Common <= 0 {
			s.error(w, r, http.StatusBadRequest, errNoPacks)
			return
		}

		p, err := s.store.PlayerCards().OpenCommonPack(s.store)
		u.Packs.Common -= 1
		u.CardsCollection = append(u.CardsCollection, p.Cards[:]...)
		s.store.User().ReplaseUser(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, p)
	}
}

func (s *APIServer) collection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(cxtKeyUser).(*models.User)
		s.respond(w, r, http.StatusOK, u.CardsCollection)
	}
}

func (s *APIServer) disenchant() http.HandlerFunc {
	type card struct {
		ID primitive.ObjectID `json:"card_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &card{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := r.Context().Value(cxtKeyUser).(*models.User)

		for i, v := range u.CardsCollection {
			if v.Id == req.ID {
				u.CardsCollection = removeCard(u.CardsCollection, i)
				break
			}
		}
		u.FantacyCoins += 200
		if err := s.store.User().ReplaseUser(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, u.CardsCollection)
	}
}

func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func removeCard(s []models.PlayerCard, i int) []models.PlayerCard {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
