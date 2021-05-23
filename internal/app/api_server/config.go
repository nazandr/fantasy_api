package api_server

import "github.com/nazandr/fantasy_api/internal/app/store"

type Config struct {
	IP_addr string `toml:"ip_addr"`
	Log_lvl string `toml:"log_lvl"`
	Store   *store.Config
}

func NewConfig() *Config {
	return &Config{
		IP_addr: ":8080",
		Log_lvl: "debug",
		Store:   store.NewConfig(),
	}
}
