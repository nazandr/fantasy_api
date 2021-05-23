package main

import (
	"log"

	"github.com/nazandr/fantasy_api/internal/app/api_server"
)

func main() {
	conf := api_server.Config{
		IP_addr: ":8080",
		Log_lvl: "debug",
	}
	server := api_server.New(&conf)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
