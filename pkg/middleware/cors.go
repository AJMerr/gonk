package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

func nilIfEmpty(src, fallback []string) []string {
	if len(src) == 0 {
		return fallback
	}
	return src
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	allowedSet := make(map[string]struct{}, len(cfg.AllowedOrigins))
	allowAll := false
	for _, o := range cfg.AllowedOrigins {
		if o == "*" {
			allowAll = true
			continue
		}
		allowedSet[o] = struct{}{}
	}

	allowedMethods := strings.Join(nilIfEmpty(cfg.AllowedMethods, []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}), ",")
	allowedHeaders := strings.Join(nilIfEmpty(cfg.AllowedHeaders, []string{"Content-Type", "Authorization", "X-Request-Id"}), ",")
	exposedHeaders := strings.Join(nilIfEmpty(cfg.ExposedHeaders, []string{"ETag", "X-Request-Id"}), ",")

	// Checks if an origin is allowed
	originAllowed := func(origin string) bool {
		if origin == "" {
			return false
		}
		if allowAll {
			return true
		}
		_, ok := allowedSet[origin]
		return ok
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			origin := r.Header.Get("Origin")

			w.Header().Add("Vary", "Origin")
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			// If there's no Origin header, it's not a CORS request—just pass through.
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			allowed := originAllowed(origin)

			// Determine which value to use for Access-Control-Allow-Origin.
			allowOriginValue := ""
			if allowed {
				if cfg.AllowCredentials && allowAll {
					allowOriginValue = origin
				} else if allowAll {
					allowOriginValue = "*"
				} else {
					allowOriginValue = origin
				}
			}

			// Handle preflight requests (OPTIONS + Access-Control-Request-Method).
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				if allowed {
					w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
					w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

					if exposedHeaders != "" {
						w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
					}

					if allowOriginValue != "" {
						w.Header().Set("Access-Control-Allow-Origin", allowOriginValue)
					}

					if cfg.AllowCredentials {
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					}

					if cfg.MaxAge > 0 {
						w.Header().Set("Access-Control-Max-Age", strconv.FormatInt(int64(cfg.MaxAge/time.Second), 10))
					}
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			if allowed {
				if allowOriginValue != "" {
					w.Header().Set("Access-Control-Allow-Origin", allowOriginValue)
				}
				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				if exposedHeaders != "" {
					w.Header().Set("Access-Control-Expose-Headers", exposedHeaders)
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
