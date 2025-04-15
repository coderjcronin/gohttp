package main

import (
	"strings"
)

/*
func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
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
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: profanityCheck(params.Body),
	})

} */

func profanityCheck(rawBody string) string {
	words := strings.Fields(rawBody)
	profane := []string{"kerfuffle", "sharbert", "fornax"}
	censor := "****"
	returnStr := []string{}

	for _, word := range words {
		returnWord := word
		for _, swear := range profane {
			if strings.ToLower(word) == swear {
				returnWord = censor
				break
			}
		}
		returnStr = append(returnStr, returnWord)
	}

	return strings.Join(returnStr, " ")
}
