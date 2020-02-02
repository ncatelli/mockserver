package router

import (
	"log"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	goodConfigPath string = "test_fixtures/good.yaml"
)

var expectedRoutes = []*Route{
	&Route{
		Path:   "/test/weighted/errors",
		Method: "GET",
		Handlers: []Handler{
			Handler{
				Weight: 2,
				ResponseHeaders: map[string]string{
					"content-type": "application/json",
				},
				StaticResponse: "{\"resp\": \"Ok\"}",
				ResponseStatus: 200,
			},
			Handler{
				Weight: 1,
				ResponseHeaders: map[string]string{
					"content-type": "text/plain",
				},
				ResponseStatus: 500,
			},
		},
	},
}

func TestLoadFromFileShould(t *testing.T) {
	t.Run("load a valid configuration", func(t *testing.T) {
		gp, err := filepath.Rel("", goodConfigPath)
		if err != nil {
			log.Fatal(err)
		}

		routes, err := LoadFromFile(gp)
		if err != nil {
			t.Errorf(errFmt, expectedRoutes, err)
		} else if !reflect.DeepEqual(routes, expectedRoutes) {
			t.Errorf(errFmt, expectedRoutes, routes)
		}
	})

	t.Run("return an error a non-yaml parseable file.", func(t *testing.T) {
		bp, err := filepath.Rel("", goodConfigPath)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := LoadFromFile(bp); err != nil {
			t.Errorf(errFmt, "an error", err)
		}
	})
}
