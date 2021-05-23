package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/nazandr/fantasy_api/internal/app/api_server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := api_server.NewConfig()

	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config.Store)
	server := api_server.New(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
