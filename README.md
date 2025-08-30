# Gonk
A tiny, ergonomic HTTP framework for Go — built on top of the standard library.

Gonk wraps `net/http` and Go 1.22’s enhanced `ServeMux` with clean routing, a simple middleware pipeline, and JSON helpers.

## Features:

### Ergonomic routing
- `r.GET("/todos/{id}", handler)` using Go 1.22 patterns (`"METHOD /path"` and `PathValue`).

### Middleware pipeline
- `Use(middleware)` with first registered = outermost at request time.

### Batteries included middleware

- Request ID: generates/propagates `X-Request-Id` and puts it in context.

- Logger: JSON logs (method, path, status, bytes, latency, request_id).

- Recovery: catches panics → logs → `500` JSON (no stack leak).

- CORS: configurable allow/reflect behavior, preflight handling.

### JSON helpers
- `WriteJSON`, `WriteError`, `DecodeJSON` (strict: unknown fields rejected).

### Standard library first 
- No external deps!

## Installation: 
```
go get github.com/AJMerr/gonk
```
**NOTE** Requires Go 1.22+

## Example Quick Start:
```
package main

import (
	"net/http"
	"time"

	"github.com/AJMerr/gonk/pkg/middleware"
	"github.com/AJMerr/gonk/pkg/router"
	"github.com/AJMerr/gonk/pkg/jsonutil"
)

func main() {
	r := router.NewRouter()

	// Middleware (first Use is outermost)
	r.Use(middleware.Recover)      // catch panics → 500
	r.Use(middleware.ReqID)        // add X-Request-Id + context value
	r.Use(middleware.Logger)       // JSON request logs
	r.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET","POST","PUT","PATCH","DELETE","OPTIONS"},
		AllowedHeaders:   []string{"Content-Type","Authorization","X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id","ETag"},
		AllowCredentials: true,
		MaxAge:           10 * time.Minute,
	}))

	// Routes
	r.GET("/healthz", func(w http.ResponseWriter, r *http.Request) {
		jsonutil.WriteJSON(w, http.StatusOK, map[string]string{"status":"ok"})
	})

	http.ListenAndServe(":8080", r)
}
```

## Routing:
Gonk delegates to Go's `ServeMux` but adds methood helpers for an ergonomical approach
```
r.GET("/todos", listTodos)
r.POST("/todos", createTodo)
r.GET("/todos/{id}", getTodo)         // use r.PathValue("id")
r.PUT("/todos/{id}", updateTodo)
r.DELETE("/todos/{id}", deleteTodo)
```

### Inside a Handler:
```
func getTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // from /todos/{id}
	// ...
}
```

## Middleware:
Gonk’s middleware signature is the usual `func(next http.Handler) http.Handler`.
**Order matters**: the first Use wraps everything else (outermost).
```
r.Use(middleware.Recover)
r.Use(middleware.ReqID)
r.Use(middleware.Logger)
r.Use(middleware.CORS(cfg))
```

## Request Id: 
- Response header: `X-Request-Id`
- Context helper: 
```
id, _ := middleware.ReqIDFromCtx(r.Context())
```
## Logger:
Single JSON line per request:
```
{
  "level": "info",
  "request_id": "abc123",
  "method": "GET",
  "path": "/healthz",
  "status": 200,
  "bytes": 12,
  "latency_ms": 1,
  "remote_ip": "127.0.0.1:12345",
  "user_agent": "curl/8.0.1",
  "ts": "2025-08-30T12:34:56.789Z"
}
```

## Panic Recovery:
Catches panics, logs `panic` + optional stack (logs only), returns:
`{"error":"internal_server_error"}`
With a `500` error code

## CORS:
```
r.Use(middleware.CORS(middleware.CORSConfig{
  AllowedOrigins:   []string{"http://localhost:5173","https://app.example.com"},
  AllowedMethods:   []string{"GET","POST","PUT","PATCH","DELETE","OPTIONS"},
  AllowedHeaders:   []string{"Content-Type","Authorization","X-Request-Id"},
  ExposedHeaders:   []string{"X-Request-Id","ETag"},
  AllowCredentials: true,          // if true, framework reflects allowed origin (not "*")
  MaxAge:           10*time.Minute // preflight cache
}))
```
Handles both **simple** and **preflight** requests, sets `Vary` headers appropriately.

## JSON Utilities
```
// Write a JSON response
jsonutil.WriteJSON(w, 201, data)

// Consistent error shape
jsonutil.WriteError(w, 400, "invalid_json")

// Strict decode: unknown fields rejected, trailing data rejected
var in CreateTodo
if err := jsonutil.DecodeJSON(r, &in); err != nil {
    code := jsonutil.StatusFromDecodeError(err)
    jsonutil.WriteError(w, code, "invalid_json", err.Error())
    return
}
```
Sentinel errors & status mapping:

- `400`: empty body, syntax error, unknown field, multiple JSON values
- `413`: body too large (via http.MaxBytesReader)
- `422`: type mismatch
- `500`: fallback

## Project Layout:
```
pkg/
  router/
    router.go           # Router, NewRouter, Handle, ServeHTTP, GET/POST/PUT/PATCH/DELETE, Use
  middleware/
    requestid.go        # ReqID, ReqIDFromCtx
    logging.go          # Logger (JSON)
    recovery.go         # Recover
    cors.go             # CORS + CORSConfig
  jsonutil/
    json.go             # WriteJSON, WriteError, DecodeJSON, StatusFromDecodeError, errors
cmd/
  gonk/
    main.go             # example app / playground (optional)
```