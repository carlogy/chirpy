package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/carlogy/chirpy/internal/auth"
	"github.com/carlogy/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	rt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Received err while getting bearer token:\t%v", err)
		w.WriteHeader(401)
	}

	revokeToken, err := cfg.dbQueries.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: rt,
	})
	if err != nil {
		log.Printf("Experienced error while revoking access token:\t%v", err)
		w.WriteHeader(500)
	}

	emptyRow := database.RefreshToken{}
	if revokeToken == emptyRow {
		log.Printf("Received empty row back from sql query")
	}

	w.WriteHeader(204)
}
