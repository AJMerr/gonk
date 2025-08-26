package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Shape for the Req ID Key
// Prevents overrwrites
type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

// Creates a new Req ID as a hex encoded string
func newReqID() string {
	buf := make([]byte, 16)

	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 16)
	}
	return hex.EncodeToString(buf)
}

// Retrieve the Req ID from ctx
func ReqIDFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(requestIDKey)
	s, ok := v.(string)
	return s, ok
}

// Stores a Req ID in ctx
func storeReqID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// Ensure that every req has an X-Request-ID
// - Uses incoming Req ID if present
// - Generates one otherwise
// - Sets response header
func ReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reads the incoming header, checks for ReqID
		id := r.Header.Get("X-Request-Id")

		// Checks the len of the incoming ID
		// >128, creates a new ReqID
		if len(strings.TrimSpace(id)) == 0 || len(id) > 128 {
			id = newReqID()
		}

		// Make it visible to client and downstream middleware
		w.Header().Set("X-Request-Id", id)
		r = r.WithContext(storeReqID(r.Context(), id))

		next.ServeHTTP(w, r)
	})
}
