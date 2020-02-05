package logging

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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
		logMiddleWare := &Logger{}
		logMiddleWare.Init(map[string]interface{}{
			"target": logBuffer,
		})

		router.Use(logMiddleWare.Middleware)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if !bytes.Contains(logBuffer.Bytes(), []byte("GET /")) {
			t.Errorf(errFmt, "GET /", logBuffer.Bytes())
		}
	})
}
