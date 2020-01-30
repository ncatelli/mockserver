package router

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestHandlerUnmarshaling(t *testing.T) {
	t.Run("should unmarshal to the correct keys", func(t *testing.T) {
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
