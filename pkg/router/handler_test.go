package router

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func TestHandlerUnmarshalingShould(t *testing.T) {
	t.Run("unmarshal to the correct keys", func(t *testing.T) {
		expectedHandler := Handler{
			Weight: 1,
		}

		handler := Handler{}
		rawHandler := []byte(`{"weight": 1}`)
		err := yaml.Unmarshal(rawHandler, &handler)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expectedHandler, handler) {
			t.Errorf(errFmt, expectedHandler, handler)
		}
	})
}

func TestHandlerServeHTTPMethodShould(t *testing.T) {
	t.Run("serve a static response when set", func(t *testing.T) {
		h := &Handler{
			StaticResponse: "Ok",
			ResponseStatus: 200,
		}

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.Handle("/", h).Methods("GET")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the response body is what we expect.
		expected := `Ok`
		if rr.Body.String() != expected {
			t.Errorf(errFmt,
				expected, rr.Body.String())
		}
	})

	t.Run("set response headers when specified on the handler", func(t *testing.T) {
		er := `"Ok"`
		h := &Handler{
			StaticResponse: er,
			ResponseHeaders: map[string]string{
				"content-type": "application/json",
			},
			ResponseStatus: 200,
		}

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.Handle("/", h).Methods("GET")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf(errFmt, http.StatusOK, status)
		}

		// Check the response body is what we expect.
		if rr.Body.String() != er {
			t.Errorf(errFmt, er, rr.Body.String())
		}

		ct := rr.Header().Get("content-type")
		if ct != "application/json" {
			t.Errorf(errFmt, "application/json", ct)
		}
	})

	t.Run("serve a template from a file when specified", func(t *testing.T) {
		h := &Handler{
			ResponseStatus: 200,
			ResponsePath:   "test_fixtures/good_simple_string.txt",
		}

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := mux.NewRouter()
		router.Handle("/", h).Methods("GET")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf(errFmt, http.StatusOK, status)
		}

		// Check the response body is what we expect.
		expected := `Ok`
		if rr.Body.String() != expected {
			t.Errorf(errFmt,
				expected, rr.Body.String())
		}
	})
}
