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

type Config struct {
	IP_addr string `toml:"ip_addr"`
	Log_lvl string `toml:"log_lvl"`
	Store   *store.Config
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
		store:  config.Store.New(),
	}
}

func (s *APIServer) Start() error {
	if err := s.loggerConfig(); err != nil {
		return err
	}

	s.routerConfig()
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

func (s *APIServer) routerConfig() {
	s.router.HandleFunc("/players", s.handlePlayerData).Methods("GET")
}

func (s *APIServer) handlePlayerData(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "players[]")
}
