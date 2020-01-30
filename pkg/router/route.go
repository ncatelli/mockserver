package router

import "net/http"

// Route includes all routing data to build a route and forward to an
// appropriate router. This is handed off to the router for the live routing.
type Route struct {
	Path           string            `yaml:"path"`
	Method         string            `yaml:"method"`
	QueryParams    map[string]string `yaml:"query_params"`
	RequestHeaders map[string]string `yaml:"request_headers"`
	Handlers       []Handler         `yaml:"handlers"`
	totalWeight    uint
}

// ServeHTTP implements the http.Handler interface for pipelining a request
// further into a handler.
func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
