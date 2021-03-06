package router

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/ncatelli/mockserver/pkg/router/middleware"
)

const (
	maxInt64 = 1<<63 - 1
)

// ErrInvalidWeight is thrown when a handler has a weight outside the
// acceptable bounds.
type ErrInvalidWeight struct {
	handler *Handler
}

func (e ErrInvalidWeight) Error() string {
	return fmt.Sprintf("handler %v exceeds maximum total weight of %d", *e.handler, maxInt64)
}

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
	tw, err := calculateTotalWeightofHandlers(route.Handlers)
	if err != nil {
		return err
	}

	route.totalWeight = tw

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
	handler := route.selectHandler(rand.Intn(route.totalWeight + 1))

	// Generate handler chain with middlewares
	for i := len(route.middlewareHandlers) - 1; i >= 0; i-- {
		handler = route.middlewareHandlers[i].Middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// selectHandler randomly selects and returns a handler from the route
// handlers pool.
func (route *Route) selectHandler(hw int) http.Handler {
	var handler http.Handler

	// make handler selection based on weight.
	for _, h := range route.Handlers {
		hw -= h.Weight
		if hw <= 0 {
			handler = &h
			break
		}
	}

	return handler
}

// calculateTotalWeightofHandlers iterates over all handlers assigned to the
// route and sums their total weight.
func calculateTotalWeightofHandlers(handlers []Handler) (int, error) {
	var totalWeight int

	for _, h := range handlers {
		maxhandlerWeight := maxInt64 - totalWeight

		if h.Weight > maxInt64 || h.Weight > maxhandlerWeight || h.Weight < 0 {
			return -1, ErrInvalidWeight{
				handler: &h,
			}
		}

		totalWeight += h.Weight
	}

	return totalWeight, nil
}
