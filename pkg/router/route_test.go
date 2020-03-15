package router

import (
	"fmt"
	"math"
	"math/rand"
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
		got, _ := calculateTotalWeightofHandlers(tHandlers)

		if expected != got {
			t.Errorf(errFmt, expected, got)
		}
	})

	t.Run("throw an error if the max handler weight would overflow an int64", func(t *testing.T) {
		tHandlers := []Handler{
			Handler{Weight: maxInt64},
			Handler{Weight: 2},
		}

		got, err := calculateTotalWeightofHandlers(tHandlers)

		if got != -1 || err == nil {
			t.Errorf(errFmt, ErrInvalidWeight{handler: &tHandlers[1]}, nil)
		}
	})

	t.Run("throw an error if a negative weight is defined.", func(t *testing.T) {
		tHandlers := []Handler{
			Handler{Weight: -1},
		}

		got, err := calculateTotalWeightofHandlers(tHandlers)

		if got != -1 || err == nil {
			t.Errorf(errFmt, ErrInvalidWeight{handler: &tHandlers[0]}, nil)
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

func BenchmarkRouterHandlerSelectionWith(b *testing.B) {
	for x := 0.0; x <= 6; x++ {
		pow := math.Pow(2, x)
		h := generateTestHandlerSlice(1, int(pow))
		b.Run(fmt.Sprintf("%.0f equally-weighted handlers", pow), func(b *testing.B) {
			r := Route{
				Handlers: h,
			}

			for i := 0; i < b.N; i++ {
				r.selectHandler(rand.Intn(r.totalWeight + 1))
			}
		})
	}
}

func generateTestHandlerSlice(weight int, count int) []Handler {
	h := make([]Handler, count)

	for i := 0; i < count; i++ {
		h = append(h, Handler{Weight: weight})
	}

	return h
}
