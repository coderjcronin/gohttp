package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/coderjcronin/gohttp/internal/auth"
	"github.com/coderjcronin/gohttp/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body  string `json:"body"`
		Token string `json:"token"`
	}
	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("JSON decode error: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//Check token
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Failed to check for bearer token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not find bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		log.Printf("Failed to validate token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not validate token", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}

	ch, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   profanityCheck(params.Body),
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp.", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		Id:        ch.ID,
		CreatedAt: ch.CreatedAt.String(),
		UpdatedAt: ch.UpdatedAt.String(),
		Body:      ch.Body,
		UserId:    ch.UserID,
	})
}

func (cfg *apiConfig) apiGetChirps(w http.ResponseWriter, r *http.Request) {
	type returnChirp struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	ch, err := cfg.db.RetrieveAllChirps(r.Context())
	if err != nil {
		log.Printf("Error retrieve all chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps", err)
		return
	}

	returnSlice := []returnChirp{}

	for _, chirp := range ch {
		appendChirp := returnChirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt.String(),
			UpdatedAt: chirp.UpdatedAt.String(),
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
		returnSlice = append(returnSlice, appendChirp)
	}

	respondWithJSON(w, http.StatusOK, returnSlice)
}

func (cfg *apiConfig) apiGetChirp(w http.ResponseWriter, r *http.Request) {
	type returnChirp struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Failed to parse UUID for chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to parse chirp UUID for lookup.", err)
		return
	}

	ch, err := cfg.db.RetrieveSelectChirp(r.Context(), chirpID)
	if err == sql.ErrNoRows {
		log.Printf("No rows found, returning 404.")
		respondWithError(w, http.StatusNotFound, "Chirp not found.", err)
		return
	} else if err != nil {
		log.Printf("Failed to query database chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to query databse.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnChirp{
		Id:        ch.ID,
		CreatedAt: ch.CreatedAt.String(),
		UpdatedAt: ch.UpdatedAt.String(),
		Body:      ch.Body,
		UserId:    ch.UserID,
	})
}
