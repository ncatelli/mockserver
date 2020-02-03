package router

import (
	"github.com/gorilla/mux"
)

// New takes a list of routes and attempts to return a router with all of these
// routes registered to it.
func New(routes []*Route) *mux.Router {
	m := mux.NewRouter()

	for _, r := range routes {
		route := m.Handle(r.Path, r).Methods(r.Method)

		for k, v := range r.RequestHeaders {
			route.Headers(k, v)
		}

		for k, v := range r.QueryParams {
			route.Queries(k, v)
		}
	}

	return m
}
