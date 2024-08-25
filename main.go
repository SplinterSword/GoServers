package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const port = "8080"
const filePathRoot = "."

func (db *DB) ValidateHandler(w http.ResponseWriter, r *http.Request) {

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

	responseBody, _ := db.CreateChirp(cleanBody)

	respondWithJSON(w, 201, responseBody)
}

func (db *DB) getData(w http.ResponseWriter, req *http.Request) {
	data, _ := db.GetChirps()
	respondWithJSON(w, 200, data)
}

func main() {

	// http.NewServeMux() creates a server multiplexer
	mux := http.NewServeMux()

	// Create Database
	db, _ := NewDB("database.json")

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
