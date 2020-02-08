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

		expected := 3
		got := tRoute.totalWeight

		if expected != got {
			t.Errorf(errFmt, expected, got)
		}
	})
}
func TestRouteWeightCalculationShould(t *testing.T) {
	t.Run("return the sum of all handler weights", func(t *testing.T) {
		tHandlers := []Handler{
			Handler{Weight: 1},
			Handler{Weight: 10},
			Handler{Weight: 3},
		}

		expected := 14
		got := calculateTotalWeightofHandlers(tHandlers)

		if expected != got {
			t.Errorf(errFmt, expected, got)
		}
	})
}
