package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Get ID from ctx
				id, _ := ReqIDFromCtx(r.Context())

				// Logs errors
				errRec := map[string]any{
					"level":      "error",
					"request_id": id,
					"panic":      fmt.Sprint(rec),
					"stack":      string(debug.Stack()),
					"method":     r.Method,
					"path":       r.URL.Path,
				}
				_ = json.NewEncoder(os.Stdout).Encode(errRec)

				// Returns error response
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]any{"error": "internal_server_error"})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
