package openapi

import (
	"bytes"
	"testing"

	"github.com/ncatelli/mockserver/pkg/router"
)

const (
	errFmt string = "want %v, got %v"
)

var (
	goodConfig = []byte(`
swagger: "2.0"
info:
  version: "1.0.0"
  title: "test"
host: "test.com"
basePath: "/"
schemes:
  - "http"
paths:
  /test/weighted/errors:
    get:
      summary: "Update an existing pet"
      description: ""
      operationId: "getTest"
      produces:
        - "application/json"
      responses:
        500:
          description: "Invalid Host"
`)
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
		Load(bytes.NewReader(goodConfig))
	})
}
