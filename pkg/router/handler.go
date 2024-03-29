package router

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/leekchan/gtf"
	"github.com/ncatelli/mockserver/pkg/router/generator"
)

type templateVariables struct {
	Request  *http.Request
	PathVars map[string]string
}

// Generate plugins
//go:generate go run ./generator/gen.go

// Handler includes all the metadata to decide on and serve a response.
type Handler struct {
	Weight          uint              `yaml:"weight"`
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
		bb, err := os.ReadFile(handler.ResponsePath)
		if err != nil {
			return nil, err
		}

		body = string(bb)
	}

	t, err := template.New("").Funcs(gtf.GtfFuncMap).Funcs(generator.PluginsFuncMap()).Parse(body)
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
		http.Error(w, "", http.StatusInternalServerError)
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
