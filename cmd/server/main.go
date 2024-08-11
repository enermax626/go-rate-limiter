package main

import (
	"encoding/json"
	"fmt"
	"github.com/enermax626/go-ratelimiter/limiter"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	rateLimiter := limiter.InitializeRateLimiters()

	router.Use(rateLimiter.Middleware)
	router.Get("/hello", HelloHandler)

	serverAddr := fmt.Sprintf("0.0.0.0:8080")
	srv := &http.Server{
		Handler: router,
		Addr:    serverAddr,
	}
	log.Fatal(srv.ListenAndServe(), router)
}

type HelloResponse struct {
	Hello string `json:"hello"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	status := HelloResponse{
		Hello: "world",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
