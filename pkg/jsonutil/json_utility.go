package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Sets http response, header to JSON
func WriteJSON(w http.ResponseWriter, code int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, `{"error":"encoding_failed}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(b)
	_, _ = w.Write([]byte("\n"))
}

func WriteError(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, code, map[string]any{"error": message})
}

var (
	ErrEmptyBody          = errors.New("request body must not be empty")
	ErrSyntax             = errors.New("malformed JSON")
	ErrUnknownField       = errors.New("request contained unknown field")
	ErrTooLarge           = errors.New("request body too large")
	ErrMultipleJSONValues = errors.New("body must contain only a single JSON object")
)

func DecodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return ErrEmptyBody
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // reject extra fields

	// Try decoding
	if err := dec.Decode(dst); err != nil {
		// Classify common error cases
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("%w at byte %d", ErrSyntax, syntaxErr.Offset)
		case errors.As(err, &unmarshalTypeErr):
			return fmt.Errorf("bad value for field %q: %w", unmarshalTypeErr.Field, err)
		case errors.Is(err, io.EOF):
			return ErrEmptyBody
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fmt.Errorf("%w: %s", ErrUnknownField, err)
		default:
			return err // propagate
		}
	}

	// Ensure there’s only a single JSON value in the body
	if dec.More() {
		return ErrMultipleJSONValues
	}

	return nil
}
