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
		ta := "0.0.0.0:8080"

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
}
