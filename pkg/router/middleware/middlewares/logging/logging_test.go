package logging

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
)

const (
	errFmt string = "want %v, got %v"
)

func TestLoggingMiddlewareShould(t *testing.T) {
	t.Run("log output with the expected starting prefix", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")

		logBuffer := new(bytes.Buffer)
		logMiddleware := &Middleware{}
		logMiddleware.target = logBuffer

		router.Use(logMiddleware.Middleware)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if !bytes.Contains(logBuffer.Bytes(), []byte("GET /")) {
			t.Errorf(errFmt, "GET /", logBuffer.Bytes())
		}
	})

	t.Run("logger target should default to os.Stdout if no target is passed to Init", func(t *testing.T) {
		logMiddleware := &Middleware{}
		logMiddleware.Init(map[string]string{})

		if logMiddleware.target != os.Stdout {
			t.Errorf(errFmt, os.Stdout, logMiddleware.target)
		}
	})

	t.Run("logger target should be set to os.Stdout if stdout is passed in Init", func(t *testing.T) {
		logMiddleware := &Middleware{}
		logMiddleware.Init(map[string]string{
			"target": "stdout",
		})

		if logMiddleware.target != os.Stdout {
			t.Errorf(errFmt, os.Stdout, logMiddleware.target)
		}
	})

	t.Run("throw an error if an unknown target is specified", func(t *testing.T) {
		logMiddleware := &Middleware{}
		err := logMiddleware.Init(map[string]string{
			"target": "this target doesn't exist",
		})

		if err == nil {
			t.Errorf(errFmt, "error", nil)
		}
	})
}
