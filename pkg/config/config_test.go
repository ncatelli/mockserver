package config

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

const (
	errFmt              = "want %v, got %v"
	goodTestFixturePath = "test_fixtures/good.yaml"
	goodResponseBody    = `- path: "/test"
  method: GET
  handlers:
    - weight: 1
      response_headers:
        content-type: application/json
      static_response: '{"resp": "Ok"}'
      response_status: 200
`
)

func TestInitializingAConfigShould(t *testing.T) {
	t.Run("return parse Addr field when an env is specified.", func(t *testing.T) {
		ta := "127.0.0.1:8080"

		oe := os.Getenv("ADDR")
		if oe == "" {
			defer os.Unsetenv("ADDR")
		} else {
			defer os.Setenv("ADDR", oe)
		}

		os.Setenv("ADDR", ta)
		c, err := New()
		if err != nil {
			t.Error(err)
		}

		if c.Addr != ta {
			t.Errorf(errFmt, ta, c.Addr)
		}
	})

	t.Run("return the default Addr field if no env is passed", func(t *testing.T) {
		ea := "0.0.0.0:8080"

		c, err := New()
		if err != nil {
			t.Error(err)
		}

		if c.Addr != ea {
			t.Errorf(errFmt, ea, c.Addr)
		}
	})
}

func TestConfigurationLoadingShould(t *testing.T) {
	t.Run("return an ErrUnspecifiedConfig when no config option is set", func(t *testing.T) {
		c := Config{}

		_, err := c.Load()
		if err == nil {
			t.Errorf(errFmt, "error", err)
		}
	})

	t.Run("load a configuration from a URL when specified", func(t *testing.T) {
		// generate a test server so we can capture and inspect the request
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(goodResponseBody))
		}))
		defer func() { testServer.Close() }()

		// setup Config
		testURL, _ := url.Parse(testServer.URL)
		c := Config{
			ConfigURL: *testURL,
		}

		reader, err := c.Load()
		if err != nil {
			t.Errorf(errFmt, nil, err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		s := buf.String()

		if !reflect.DeepEqual(goodResponseBody, s) {
			t.Errorf(errFmt, goodResponseBody, s)
		}
	})

	t.Run("load a configuration from a filepath when specified", func(t *testing.T) {
		c := Config{
			ConfigPath: goodTestFixturePath,
		}

		reader, err := c.Load()
		if err != nil {
			t.Errorf(errFmt, nil, err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		s := buf.String()

		if !reflect.DeepEqual(goodResponseBody, s) {
			t.Errorf(errFmt, goodResponseBody, s)
		}
	})
}
