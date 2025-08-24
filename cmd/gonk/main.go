package main

import (
	"encoding/json"
	"net/http"

	"github.com/AJMerr/gonk/pkg/router"
)

var Version string

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type:", "application/json, charset-utf8")

	status := struct {
		OK bool `json:"ok"`
	}{
		OK: true,
	}

	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, "{ok: false}", http.StatusInternalServerError)
	}
}

func main() {
	r := router.NewRouter()

	r.Handle("GET /healthz", http.HandlerFunc(healthzHandler))

	http.ListenAndServe(":8080", r)
}
