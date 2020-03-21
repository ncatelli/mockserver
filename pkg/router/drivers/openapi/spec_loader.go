package openapi

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ncatelli/mockserver/pkg/router"
)

// Load takes an io.Reader and attempts to unmarshal the configuration into a
// route slice. On success, a slice of routes and nil is returned, otherwise an
// error is returned.
func Load(data io.Reader) ([]*router.Route, error) {
	routes := make([]*router.Route, 0)
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return routes, err
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(b)
	if err != nil {
		return routes, err
	}

	fmt.Print(swagger.Info.Description)

	return routes, nil
}
