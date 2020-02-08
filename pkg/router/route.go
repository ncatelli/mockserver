package router

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Route includes all routing data to build a route and forward to an
// appropriate router. This is handed off to the router for the live routing.
type Route struct {
	Path           string            `yaml:"path"`
	Method         string            `yaml:"method"`
	QueryParams    map[string]string `yaml:"query_params"`
	RequestHeaders map[string]string `yaml:"request_headers"`
	Handlers       []Handler         `yaml:"handlers"`
	totalWeight    int
}

// ServeHTTP implements the http.Handler interface for pipelining a request
// further into a handler.
func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if route.totalWeight == 0 {
		route.totalWeight = calculateTotalWeightofHandlers(route.Handlers)
	}

	hw := rand.Intn(route.totalWeight + 1)

	for _, h := range route.Handlers {
		hw -= h.Weight
		if hw <= 0 {
			h.ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "unable to determine routable handler")
}

// calculateTotalWeightofHandlers iterates over all handlers assigned to the
// route and sums their total weight.
func calculateTotalWeightofHandlers(handlers []Handler) int {
	var weight int

	for _, h := range handlers {
		weight += h.Weight
	}

	return weight
}
