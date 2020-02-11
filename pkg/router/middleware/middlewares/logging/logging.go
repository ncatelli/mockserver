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
	target string
}

func (e ErrUnknownTarget) Error() string {
	return fmt.Sprintf("target %s unknown", e.target)
}

// Middleware is a logging middleware that outputs in CLF format
type Middleware struct {
	target io.Writer
}

// Init takes a configuration map of strings to configure the middleware. The
// current only accepted parameter is "target" representing the output target
// for log data. Currently, "stdout" is the only supported target though this
// is subject ot change in the future. If no value is specified, "stdout" is
// default to for the target.
func (logger *Middleware) Init(conf map[string]string) error {
	if t, prs := conf["target"]; prs == true {
		switch t {
		case "stdout":
			logger.target = os.Stdout
		default:
			return ErrUnknownTarget{
				target: t,
			}
		}
	} else {
		// default to os.Stdout if no target is specified.
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
