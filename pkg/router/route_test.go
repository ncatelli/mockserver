package router

import (
	"reflect"
	"testing"

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

func TestRouteInitShould(t *testing.T) {
	t.Run("assign the total weight of all handlers to the route", func(t *testing.T) {
		tRoute := Route{
			Handlers: []Handler{
				Handler{Weight: 1},
				Handler{Weight: 2},
			},
		}

		tRoute.Init()

		expected := []int{2, 1}
		got := make([]int, 0, len(tRoute.strideHandlers))
		for _, sH := range tRoute.strideHandlers {
			got = append(got, sH.stride)
		}

		for i, expectedStride := range expected {
			if expectedStride != got[i] {
				t.Errorf(errFmt, expectedStride, got[i])
			}
		}

	})
}

func TestHandlerSelectionShould(t *testing.T) {
	t.Run("return the first handler that brings the handler weight to 0", func(t *testing.T) {
		r := Route{
			Handlers: []Handler{
				Handler{Weight: 2},
				Handler{Weight: 1},
			},
		}

		expected := &r.Handlers[0]
		got := r.selectHandler(2)

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(errFmt, expected, got)
		}
	})
}
