package api_server

import (
	"encoding/json"
	"net/http"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/nazandr/fantasy_api/internal/app/store"
)

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/singup", s.handleUserCreate()).Methods("POST")
}

func (s *APIServer) handleUserCreate() http.HandlerFunc {
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

		u := &models.User{
			Email:    req.Email,
			Password: req.Password,
		}
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

func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
