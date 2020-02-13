package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/mux"

	"net/http"

	"github.com/ncatelli/mockserver/pkg/config"
	"github.com/ncatelli/mockserver/pkg/router"
	"github.com/ncatelli/mockserver/pkg/router/drivers/simple"
)

func buildRouterFromConfig(c *config.Config) *mux.Router {
	routes, err := simple.LoadFromFile(c.ConfigPath)
	if err != nil {
		panic(err)
	}

	router, err := router.New(routes)
	if err != nil {
		panic(err)
	}

	return router
}

func startHTTPServer(c *config.Config, wg *sync.WaitGroup) *http.Server {
	router := buildRouterFromConfig(c)
	router.HandleFunc(`/healthcheck`, healthHandler).Methods("GET")
	srv := &http.Server{
		Addr:    c.Addr,
		Handler: router,
	}

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// unexpected error. port in use?
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	for {
		c, e := config.New()
		if e != nil {
			log.Fatal("unable to parse config params")
		}

		log.Printf("Starting server on %s\n", c.Addr)

		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := startHTTPServer(&c, httpServerExitDone)

		// blocks for shutdown. If a SIGHUP happens it will gracefully
		// restart the server.
		<-sigs

		log.Println("reloading configuration...")

		if err := srv.Shutdown(context.TODO()); err != nil {
			panic(err) // failure/timeout shutting down the server gracefully
		}

		// wait for goroutine started in startHttpServer() to stop
		httpServerExitDone.Wait()
	}
}
