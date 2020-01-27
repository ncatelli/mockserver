package main

import (
	"log"

	"net/http"

	"github.com/gorilla/mux"
)

const (
	addr string = "0.0.0.0:8080"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")

	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe("addr", mux); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
