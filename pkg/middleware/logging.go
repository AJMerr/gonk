package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
	wrote  bool
}

func (sw *statusWriter) WriteHeader(code int) {
	if !sw.wrote {
		sw.status = code
		sw.wrote = true
	}
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	if !sw.wrote {
		sw.status = 200
		sw.wrote = true
	}
	n, err := sw.ResponseWriter.Write(b)
	sw.bytes += n
	return n, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(wrapped, r)
		id, _ := ReqIDFromCtx(r.Context())
		rec := map[string]any{
			"ts":         time.Now().Format(time.RFC3339Nano),
			"level":      "info",
			"request_id": id,
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     wrapped.status,
			"bytes":      wrapped.bytes,
			"latency_ms": time.Since(start).Milliseconds(),
			"remote_ip":  r.RemoteAddr,
			"user_agent": r.UserAgent(),
		}
		_ = json.NewEncoder(os.Stdout).Encode(rec)
	})
}
