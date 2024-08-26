package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

const port = "8080"
const filePathRoot = "."

func (db *DB) MakeUser(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body     string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	data, err := db.CreateMail(params.Body, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create a new User")
		return
	}

	type responseStructure struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	response := responseStructure{
		Id:    data.Id,
		Email: data.Email,
	}

	respondWithJSON(w, http.StatusOK, response)

	respondWithJSON(w, http.StatusCreated, response)
}

func (db *DB) HandleLogin(w http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body     string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	AllUsers, err := db.GetMails()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to Get all users User")
		return
	}

	var idNumber int

	for i := 0; i < len(AllUsers); i++ {
		if AllUsers[i].Email == params.Body {
			err := bcrypt.CompareHashAndPassword(AllUsers[i].Password, []byte(params.Password))

			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Invalid Password")
				return
			}

			idNumber = AllUsers[i].Id
			break
		}
	}

	type responseStructure struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	response := responseStructure{
		Id:    idNumber,
		Email: params.Body,
	}

	respondWithJSON(w, http.StatusOK, response)
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
	mux.HandleFunc("GET /api/chirps/{chirpID}", db.getSpecificData)
	mux.HandleFunc("POST /api/users", db.MakeUser)
	mux.HandleFunc("POST /api/login", db.HandleLogin)

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
