package ncsa

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// Logger is a logging middleware that outputs in NCSA format
type Logger struct{}

// Init takes 0 parameters thus this will always return nil.
func (ncsa *Logger) Init(conf map[string]interface{}) error {
	return nil
}

// Middleware iplements the Middleware interface and executes the process of
// outputing a log before handing off the request to the next handler in the
// chain.
func (ncsa *Logger) Middleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
