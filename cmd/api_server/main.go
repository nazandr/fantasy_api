package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nazandr/fantasy_api/internal/app/api_server"
	"github.com/nazandr/fantasy_api/internal/app/store"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := &api_server.Config{
		IP_addr: ":8080",
		Log_lvl: "debug",
		Store:   &store.Config{},
	}

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	server := api_server.New(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
