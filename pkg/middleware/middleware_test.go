package middleware

import (
	"net/http"
	"testing"
)

const (
	errFmt string = "want %v, got %v"
)

type testMiddleware struct{}

func (tm *testMiddleware) Init(conf map[string]string) error {
	return nil
}

func (tm *testMiddleware) Middleware(next http.Handler) http.Handler {
	return next
}

func TestMiddlewareLookupShould(t *testing.T) {
	t.Run("return a middleware if it exists in the map", func(t *testing.T) {
		em := &testMiddleware{}
		middlewares["test"] = em
		defer delete(middlewares, "test")

		if m := Lookup("test"); m == nil {
			t.Errorf(errFmt, em, m)
		}
	})

	t.Run("return nil if the middleware isn't registered in the map", func(t *testing.T) {
		if m := Lookup("test_middleware_shouldn't_exist"); m != nil {
			t.Errorf(errFmt, nil, m)
		}
	})
}
