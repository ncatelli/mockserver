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

type StrideHandler struct {
	pass    uint
	stride  uint
	handler Handler
}

// ServeHTTP wraps an inner handler's ServeHTTP but increments the enclosed
// pass by it's stride.
func (sH *StrideHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// incrememt pass by stride
	sH.pass += sH.stride

	sH.handler.ServeHTTP(w, r)
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
	strideHandlers     []StrideHandler
	middlewareHandlers []middleware.Middleware
}

// Init performs any setup and initialization around the route.
func (route *Route) Init() error {

	handlerCount := len(route.Handlers)
	weights := make([]uint, 0, handlerCount)
	for _, h := range route.Handlers {
		weights = append(weights, h.Weight)
	}

	var weightsLCM uint = 0
	if handlerCount == 1 {
		weightsLCM = route.Handlers[0].Weight
	} else if handlerCount > 1 {
		a := weights[0]
		b := weights[1]
		rem := weights[2:]

		weightsLCM = lcm(a, b, rem...)
	}

	for _, h := range route.Handlers {
		stride := weightsLCM / h.Weight
		sH := StrideHandler{
			stride:  stride,
			pass:    0,
			handler: h,
		}

		route.strideHandlers = append(route.strideHandlers, sH)
	}

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
	// set to max of uint so first value is guaranteed to be <= value.
	var lowestPass = ^uint(0)
	lowestPassIdx := 0
	for idx, h := range route.strideHandlers {
		if h.pass < lowestPass {
			lowestPass = h.pass
			lowestPassIdx = idx
		}
	}

	var handler http.Handler = &(route.strideHandlers[lowestPassIdx])

	// Generate handler chain with middlewares
	for i := len(route.middlewareHandlers) - 1; i >= 0; i-- {
		handler = route.middlewareHandlers[i].Middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

func gcd(a, b uint) uint {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b uint, integers ...uint) uint {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
