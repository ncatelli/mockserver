package router

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type templateVariables struct {
	Request  *http.Request
	PathVars map[string]string
}

// Handler includes all the metadata to decide on and serve a response.
type Handler struct {
	Weight          int               `yaml:"weight"`
	ResponseHeaders map[string]string `yaml:"response_headers"`
	StaticResponse  string            `yaml:"static_response"`
	ResponseStatus  int               `yaml:"response_status"`
	ResponsePath    string            `yaml:"response_path"`
	bodyTemplate    *template.Template
}

// getBodyTemplate will attempt to retrieve, preferably from a cache field, the
// template used to generate the response body of a Handler.
func (handler *Handler) getBodyTemplate() (*template.Template, error) {
	// short circut if the template is cached
	if handler.bodyTemplate != nil {
		return handler.bodyTemplate, nil
	}

	// placeholder for future template data.
	var body string

	// static response is highest precedence
	if len(handler.StaticResponse) > 0 {
		body = handler.StaticResponse
	} else if len(handler.ResponsePath) > 0 {
		bb, err := ioutil.ReadFile(handler.ResponsePath)
		if err != nil {
			return nil, err
		}

		body = string(bb)
	}

	t, err := template.New("").Parse(body)
	if err != nil {
		return nil, err
	}

	handler.bodyTemplate = t
	return t, nil
}

// ServeHTTP implements the http.Handler interface eventually serving a request.
func (handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := handler.getBodyTemplate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for h, v := range handler.ResponseHeaders {
		w.Header().Set(h, v)
	}

	w.WriteHeader(handler.ResponseStatus)
	t.Execute(w, &templateVariables{
		Request:  r,
		PathVars: mux.Vars(r),
	})
}
