package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

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

func (db *DB) getSpecificData(w http.ResponseWriter, req *http.Request) {

	ID_Number, err := strconv.Atoi(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get the Id number")
		return
	}

	Chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
		return
	}

	Chi := Chirp{}

	for i := 0; i < len(Chirps); i++ {
		if Chirps[i].Id == ID_Number {
			Chi = Chirps[i]
			break
		}
	}

	respondWithJSON(w, http.StatusOK, Chi)
}
