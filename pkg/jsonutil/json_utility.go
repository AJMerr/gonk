package jsonutil

import (
	"encoding/json"
	"net/http"
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
