package logging

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// ErrUnknownTarget represents an error trigger when a logger target doesn't
// implement the io.Writer interface.
type ErrUnknownTarget struct {
	target interface{}
}

func (e ErrUnknownTarget) Error() string {
	return fmt.Sprintf("target %v doesn't implement io.Writer", e.target)
}

// Middleware is a logging middleware that outputs in NCSA format
type Middleware struct {
	target io.Writer
}

// Init takes 0 parameters thus this will always return nil.
func (logger *Middleware) Init(conf map[string]interface{}) error {
	if t, prs := conf["target"]; prs == true {
		v, ok := t.(io.Writer)
		if !ok {
			return ErrUnknownTarget{
				target: t,
			}
		}

		logger.target = v
	} else {
		logger.target = os.Stdout
	}

	return nil
}

// Middleware iplements the Middleware interface and executes the process of
// outputing a log before handing off the request to the next handler in the
// chain.
func (logger *Middleware) Middleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(logger.target, next)
}
