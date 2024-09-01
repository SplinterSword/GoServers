package main

import (
	"net/http"

	"github.com/SplinterSword/WebServers/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	err = cfg.DB.RevokeRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
