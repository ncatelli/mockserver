package router

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

type ErrRouteMatchFailure struct{}

func generateTestHandler() []Handler {
	return []Handler{Handler{
		Weight:         1,
		StaticResponse: "Ok",
		ResponseStatus: 200,
	}}
}

func routerHelper(req *http.Request, route *Route) bool {
	rm := &mux.RouteMatch{}
	router, _ := New([]*Route{route})

	return router.Match(req, rm)
}

func TestRouterShouldMatch(t *testing.T) {
	t.Run("when a route with a valid path and method", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Error(err)
		}

		route := &Route{
			Path:   "/test",
			Method: "GET",
		}

		if !routerHelper(req, route) {
			t.Errorf(errFmt, true, false)
		}
	})

	t.Run("when a route with a path var is present", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test/1", nil)
		if err != nil {
			t.Error(err)
		}

		route := &Route{
			Path:   "/test/{key}",
			Method: "GET",
		}

		if !routerHelper(req, route) {
			t.Errorf(errFmt, true, false)
		}
	})

	t.Run("when a route with a header specification exists", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Error(err)
		}

		req.Header.Add("TestHeader", "present")

		route := &Route{
			Path:   "/test",
			Method: "GET",
			RequestHeaders: map[string]string{
				"TestHeader": "present",
			},
		}

		if !routerHelper(req, route) {
			t.Errorf(errFmt, true, false)
		}
	})

	t.Run("when a route with a query param specification exists", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test?testparam=present", nil)
		if err != nil {
			t.Error(err)
		}

		route := &Route{
			Path:   "/test",
			Method: "GET",
			QueryParams: map[string]string{
				"testparam": "present",
			},
		}

		if !routerHelper(req, route) {
			t.Errorf(errFmt, true, false)
		}
	})
}

func TestRouterShouldNotMatch(t *testing.T) {
	t.Run("not match when a route with a valid path but invalid method are specified", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Error(err)
		}

		route := &Route{
			Path:   "/test",
			Method: "POST",
		}

		if routerHelper(req, route) {
			t.Errorf(errFmt, true, false)
		}
	})

	t.Run("not match when a route has a header requirement and the header isn't set", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Error(err)
		}

		route := &Route{
			Path:   "/test",
			Method: "GET",
			RequestHeaders: map[string]string{
				"TestHeader": "present",
			},
		}

		if routerHelper(req, route) {
			t.Errorf(errFmt, true, false)
		}
	})

	t.Run("not match when a route has a query parameter requirement and the header isn't set", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Error(err)
		}

		route := &Route{
			Path:   "/test",
			Method: "GET",
			QueryParams: map[string]string{
				"testparam": "present",
			},
		}

		if routerHelper(req, route) {
			t.Errorf(errFmt, true, false)
		}
	})
}
