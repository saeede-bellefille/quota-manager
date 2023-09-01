package main

import (
	"log"

	"github.com/saeede-bellefille/quota-manager/pkg/config"
	"github.com/saeede-bellefille/quota-manager/pkg/server"
)

func main() {
	conf := config.Load()
	server := server.New(conf)
	log.Fatal(server.Start())
}
