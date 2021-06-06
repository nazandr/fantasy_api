package server

import (
	"github.com/nazandr/fantasy_api/internal/app/store"
)

type Config struct {
	IP_addr         string `toml:"ip_addr"`
	Log_lvl         string `toml:"log_lvl"`
	AcssesTokenExp  int    `toml:"acsses_token_exp"`
	RefreshTokenExp int    `toml:"refresh_token_exp"`
	SignatureKey    string
	Store           *store.Config
}

func NewConfig() *Config {
	return &Config{
		IP_addr:         ":8080",
		Log_lvl:         "debug",
		AcssesTokenExp:  15,
		RefreshTokenExp: 30,
		SignatureKey:    "secret",
		Store:           store.NewConfig(),
	}
}
