package main

import (
	"encoding/json"
	"net/http"

	"github.com/AJMerr/gonk/pkg/middleware"
	"github.com/AJMerr/gonk/pkg/router"
)

var Version string

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type:", "application/json, charset=utf-8")

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
	/* Define your CORS policy (cfg)
	cfg := middleware.CORSConfig{
		AllowedOrigins:   []string{"http://localhost:5173", "https://app.example.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id", "ETag"},
		AllowCredentials: true,
		MaxAge:           10 * time.Minute,
	} */

	r := router.NewRouter()

	r.Use(middleware.Recover)
	r.Use(middleware.ReqID)
	r.Use(middleware.Logger)
	//	r.Use(middleware.CORS(cfg))

	r.GET("/healthz", healthzHandler)
	r.GET("/panic", func(w http.ResponseWriter, r *http.Request) { panic("AAAAAHHH, BEEESSS") })

	http.ListenAndServe(":8080", r)
}
