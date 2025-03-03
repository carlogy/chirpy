package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/carlogy/chirpy/internal/auth"
	"github.com/carlogy/chirpy/internal/database"
	id "github.com/google/uuid"
)

type Chirp struct {
	ID         id.UUID   `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	UserID     id.UUID   `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {

	token, err := auth.GetBearerToken(req.Header)

	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		w.Write([]byte("Error getting bearer token"))
	}

	validId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Print(err)
		w.WriteHeader(401)
		w.Write([]byte("Invalid token provided"))
	}

	type chirp struct {
		Body string `json:"body"`
		// User_Id id.UUID `json:"user_id"`
	}

	type errorJSON struct {
		Error string `json:"error"`
	}

	c := chirp{}

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal("Error decoding body: %w", err)
		w.WriteHeader(500)
		w.Write([]byte("Error while decoding json body"))
	}

	if len(c.Body) > 140 {

		e := errorJSON{
			Error: "Chirp is too long",
		}

		respBody, err := json.Marshal(e)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error marshalling error json"))
			return

		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(respBody)
		return

	}

	badwords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"formax":    {},
	}

	c.Body = scrubBody(c.Body, badwords)

	createdChirp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		ID:        id.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      c.Body,
		UserID:    validId,
	})

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error saving chirp"))
	}

	jsonChirp := convertToJSONChirp(createdChirp)

	respBody, err := json.Marshal(jsonChirp)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)

}

func scrubBody(body string, badwords map[string]struct{}) string {
	words := strings.Split(body, " ")

	for i, word := range words {
		if _, ok := badwords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}
	scrubbed := strings.Join(words, " ")
	return scrubbed
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	sort, sortOk := queryParams["sort"]
	authorID, authorIDOk := queryParams["author_id"]

	switch {
	case authorIDOk && sortOk && strings.ToUpper(sort[0]) == "DESC":
		userId, err := id.Parse(authorID[0])
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("No author"))
		}

		chirps, err := cfg.dbQueries.GetChirpsByAuthorDESC(r.Context(), userId)

		jsonChirps := make([]Chirp, len(chirps))

		for i, chirp := range chirps {
			jsonChirps[i] = convertToJSONChirp(chirp)
		}

		respBody, err := json.Marshal(jsonChirps)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error marshaling chirp list"))
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
		return

	case authorIDOk && !sortOk:
		userId, err := id.Parse(authorID[0])
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("No author"))
		}

		chirps, err := cfg.dbQueries.GetChirpsByAuthor(r.Context(), userId)

		jsonChirps := make([]Chirp, len(chirps))

		for i, chirp := range chirps {
			jsonChirps[i] = convertToJSONChirp(chirp)
		}

		respBody, err := json.Marshal(jsonChirps)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error marshaling chirp list"))
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
		return

	case !authorIDOk && strings.ToUpper(sort[0]) == "DESC":
		chirps, err := cfg.dbQueries.GetChirpsDesc(r.Context())
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error get chirps"))
		}

		jsonChirps := make([]Chirp, len(chirps))

		for i, chirp := range chirps {
			jsonChirps[i] = convertToJSONChirp(chirp)
		}

		respBody, err := json.Marshal(jsonChirps)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error marshaling chirp list"))
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
	default:
		chirps, err := cfg.dbQueries.GetChirps(r.Context())
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error get chirps"))
		}

		jsonChirps := make([]Chirp, len(chirps))

		for i, chirp := range chirps {
			jsonChirps[i] = convertToJSONChirp(chirp)
		}

		respBody, err := json.Marshal(jsonChirps)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error marshaling chirp list"))
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)

	}
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	cid := r.PathValue("chirpID")

	if cid == "" {
		log.Print("No id was provided to single call")

		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{ "Error": "No id for chirp was provided"}`))
		return
	}

	ChirpID, err := id.Parse(cid)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	c, err := cfg.dbQueries.GetChirp(r.Context(), ChirpID)
	if err != nil {

		w.WriteHeader(404)
		w.Header().Set("Content-Type", "applcation/json")
		w.Write([]byte(`{ "Error": "No chirp found"}`))
		return
	}

	jsonChirp := convertToJSONChirp(c)

	respBody, err := json.Marshal(jsonChirp)
	if err != nil {
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{ "Error": "Error marshaling chirp"}`))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)

}

func convertToJSONChirp(c database.Chirp) Chirp {

	return Chirp{
		ID:         c.ID,
		Created_at: c.CreatedAt,
		Updated_at: c.UpdatedAt,
		Body:       c.Body,
		UserID:     c.UserID,
	}

}

func (cfg *apiConfig) handlerDeletechirp(w http.ResponseWriter, r *http.Request) {

	cid, err := id.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing chirpID:\t%v\n", err)
		w.WriteHeader(401)
		return
	}

	at, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token:\t%v\n", err)
		w.WriteHeader(401)
		return
	}

	userID, err := auth.ValidateJWT(at, cfg.secret)
	if err != nil {
		log.Printf("Error validating token:\t%v\n", err)
		w.WriteHeader(401)
		return
	}

	c, err := cfg.dbQueries.GetChirp(r.Context(), cid)

	if c.UserID != userID {
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Error": Not auther of chirp unable to delete}`))
		return
	}

	_, err = cfg.dbQueries.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     cid,
		UserID: c.UserID,
	})

	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(204)
	return

}
