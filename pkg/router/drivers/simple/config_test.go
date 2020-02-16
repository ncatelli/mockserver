package simple

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ncatelli/mockserver/pkg/router"
)

const (
	goodConfigPath string = "test_fixtures/good.yaml"
	badConfigPath  string = "test_fixtures/bad.yaml"
	errFmt         string = "want %v, got %v"
)

var expectedRoutes = []*router.Route{
	&router.Route{
		Path:   "/test/weighted/errors",
		Method: "GET",
		Handlers: []router.Handler{
			router.Handler{
				Weight: 2,
				ResponseHeaders: map[string]string{
					"content-type": "application/json",
				},
				StaticResponse: "{\"resp\": \"Ok\"}",
				ResponseStatus: 200,
			},
			router.Handler{
				Weight: 1,
				ResponseHeaders: map[string]string{
					"content-type": "text/plain",
				},
				ResponseStatus: 500,
			},
		},
	},
}

func TestLoadShould(t *testing.T) {
	t.Run("load a valid configuration", func(t *testing.T) {
		gp, err := filepath.Rel("", goodConfigPath)
		file, err := os.Open(gp)
		defer file.Close()

		routes, err := Load(file)
		if err != nil {
			t.Errorf(errFmt, expectedRoutes, err)
		} else if !reflect.DeepEqual(routes, expectedRoutes) {
			t.Errorf(errFmt, expectedRoutes, routes)
		}
	})

	t.Run("return an error a non-yaml parseable file.", func(t *testing.T) {
		bp, err := filepath.Rel("", badConfigPath)
		file, err := os.Open(bp)
		defer file.Close()

		if err != nil {
			log.Fatal(err)
		}

		if _, err := LoadFromFile(bp); err == nil {
			t.Errorf(errFmt, "an error", err)
		}
	})
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
		bp, err := filepath.Rel("", badConfigPath)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := LoadFromFile(bp); err == nil {
			t.Errorf(errFmt, "an error", err)
		}
	})
}
