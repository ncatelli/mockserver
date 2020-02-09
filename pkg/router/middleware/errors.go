package middleware

import "fmt"

// ErrUndefinedMiddleware represents lookup against a middleware that has yet
// to be defined has failed.
type ErrUndefinedMiddleware struct {
	ID string
}

func (e ErrUndefinedMiddleware) Error() string {
	return fmt.Sprintf("the middleware %s is undefined", e.ID)
}
