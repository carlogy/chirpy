package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	auth "github.com/carlogy/chirpy/internal/auth"
	"github.com/carlogy/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type jsonUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	body := jsonUser{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&body)

	if err != nil {
		errfmt := fmt.Errorf("Experienced the following error while unmarshaling response body %w", err)
		w.WriteHeader(500)
		w.Write([]byte(errfmt.Error()))
	}

	// if body.Password != "" || body.Email == "" {
	// 	w.WriteHeader(401)
	// 	w.Write([]byte("empty: Incorrect email or password"))
	// 	return
	// }

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), body.Email)

	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("query user: Incorrect email or password"))
		return
	}

	err = auth.CheckPasswordHash(body.Password, user.HashedPassword)

	if err != nil {
		log.Print(err)
		w.WriteHeader(401)
		w.Write([]byte("check Hash: Incorrect email or password"))
		return
	}

	jwt_token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(time.Hour))
	rtoken, err := auth.MakeRefreshToken()

	rt, err := cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     rtoken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	if err != nil {
		log.Print(err)
		w.WriteHeader(401)
		w.Write([]byte("Unable to generate token"))
	}

	u := User{
		Id:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Email:         user.Email,
		JWT_Token:     jwt_token,
		Refresh_Token: rt.Token,
		Is_chirpy_red: user.IsChirpyRed.Bool,
	}

	respBody, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error marshaling user"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respBody)

}
