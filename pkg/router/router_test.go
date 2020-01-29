package router

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handlerTest(method, path string, reqBody io.Reader, respCode int, respBody string, h http.HandlerFunc) error {
	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		return err
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := New()
	router.Mux.HandleFunc(path, h).Methods(method)

	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != respCode {
		return fmt.Errorf("handler returned wrong status code: got %v want %v",
			status, respCode)
	}

	if respBody != "" {
		// Check the response body is what we expect.
		if rr.Body.String() != respBody {
			return fmt.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), respBody)
		}
	}

	return nil
}

func TestRouterShould(t *testing.T) {
	t.Run("correctly hand off routing to mux", func(t *testing.T) {
		err := handlerTest("GET", "/test", nil, http.StatusOK, `Ok`, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `Ok`)
		})
		if err != nil {
			t.Error(err)
		}
	})
}
