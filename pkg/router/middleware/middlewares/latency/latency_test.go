package latency

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

const (
	errFmt string = "want %v, got %v"
)

func TestLatencyMiddlewareInitShould(t *testing.T) {
	t.Run("parse a valid latency parameter", func(t *testing.T) {
		unparsedLatency := "50"
		expectedLatency := 50

		m := &Middleware{}
		conf := map[string]string{"latency": unparsedLatency}

		err := m.Init(conf)
		if err != nil {
			t.Errorf(errFmt, nil, err)
		}

		if m.Latency != expectedLatency {
			t.Errorf(errFmt, expectedLatency, m.Latency)
		}
	})

	t.Run("throw an err when the latency parameter is invalid", func(t *testing.T) {
		unparsedInvalidLatency := "invalidParam"

		m := &Middleware{}
		conf := map[string]string{"latency": unparsedInvalidLatency}

		err := m.Init(conf)
		if err == nil {
			t.Errorf(errFmt, "error", nil)
		}
	})

	t.Run("parse valid min/max parameters", func(t *testing.T) {
		unparsedMin := "50"
		expectedMin := 50
		unparsedMax := "100"
		expectedMax := 100

		expectedMiddleware := &Middleware{
			Min: expectedMin,
			Max: expectedMax,
		}

		m := &Middleware{}
		conf := map[string]string{
			"min": unparsedMin,
			"max": unparsedMax,
		}

		err := m.Init(conf)
		if err != nil {
			t.Errorf(errFmt, nil, err)
		}

		if !reflect.DeepEqual(expectedMiddleware, m) {
			t.Errorf(errFmt, expectedMiddleware, m)
		}
	})

	t.Run("throw an err when the min/max parameters is invalid", func(t *testing.T) {
		unparsedInvalidParam := "invalidParam"

		m := &Middleware{}
		conf := map[string]string{
			"min": unparsedInvalidParam,
			"max": unparsedInvalidParam,
		}

		err := m.Init(conf)
		if err == nil {
			t.Errorf(errFmt, "error", nil)
		}
	})
}

func TestLatencyMiddlewareShould(t *testing.T) {
	t.Run("stay within a range of responses of latency setting", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 1000; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}

				router := mux.NewRouter()
				router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

				latencyMiddleware := &Middleware{Latency: 1000}

				router.Use(latencyMiddleware.Middleware)

				rr := httptest.NewRecorder()

				// start timer
				start := time.Now()
				router.ServeHTTP(rr, req)
				duration := time.Since(start)

				if duration.Milliseconds() < 1000 || duration.Milliseconds() > 1500 {
					t.Errorf(errFmt, "between 1000 - 1500", duration.Milliseconds())
				}
			}()
		}

		wg.Wait()
	})

	t.Run("prefer explicit latency setting over min/max", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		latencyMiddleware := &Middleware{
			Latency: 1000,
			Min:     2000,
			Max:     3000,
		}

		router.Use(latencyMiddleware.Middleware)

		rr := httptest.NewRecorder()

		// start timer
		start := time.Now()
		router.ServeHTTP(rr, req)
		duration := time.Since(start)

		if duration.Milliseconds() < 1000 || duration.Milliseconds() > 1500 {
			t.Errorf(errFmt, "between 1000 - 1500", duration.Milliseconds())
		}
	})

	t.Run("stay within a range of responses within min/max", func(t *testing.T) {
		var expectedMin = 1000
		var expectedMax = 2000

		var wg sync.WaitGroup

		for i := 0; i < 1000; i++ {
			wg.Add(1)

			go func() {
				defer wg.Done()
				req, err := http.NewRequest("GET", "/", nil)
				if err != nil {
					t.Fatal(err)
				}

				router := mux.NewRouter()
				router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

				latencyMiddleware := &Middleware{
					Min: expectedMin,
					Max: expectedMax,
				}

				router.Use(latencyMiddleware.Middleware)

				rr := httptest.NewRecorder()

				// start timer
				start := time.Now()
				router.ServeHTTP(rr, req)
				duration := time.Since(start)

				if duration.Milliseconds() < int64(expectedMin) || duration.Milliseconds() > int64(expectedMax+100) {
					t.Errorf(errFmt, fmt.Sprintf("between %d - %d", expectedMin, expectedMax), duration.Milliseconds())
				}
			}()
		}

		wg.Wait()
	})
}
