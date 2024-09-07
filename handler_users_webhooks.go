package main

import (
	"encoding/json"
	"net/http"

	"github.com/SplinterSword/WebServers/internal/auth"
)

type UserID struct {
	UserID int `json:"user_id"`
}

func (cfg *apiConfig) handlerUsersWebHook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  UserID `json:"data"`
	}

	PolkaKey, err := auth.GetPolkaKey(r.Header)
	if err != nil || PolkaKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	params := parameters{}
	Decoder := json.NewDecoder(r.Body)
	err = Decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrive the body of the request")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id := params.Data.UserID

	_, err = cfg.DB.UpdateUserStatus(id)
	if err != nil {
		if err.Error() == "user already using chirpy red" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithError(w, http.StatusNotFound, err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
