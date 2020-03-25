package openapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/ncatelli/mockserver/pkg/router"
)

// opMap represents a kv pairing of specOperations to their corresponding
// http method.
type opMap map[string]*spec.Operation

// op stores a method name and swagger.Operation pairing
type op struct {
	Name string
	Op   *spec.Operation
}

// newOp initializes and returns a pointer to an op after formatting its
// parameters.
func newOp(n string, o *spec.Operation) *op {
	return &op{
		Name: strings.ToUpper(n),
		Op:   o,
	}
}

// Load takes an io.Reader and attempts to unmarshal the swagger spec into a
// route slice. On success, a slice of routes and nil is returned, otherwise an
// error is returned.
func Load(data io.Reader) ([]*router.Route, error) {
	routes := make([]*router.Route, 0)
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return routes, err
	}

	// load and expand the spec
	swagger, err := loads.Analyzed(json.RawMessage(b), "")
	if err != nil {
		return routes, err
	}

	ss := swagger.Spec()
	for path, params := range ss.Paths.Paths {
		for _, operation := range getPathItemOperations(params) {
			r, err := generateRouteFromSwaggerPathOperation(path, operation)
			if err != nil {
				return routes, err
			}

			routes = append(
				routes,
				r,
			)
		}
	}

	return routes, nil
}

// Takes swagger spec parameters and attempts to package them into a
// corresponding mock server route.
func generateRouteFromSwaggerPathOperation(urlPath string, o *op) (*router.Route, error) {
	return &router.Route{
		Path:   urlPath,
		Method: o.Name,
		Handlers: []router.Handler{
			router.Handler{
				ResponseHeaders: map[string]string{
					"content-type": o.Op.Produces[0],
				},
			},
		},
	}, nil
}

// getPathItemOperations takes a PathItem and iterates over all
// Operations defined on the map, returning any non-nil Operations.
func getPathItemOperations(p spec.PathItem) []*op {
	props := p.PathItemProps
	ops := make([]*op, 0, 1)

	sV := reflect.ValueOf(props)
	sT := reflect.TypeOf(props)

	for i := 0; i < sV.NumField(); i++ {
		method := sT.Field(i).Name
		if method == "Parameters" {
			continue
		}

		sop := sV.Field(i).Interface().(*spec.Operation)
		if sop != nil {
			ops = append(ops, newOp(method, sop))
		}
	}
	return ops
}
