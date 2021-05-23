package main

import (
	"flag"
	"fmt"
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

	config := new(api_server.Config)
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)
	// conf := api_server.Config{
	// 	IP_addr: ":8080",
	// 	Log_lvl: "debug",
	// 	Store:   store.NewConfig(),
	// }
	server := api_server.New(config)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
