package router

import (
	"fmt"
	"net/http"
)

// Handler includes all the metadata to decide on and serve a response.
type Handler struct {
	Weight          int               `yaml:"weight"`
	ResponseHeaders map[string]string `yaml:"response_headers"`
	StaticResponse  string            `yaml:"static_response"`
	ResponseStatus  int               `yaml:"response_status"`
	ResponsePath    string            `yaml:"response_path"`
}

// ServeHTTP implements the http.Handler interface eventually serving a request.
func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok") // STUB
}
