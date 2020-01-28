package router

type Handler struct {
	Weight          int
	ResponseHeaders map[string]string
	StaticResponse  string
	ResponseStatus  int
}
