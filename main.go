package main

import (
	"log"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/ncatelli/mockserver/pkg/config"
)

func main() {
	c, e := config.New()
	if e != nil {
		log.Fatal("Unable to parse config file.")
	}

	mux := mux.NewRouter()
	mux.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")

	log.Printf("Starting server on %s\n", c.Addr)
	if err := http.ListenAndServe(c.Addr, mux); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
