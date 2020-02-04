package middleware

import "net/http"

var (
	middlewares = make(map[string]Middleware)
)

// Middleware defines the necessary functions to configure and implement a
// middleware for use on a route.
type Middleware interface {
	Init(map[string]interface{}) error
	Middleware(http.Handler) http.Handler
}

// Lookup takes an id and attempts to return the corresponding middleware if
// the middleware is undefined nil is returned.
func Lookup(id string) Middleware {
	if m, prs := middlewares["id"]; prs == true {
		return m
	}

	return nil
}
