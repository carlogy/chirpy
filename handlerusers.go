package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	auth "github.com/carlogy/chirpy/internal/auth"
	"github.com/carlogy/chirpy/internal/database"
	id "github.com/google/uuid"
)

type User struct {
	Id            id.UUID   `json:"id"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	JWT_Token     string    `json:"token"`
	Refresh_Token string    `json:"refresh_token"`
	Is_chirpy_red bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, req *http.Request) {

	type jsonUser struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type errorVal struct {
		Error string `json:"error"`
	}

	body := jsonUser{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	hashpw, err := auth.HashPassword(body.Password)
	if err != nil {

		errstring := fmt.Errorf("Experienced error while hashing password, %w", err)
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(errstring.Error()))
	}

	user, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		ID:             id.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          body.Email,
		HashedPassword: string(hashpw),
	},
	)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error: creating requested user"))
		return
	}

	u := User{
		Id:            user.ID,
		Created_at:    user.CreatedAt,
		Updated_at:    user.UpdatedAt,
		Email:         user.Email,
		Is_chirpy_red: user.IsChirpyRed.Bool,
	}

	respBody, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error marshaling created user"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(respBody)
}

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {

	at, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Experienced error getting bearer token:\t%v\t", err)
		w.WriteHeader(401)
		return
	}

	userID, err := auth.ValidateJWT(at, cfg.secret)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	type updatedCreds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	uc := updatedCreds{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&uc)

	if err != nil {
		log.Printf("Experienced error unmarshalling reponse body:\t%v\n", err)
		w.WriteHeader(500)
		return
	}

	hashPW, err := auth.HashPassword(uc.Password)
	if err != nil {
		log.Printf("Error hashing pw:\t%v\n", err)
		w.WriteHeader(500)
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUserDetails(r.Context(), database.UpdateUserDetailsParams{
		UpdatedAt:      time.Now(),
		Email:          uc.Email,
		HashedPassword: hashPW,
		ID:             userID,
	})

	if err != nil {
		log.Printf("Error updated user details:\t%v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Error updating user"))
	}

	u := User{
		Id:            updatedUser.ID,
		Created_at:    updatedUser.CreatedAt,
		Updated_at:    updatedUser.UpdatedAt,
		Email:         updatedUser.Email,
		Is_chirpy_red: updatedUser.IsChirpyRed.Bool,
	}

	respBody, err := json.Marshal(u)
	if err != nil {
		log.Printf("Error marshalling repsonse body:\t%v\n", err)
		w.WriteHeader(500)
		w.Write([]byte("Experienced issue marshalling reponse body"))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
	return

}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Experienced error getting authorization header:\t%v\n", err)
		w.WriteHeader(401)
		return
	}

	if key != cfg.polka_key {
		log.Printf("Attempt to upgrage with invalid APIKey:\t%s\n", key)
		w.WriteHeader(401)
	}

	type polkaReqBody struct {
		Event string            `json:"event"`
		Data  map[string]string `json:"data"`
	}

	prb := polkaReqBody{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&prb)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error unmarshalling respose:\t%v\n", err)))
	}

	if prb.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	user_id, err := id.Parse(prb.Data["user_id"])
	if err != nil {
		log.Printf("Error parsing polka response body for uuid:\t%v\n", err)
		w.WriteHeader(500)
	}

	// respBody, err := json.Marshal(prb)
	// if err != nil {
	// 	w.WriteHeader(500)
	// 	w.Write([]byte(fmt.Sprintf("Error marshalling respose:\t%v\n", err)))
	// 	return
	// }

	_, err = cfg.dbQueries.UpgradeUserToRed(r.Context(), database.UpgradeUserToRedParams{
		UpdatedAt: time.Now(),
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ID: user_id,
	})

	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
	w.Write([]byte(""))
	// w.WriteHeader(200)
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(respBody)
}
