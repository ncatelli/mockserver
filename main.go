package main

import (
	"log"

	"net/http"

	"github.com/ncatelli/mockserver/pkg/config"
	"github.com/ncatelli/mockserver/pkg/router"
)

func main() {
	c, e := config.New()
	if e != nil {
		log.Fatal("Unable to parse config file.")
	}

	router := router.New()
	router.Mux.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")

	log.Printf("Starting server on %s\n", c.Addr)
	if err := http.ListenAndServe(c.Addr, router); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
