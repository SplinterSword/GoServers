package main

import "net/http"

func (cfg *apiConfig) Reset(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileServerHits = 0
	w.Write([]byte("Resetted"))
}
