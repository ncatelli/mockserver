package router

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/ncatelli/mockserver/pkg/router/middleware"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Route includes all routing data to build a route and forward to an
// appropriate router. This is handed off to the router for the live routing.
type Route struct {
	Path               string                       `yaml:"path"`
	Method             string                       `yaml:"method"`
	QueryParams        map[string]string            `yaml:"query_params"`
	RequestHeaders     map[string]string            `yaml:"request_headers"`
	Middleware         map[string]map[string]string `yaml:"middleware"`
	Handlers           []Handler                    `yaml:"handlers"`
	middlewareHandlers []middleware.Middleware
	totalWeight        int
}

// Init performs any setup and initialization around the route.
func (route *Route) Init() error {
	route.totalWeight = calculateTotalWeightofHandlers(route.Handlers)

	for k, v := range route.Middleware {
		m := middleware.Lookup(k)
		if m == nil {
			return middleware.ErrUndefinedMiddleware{ID: k}
		}

		if err := m.Init(v); err != nil {
			return err
		}

		route.middlewareHandlers = append(route.middlewareHandlers, m)
	}

	return nil
}

// ServeHTTP implements the http.Handler interface for pipelining a request
// further into a handler.
func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler
	hw := rand.Intn(route.totalWeight + 1)

	// make handler selection based on weight.
	for _, h := range route.Handlers {
		hw -= h.Weight
		if hw <= 0 {
			handler = &h
		}
	}

	// Generate handler chain with middlewares
	for i := len(route.middlewareHandlers) - 1; i >= 0; i-- {
		handler = route.middlewareHandlers[i].Middleware(handler)
	}

	handler.ServeHTTP(w, r)
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
