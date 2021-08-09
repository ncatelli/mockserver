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
	t.Run("generate the expected stride value for each route.", func(t *testing.T) {
		tRoute := Route{
			Handlers: []Handler{
				{Weight: 1},
				{Weight: 2},
			},
		}

		tRoute.Init()

		expected := []uint{2, 1}
		got := make([]uint, 0, len(tRoute.strideHandlers))
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

}
