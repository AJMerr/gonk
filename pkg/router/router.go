package router

import (
	"net/http"
)

// Shape of the Router takes in a pointer to http.ServeMux
type Router struct {
	mux *http.ServeMux
}

// New Router calls to a NewServeMux for easy method safe routing
func NewRouter() *Router {
	mux := http.NewServeMux()
	return &Router{
		mux: mux,
	}
}
