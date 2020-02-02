package simple

import (
	"io/ioutil"
	"path/filepath"

	"github.com/ncatelli/mockserver/pkg/router"
	"gopkg.in/yaml.v2"
)

// LoadFromFile takes a path an attempts to unmarshal a route slice from a yaml
// file. On success, a slice of routes and nil is returned, otherwise an error
// is returned.
func LoadFromFile(path string) ([]*router.Route, error) {
	routes := make([]*router.Route, 0)
	dat, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return routes, err
	}

	if err = yaml.Unmarshal(dat, &routes); err != nil {
		return routes, err
	}

	return routes, nil
}
