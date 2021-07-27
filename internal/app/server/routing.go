package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	s.router.HandleFunc("/verify", s.verify()).Methods("GET")

	auth := s.router.PathPrefix("/auth").Subrouter()
	auth.Use(s.authenticateUser)
	auth.HandleFunc("/openCommonPack", s.openCommonPack()).Methods("POST")
	auth.HandleFunc("/collection", s.collection()).Methods("GET")
	auth.HandleFunc("/user", s.userData()).Methods("GET")
	auth.HandleFunc("/disenchant", s.disenchant()).Methods("POST")
	auth.HandleFunc("/setFantasyTeam", s.setFantacyTeam()).Methods("POST")
	auth.HandleFunc("fantacyTeamsCollection", s.fantacyTeamsCollection()).Methods("GET")
	auth.HandleFunc("/addCardsPack", s.addCardsPacks()).Methods("POST")

	admin := s.router.PathPrefix("/admin").Subrouter()
	admin.Use(s.admin)
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
		logger := s.Logger.WithFields(logrus.Fields{
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

func (s *APIServer) verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := NewToken()
		if err := json.NewDecoder(r.Body).Decode(token); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, err := token.ParseJWT(s.config)
		if err != nil {
			u, err := s.Store.User().Find(id)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			fmt.Println(u.Session.Refresh_token)
			fmt.Println(token.RefreshToken)
			if u.Session.Refresh_token == token.RefreshToken {
				token := NewToken()
				token.Auth(u.ID, s.config)
				if err := s.Store.User().UpdateRefreshToken(u.ID, token.RefreshToken, s.config.RefreshTokenExp); err != nil {
					s.error(w, r, http.StatusInternalServerError, err)
					return
				}
				s.respond(w, r, http.StatusOK, token)
				return
			} else {
				s.error(w, r, http.StatusUnauthorized, errExpiredRefreshToken)
				return
			}
		}
		u, err := s.Store.User().Find(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		token.Auth(u.ID, s.config)
		if err := s.Store.User().UpdateRefreshToken(u.ID, token.RefreshToken, s.config.RefreshTokenExp); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, token)
	}
}

func (s *APIServer) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(b))
		token := NewToken()
		if err := json.NewDecoder(bytes.NewBuffer(b)).Decode(token); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, err := token.ParseJWT(s.config)
		if err != nil {
			u, err := s.Store.User().Find(id)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if u.Session.Refresh_token == token.RefreshToken {

				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), cxtKeyUser, u)))
				return
			} else {
				s.error(w, r, http.StatusUnauthorized, errExpiredRefreshToken)
				return
			}
		}
		u, err := s.Store.User().Find(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

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
		u, err := s.Store.User().Find(req.UserId)

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

		err := s.Store.User().Create(u)
		if err == store.ErrUserAllreadyExist {
			s.error(w, r, http.StatusOK, err)
			return
		}
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitaze()
		token := NewToken()
		token.Auth(u.ID, s.config)
		if err := s.Store.User().UpdateRefreshToken(u.ID, token.RefreshToken, s.config.RefreshTokenExp); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusCreated, token)
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

		u, err := s.Store.User().FindByEmail(req.Email)
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
		if err := s.Store.User().UpdateRefreshToken(u.ID, token.RefreshToken, s.config.RefreshTokenExp); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, token)
	}
}

func (s *APIServer) addCardsPacks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(cxtKeyUser).(*models.User)
		if u.FantacyCoins < 1000 {
			s.respond(w, r, http.StatusNotModified, nil)
			return
		}

		u.FantacyCoins -= 1000
		u.Packs.Common += 1
		if err := s.Store.User().ReplaseUser(u); err != nil {
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

		p, err := s.Store.PlayerCards().OpenCommonPack()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u.Packs.Common -= 1
		packCopy := []models.PlayerCard{}

		if len(u.CardsCollection) == 0 {
			u.CardsCollection = make([][]models.PlayerCard, 5)
			for i, v := range p.Cards {
				packCopy = append(packCopy, v)
				u.CardsCollection[i] = append(u.CardsCollection[i], v)
			}
			s.Store.User().ReplaseUser(u)
			s.respond(w, r, http.StatusOK, packCopy)
			return
		}

		for i := 0; i < len(u.CardsCollection); i++ {
			if len(p.Cards) == 0 {
				break
			}

			for idx, v := range p.Cards {
				if u.CardsCollection[i][0].AccountId == v.AccountId {
					packCopy = append(packCopy, v)
					u.CardsCollection[i] = append(u.CardsCollection[i], v)
					p.Cards = removeCard(p.Cards, idx)
					break
				}
			}
		}
		for i := 0; i < len(u.CardsCollection); i++ {
			if len(p.Cards) == 0 {
				break
			}

			for idx, v := range p.Cards {
				if u.CardsCollection[i][0].AccountId == v.AccountId {
					packCopy = append(packCopy, v)
					u.CardsCollection[i] = append(u.CardsCollection[i], v)
					p.Cards = removeCard(p.Cards, idx)
					break
				}
			}
		}

		for _, v := range p.Cards {
			v.Id = primitive.NewObjectID()
			packCopy = append(packCopy, v)
			n := []models.PlayerCard{v}
			u.CardsCollection = append(u.CardsCollection, n)
		}

		s.Store.User().ReplaseUser(u)

		s.respond(w, r, http.StatusOK, packCopy)
	}
}

func (s *APIServer) collection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(cxtKeyUser).(*models.User)
		s.respond(w, r, http.StatusOK, u.CardsCollection)
	}
}

func (s *APIServer) userData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(cxtKeyUser).(*models.User)

		for i := 0; i < len(u.Teams); i++ {
			if !u.Teams[i].Calculated && !u.Teams[i].Date.UTC().Truncate(24*time.Hour).Equal(time.Now().UTC().Truncate(24*time.Hour)) {
				series, err := s.Store.Series().FindByDate(u.Teams[i].Date)
				if err != nil {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				u.Teams[i].SetPoints(series)

				u.Teams[i].Calculated = true

			}

		}
		if err := s.Store.User().ReplaseUser(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, u)
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
		rar := 1
		for i := 0; i < len(u.CardsCollection); i++ {
			for idx := 0; idx < len(u.CardsCollection[i]); idx++ {
				if u.CardsCollection[i][idx].Id == req.ID {
					rar += u.CardsCollection[i][idx].Rarity
					u.CardsCollection[i] = removeCard(u.CardsCollection[i], idx)
					if len(u.CardsCollection[i]) == 0 {
						u.CardsCollection = removeSlice(u.CardsCollection, i)
					}
					break
				}
			}
		}

		u.FantacyCoins += 100 * rar
		if err := s.Store.User().ReplaseUser(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, u.CardsCollection)
	}
}

func (s *APIServer) setFantacyTeam() http.HandlerFunc {
	type request struct {
		Team []models.PlayerCard `json:"team"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(cxtKeyUser).(*models.User)
		if u.Teams[len(u.Teams)-1].Date.Truncate(24*time.Hour) == time.Now().Truncate(24*time.Hour) {
			s.respond(w, r, http.StatusOK, nil)
		}
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		team := models.NewTeam()
		if len(req.Team) != 5 {
			s.error(w, r, http.StatusBadRequest, fmt.Errorf("len of team err"))
			return
		}

		for i := 0; i < 5; i++ {
			team.Team[i].PlayerCard = req.Team[i]
		}

		u.Teams = append(u.Teams, *team)

		if err := s.Store.User().ReplaseUser(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *APIServer) fantacyTeamsCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(cxtKeyUser).(*models.User)
		for i := 0; i < len(u.Teams); i++ {
			if !u.Teams[i].Calculated {
				series, err := s.Store.Series().FindByDate(u.Teams[i].Date)
				if err != nil {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				u.Teams[i].SetPoints(series)

				if !u.Teams[i].Date.UTC().Truncate(24 * time.Hour).Equal(time.Now().UTC().Truncate(24 * time.Hour)) {
					u.Teams[i].Calculated = true
				}
			}
		}
		if err := s.Store.User().ReplaseUser(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, u.Teams)
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
func removeSlice(s [][]models.PlayerCard, i int) [][]models.PlayerCard {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
