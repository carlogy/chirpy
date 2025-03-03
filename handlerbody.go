package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandlerBody(w http.ResponseWriter, req *http.Request) {

	type bodyParams struct {
		Body string `json:"body"`
	}

	type errorVal struct {
		Error string `json:"error"`
	}

	body := bodyParams{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {

		dat := errorVal{
			Error: "Error decoding request body",
		}

		respBody, err := json.Marshal(dat)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		log.Fatal("Error decoding body: %w", err)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
		w.WriteHeader(500)

	}

	if len(body.Body) > 140 {
		dat := errorVal{
			Error: "Chirp is too long",
		}

		respBody, err := json.Marshal(dat)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(respBody)
		return
	}

	type validatedChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	words := strings.Split(body.Body, " ")
	replaceWith := "****"

	for i, word := range words {
		word = strings.ToLower(word)

		switch word {
		case "kerfuffle":
			words[i] = replaceWith
		case "sharbert":
			words[i] = replaceWith
		case "fornax":
			words[i] = replaceWith
		default:
			continue
		}

	}

	newBody := strings.Join(words, " ")

	vc := validatedChirp{CleanedBody: newBody}

	respBody, err := json.Marshal(vc)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(respBody)
}
