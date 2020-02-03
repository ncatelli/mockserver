package main

import (
	"log"

	"net/http"

	"github.com/ncatelli/mockserver/pkg/config"
	"github.com/ncatelli/mockserver/pkg/router"
	"github.com/ncatelli/mockserver/pkg/router/drivers/simple"
)

func main() {
	c, e := config.New()
	if e != nil {
		log.Fatal("unable to parse config params")
	}

	routes, err := simple.LoadFromFile(c.ConfigPath)
	if err != nil {
		panic(err)
	}

	router := router.New(routes)
	router.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")

	log.Printf("Starting server on %s\n", c.Addr)
	if err := http.ListenAndServe(c.Addr, router); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
