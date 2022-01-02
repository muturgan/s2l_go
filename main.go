package main

import (
	"log"

	"github.com/muturgan/s2l_go/src/config"
	"github.com/muturgan/s2l_go/src/dal"
	"github.com/muturgan/s2l_go/src/server"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	dal, err := dal.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer dal.Stop()

	server.Serve(config, dal)
}
