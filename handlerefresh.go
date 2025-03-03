package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/carlogy/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type TokenResponse struct {
		Token string `json:"token"`
	}

	rt, err := auth.GetBearerToken(r.Header)

	if err != nil {
		log.Printf("Recieved err while getting bearer token:\t%v", err)
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
	}

	newRT, err := cfg.dbQueries.GetUserFromToken(r.Context(), rt)

	if err != nil {
		log.Printf("Error checking for refresh token:\t%v\n", err)
	}

	if newRT.RevokedAt.Valid || newRT.ExpiresAt.Before(time.Now()) {
		log.Print("Attempt to refresh revoked token\t", newRT.Token)
		w.WriteHeader(401)
		return
	}

	newAT, err := auth.MakeJWT(newRT.ID, cfg.secret, time.Duration(time.Hour))

	if err != nil {
		log.Print(err)
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	tr := TokenResponse{
		Token: newAT,
	}
	respBody, err := json.Marshal(tr)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error marshaling token"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respBody)

}
