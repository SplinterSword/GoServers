package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "8080"
const filePathRoot = "."

type apiConfig struct {
	fileServerHits int
}

// This is a middleware
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf("Hits: %d", cfg.fileServerHits)
	w.Write([]byte(hits))
}

func main() {

	// http.NewServeMux() creates a server multiplexer
	mux := http.NewServeMux()
	cfg := apiConfig{
		fileServerHits: 0,
	}
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("/healthz", handleReadiness)
	mux.HandleFunc("/metrics", cfg.handleMetrics)
	mux.HandleFunc("/reset", cfg.Reset)

	// use to create server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Print in terminal
	log.Printf("Serving files from %s on http://localhost:%s/app\n", filePathRoot, port)

	// catch error in the terminal
	log.Fatal(srv.ListenAndServe())
}
