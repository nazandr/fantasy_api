package api_server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nazandr/fantasy_api/internal/app/store"
	"github.com/sirupsen/logrus"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store  *store.Store
}

func New(config *Config) *APIServer {
	s := &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *APIServer) Start() error {
	if err := s.loggerConfig(); err != nil {
		return err
	}

	if err := s.storeConfig(); err != nil {
		return err
	}
	s.logger.Info("api server started")

	return http.ListenAndServe(s.config.IP_addr, s.router)
}

func (s *APIServer) loggerConfig() error {
	lvl, err := logrus.ParseLevel(s.config.Log_lvl)
	if err != nil {
		return err
	}

	s.logger.SetLevel(lvl)

	return nil
}

func (s *APIServer) storeConfig() error {
	st := store.New(s.config.Store)

	if err := st.Connect(); err != nil {
		return err
	}
	s.store = st

	s.logger.Info("Connected to DB")
	return nil
}
