package main

import (
	"net/http"
	"strconv"

	"github.com/SplinterSword/WebServers/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	Token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrive token")
		return
	}

	Author, err := auth.ValidateJWT(Token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Invalid User")
		return
	}

	AuthorID, err := strconv.Atoi(Author)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert AuthorID into int")
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	if AuthorID != chirpID {
		respondWithError(w, http.StatusForbidden, "Invalid User")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Couldn't Delete Chirp")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "Deleted")
}
