package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router provides a wrapper for gorilla mux to fascilitate registering, and
// routing to mock methods.
type Router struct {
	Mux *mux.Router
}

// New initializes and returns an instance of Router.
func New(routes []Route) *Router {
	return &Router{
		Mux: mux.NewRouter(),
	}
}

// ServeHTTP handles wrapping the mux ServeHTTP method.
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.Mux.ServeHTTP(w, r)
}
