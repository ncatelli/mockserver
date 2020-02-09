package logging

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

const (
	errFmt string = "want %v, got %v"
)

func TestLatencyMiddlewareShould(t *testing.T) {
	t.Run("stay within a range of responses.", func(t *testing.T) {
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

				if duration.Milliseconds() < 1000 || duration.Milliseconds() > 2000 {
					t.Errorf(errFmt, "between 1000 - 2000", duration.Milliseconds())
				}
			}()
		}

		wg.Wait()
	})
}
