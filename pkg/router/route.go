package router

import (
	"fmt"
	"math"
	"net/http"

	"github.com/ncatelli/mockserver/pkg/router/middleware"
)

// ErrInvalidWeight is thrown when a handler has a weight outside the
// acceptable bounds.
type ErrInvalidWeight struct {
	handler *Handler
}

func (e ErrInvalidWeight) Error() string {
	return fmt.Sprintf("handler %v exceeds maximum total weight of %v", *e.handler, math.MaxInt64)
}

// StrideHandlers wraps the Handler type with a precomputed stride and pass context.
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
	middlewareHandlers []middleware.Middleware
	handlerChan        chan http.Handler
}

// Init performs any setup and initialization around the route.
func (route *Route) Init() error {
	route.handlerChan = make(chan http.Handler, 1024)

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

	go func(handler []Handler, handlerQueue chan http.Handler) {
		var weightsLCM uint = 0
		handlerCount := len(handler)
		weights := make([]uint, 0, handlerCount)
		strideHandlers := make([]StrideHandler, 0, handlerCount)

		for _, v := range handler {
			weights = append(weights, v.Weight)
		}

		if handlerCount == 1 {
			weightsLCM = weights[0]
		} else if handlerCount > 1 {
			a := weights[0]
			b := weights[1]
			rem := weights[2:]

			weightsLCM = lcm(a, b, rem...)
		}

		for _, h := range handler {
			stride := weightsLCM / h.Weight
			sH := StrideHandler{
				stride:  stride,
				pass:    0,
				handler: h,
			}

			strideHandlers = append(strideHandlers, sH)
		}

		for {
			// set to max of uint so first value is guaranteed to be <= value.
			var lowestPass uint = math.MaxUint32
			lowestPassIdx := 0
			for idx, h := range strideHandlers {
				if h.pass < lowestPass {
					lowestPass = h.pass
					lowestPassIdx = idx
				}
			}

			sH := strideHandlers[lowestPassIdx]

			// append strideHandler to the end of the queue
			strideHandlers = append(
				append(strideHandlers[:lowestPassIdx], strideHandlers[lowestPassIdx+1:]...),
				sH)

			// Generate handler chain with middlewares
			handlerQueue <- &sH
		}
	}(route.Handlers, route.handlerChan)

	return nil
}

// ServeHTTP implements the http.Handler interface for pipelining a request
// further into a handler.
func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := <-route.handlerChan

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
