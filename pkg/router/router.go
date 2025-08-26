package router

import (
	"net/http"
)

// Shape of the Router takes in a pointer to http.ServeMux
type Router struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
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
		panic("Router is not initialized, use NewRouter()")
	}
	handler := http.Handler(r.mux)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	handler.ServeHTTP(w, req)
}

// Allows you to handle a GET request route
// Example: r.GET("/healthz", healthzHandler)
func (r *Router) GET(p string, fn func(w http.ResponseWriter, r *http.Request)) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}
	handler := http.HandlerFunc(fn)

	// Adds GET to path
	fullP := "GET " + p

	r.Handle(fullP, handler)
}

// Allows you to handle a POST request route
// Example: r.POST("/note", noteHandler)
func (r *Router) POST(p string, fn func(w http.ResponseWriter, r *http.Request)) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}

	handler := http.HandlerFunc(fn)

	// Adds POST to path
	fullP := "POST " + p

	r.Handle(fullP, handler)
}

// Allows you to handle a PATCH request route
// Example: r.PATCH("/note/{id}", noteHandler)
func (r *Router) PATCH(p string, fn func(w http.ResponseWriter, r *http.Request)) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}

	handler := http.HandlerFunc(fn)

	fullP := "PATCH " + p

	r.Handle(fullP, handler)
}

// Allows you to handle a PUT request route
// Example: r.PUT("/note/{id}", noteHandler)

func (r *Router) PUT(p string, fn func(w http.ResponseWriter, r *http.Request)) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}

	handler := http.HandlerFunc(fn)

	fullP := "UPDATE " + p

	r.Handle(fullP, handler)
}

// Allows you to handle a DELETE request route
// Example: r.DELETE("/note/{id}", noteHandler)
func (r *Router) DELETE(p string, fn func(w http.ResponseWriter, r *http.Request)) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}

	handler := http.HandlerFunc(fn)

	fullP := "DELETE " + p

	r.Handle(fullP, handler)
}

// Allows you to handle a HEAD request route
func (r *Router) HEAD(p string, fn func(w http.ResponseWriter, r *http.Request)) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}

	handler := http.HandlerFunc(fn)

	fullP := "HEAD " + p

	r.Handle(fullP, handler)
}

// Allows you to handle a OPTIONS request route
func (r *Router) OPTIONS(p string, fn func(w http.ResponseWriter, r *http.Request)) {
	if r.mux == nil {
		panic("Router is not initialized, use NewRouter()")
	}

	handler := http.HandlerFunc(fn)

	fullP := "OPTIONS " + p

	r.Handle(fullP, handler)
}

// Ensures Middleware is usable with router
func (r *Router) Use(mw func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, mw)
}
