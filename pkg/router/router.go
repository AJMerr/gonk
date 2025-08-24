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

// Handle allows you to take in a string that contains a METHOD and route, and a handler
func (r *Router) Handle(p string, h http.Handler) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}
	r.mux.Handle(p, h)
}

// Resolves the ServeHTTP requirement to use for HTTP servers.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.mux == nil {
		panic("Router is not intialized, use NewRouter()")
	}
	r.mux.ServeHTTP(w, req)
}
