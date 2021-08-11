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

	t.Run("return handlers in deterministic pattern for unequally-weighted handlers", func(t *testing.T) {
		r := &Route{
			Path:     "/",
			Method:   "GET",
			Handlers: []Handler{successHandler, failureHandler},
		}
		r.Init()

		samples := 10
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

		expectedResponseCodes := []int{200, 500, 200, 500, 200, 200, 500, 200, 200, 500}
		responseCodes := make([]int, 0, len(responses))
		for _, res := range responses {
			responseCodes = append(responseCodes, res.Code)
		}

		if !reflect.DeepEqual(expectedResponseCodes, responseCodes) {
			t.Errorf(errFmt, expectedResponseCodes, responseCodes)
		}
	})

	t.Run("return handlers in deterministic pattern for equally-weighted handlers", func(t *testing.T) {
		equalSuccessHandler := successHandler
		equalSuccessHandler.Weight = 1
		r := &Route{
			Path:     "/",
			Method:   "GET",
			Handlers: []Handler{equalSuccessHandler, failureHandler},
		}
		r.Init()

		samples := 10
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

		expectedResponseCodes := []int{200, 500, 200, 500, 200, 500, 200, 500, 200, 500}
		responseCodes := make([]int, 0, len(responses))
		for _, res := range responses {
			responseCodes = append(responseCodes, res.Code)
		}

		if !reflect.DeepEqual(expectedResponseCodes, responseCodes) {
			t.Errorf(errFmt, expectedResponseCodes, responseCodes)
		}
	})
}

func BenchmarkRouterHandlerSelectionWith(b *testing.B) {
	successHandler := Handler{
		Weight:         1,
		ResponseStatus: 200,
		StaticResponse: "Ok",
	}
	failureHandler := Handler{
		Weight:         1,
		ResponseStatus: 500,
		StaticResponse: "",
	}
	b.Run("equally-weighted handlers", func(b *testing.B) {
		r := &Route{
			Path:     "/",
			Method:   "GET",
			Handlers: []Handler{successHandler, failureHandler},
		}
		r.Init()

		for i := 0; i < b.N; i++ {
			<-r.handlerChan
		}

	})

	b.Run("unequally-weighted handlers", func(b *testing.B) {
		unequalSuccessHandler := successHandler
		unequalSuccessHandler.Weight = 2
		r := &Route{
			Path:     "/",
			Method:   "GET",
			Handlers: []Handler{unequalSuccessHandler, failureHandler},
		}
		r.Init()

		for i := 0; i < b.N; i++ {
			<-r.handlerChan
		}

	})

}
