package generator

import "math/rand"

// Generator defines the interface for types to be added to the templating
// FuncMap and exposed to the handlers.
type Generator interface {
	Init(rand.Rand) error
	ID() string
	Generate(...string) string
}
