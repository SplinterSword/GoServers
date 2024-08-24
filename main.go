package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, cfg.fileServerHits)
	w.Write([]byte(hits))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
		// the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	// Converts json to struct and put it in the mentioned struct address
	err := decoder.Decode(&params)

	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		respondWithError(w, 500, err.Error())
		return
	}

	dirtyBody := params.Body

	dirtyArray := strings.Split(dirtyBody, " ")

	for i := 0; i < len(dirtyArray); i++ {
		if strings.ToLower(dirtyArray[i]) == "kerfuffle" || strings.ToLower(dirtyArray[i]) == "sharbert" || strings.ToLower(dirtyArray[i]) == "fornax" {
			dirtyArray[i] = "****"
		}
	}

	cleanBody := strings.Join(dirtyArray, " ")

	respondWithJSON(w, 200, map[string]string{"cleaned_body": cleanBody})
}

func main() {

	// http.NewServeMux() creates a server multiplexer
	mux := http.NewServeMux()
	cfg := apiConfig{
		fileServerHits: 0,
	}
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	// If no method is mentioned the it responds to all the methods
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("/api/reset", cfg.Reset)
	mux.HandleFunc("POST /api/validate_chirp", ValidateHandler)

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
