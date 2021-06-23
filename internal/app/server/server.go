package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nazandr/fantasy_api/internal/app/store"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	Logger *logrus.Logger
	router *mux.Router
	Store  *store.Store
}

func New(config *Config) *APIServer {
	s := &APIServer{
		config: config,
		Logger: logrus.New(),
		router: mux.NewRouter(),
	}
	s.configureRouter()
	if err := s.storeConfig(); err != nil {
		log.Fatal(err)
	}
	return s
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *APIServer) Start() error {
	if err := s.loggerConfig(); err != nil {
		return err
	}

	s.Logger.Info("api server started")

	return http.ListenAndServe(s.config.IP_addr, s.router)
}

func (s *APIServer) loggerConfig() error {
	lvl, err := logrus.ParseLevel(s.config.Log_lvl)
	if err != nil {
		return err
	}

	s.Logger.SetLevel(lvl)

	return nil
}

func (s *APIServer) storeConfig() error {
	st := store.New(s.config.Store)

	if err := st.Connect(); err != nil {
		return err
	}
	s.Store = st

	s.Logger.Info("Connected to DB")
	return nil
}
