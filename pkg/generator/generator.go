package generator

// Generator defines the interface for types to be added to the templating
// FuncMap and exposed to the handlers.
type Generator interface {
	ID() string
	Generate(...interface{}) interface{}
}
