package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
)

const port = "8080"
const filePathRoot = "."

func (db *DB) ValidateHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	// Decode request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	// Convert JSON to struct
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dirtyBody := params.Body
	dirtyArray := strings.Split(dirtyBody, " ")

	for i := range dirtyArray {
		switch strings.ToLower(dirtyArray[i]) {
		case "kerfuffle", "sharbert", "fornax":
			dirtyArray[i] = "****"
		}
	}

	cleanBody := strings.Join(dirtyArray, " ")

	// Create chirp and save to database
	responseBody, err := db.CreateChirp(cleanBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, responseBody)
}

func (db *DB) getData(w http.ResponseWriter, req *http.Request) {
	data, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
		return
	}

	respondWithJSON(w, http.StatusOK, data)
}
func main() {

	// http.NewServeMux() creates a server multiplexer
	mux := http.NewServeMux()

	// Create Database
	db := DB{
		path: "database.json",
		lock: &sync.RWMutex{},
	}

	db.ensureDB()

	// Counters
	cfg := apiConfig{
		fileServerHits: 0,
	}

	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	// If no method is mentioned the it responds to all the methods
	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("/api/reset", cfg.Reset)
	mux.HandleFunc("POST /api/chirps", db.ValidateHandler)
	mux.HandleFunc("GET /api/chirps", db.getData)

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
