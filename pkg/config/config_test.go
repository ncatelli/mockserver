package config

import (
	"os"
	"testing"
)

const (
	errFmt string = "want %v, got %v"
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
}
