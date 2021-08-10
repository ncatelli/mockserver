package router

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func TestRouteUnmarshalingShould(t *testing.T) {
	t.Run("unmarshal to the correct keys", func(t *testing.T) {
		expectedRoute := Route{
			Path:   "/",
			Method: "GET",
		}

		route := Route{}
		rawRoute := []byte(`{"path": "/", "method": "GET"}`)
		err := yaml.Unmarshal(rawRoute, &route)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(expectedRoute, route) {
			t.Errorf(errFmt, expectedRoute, route)
		}
	})
}

func TestHandlerSelectionShould(t *testing.T) {
	successHandler := Handler{
		Weight:         2,
		ResponseStatus: 200,
		StaticResponse: "Ok",
	}
	failureHandler := Handler{
		Weight:         1,
		ResponseStatus: 500,
		StaticResponse: "",
	}

	t.Run("return handlers in deterministic pattern for weight", func(t *testing.T) {
		r := &Route{
			Path:     "/",
			Method:   "GET",
			Handlers: []Handler{successHandler, failureHandler},
		}
		r.Init()

		samples := 9
		responses := make([]*httptest.ResponseRecorder, 0, samples)

		for i := 0; i < samples; i++ {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Handle("/", r).Methods("GET")

			router.ServeHTTP(rr, req)
			responses = append(responses, rr)
		}

		expectedResponseCodes := []int{200, 500, 200, 500, 200, 200, 500, 200, 200}
		responseCodes := make([]int, 0, len(responses))
		for _, res := range responses {
			responseCodes = append(responseCodes, res.Code)
		}

		if !reflect.DeepEqual(expectedResponseCodes, responseCodes) {
			t.Errorf(errFmt, expectedResponseCodes, responseCodes)
		}
	})
}
