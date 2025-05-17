package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"slices"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	bodyInParts := strings.Split(params.Body, " ")
	paramsInLowerCase := strings.Split(strings.ToLower(params.Body), " ")
	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range paramsInLowerCase {
		if slices.Contains(bannedWords, word) {
			bodyInParts[i] = "****"
		}
	}
	cleanedBody := strings.Join(bodyInParts, " ")

	respondWithJSON(w, http.StatusOK, returnBody{
		CleanedBody: cleanedBody,
	})
}
