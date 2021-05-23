package api_server

import (
	"io"
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
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.loggerConfig(); err != nil {
		return err
	}

	s.routerConfig()
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

func (s *APIServer) routerConfig() {
	s.router.HandleFunc("/players", s.handlePlayerData).Methods("GET")
}

func (s *APIServer) handlePlayerData(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "players[]")
}
