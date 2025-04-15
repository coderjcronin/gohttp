package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/coderjcronin/gohttp/internal/auth"
	"github.com/coderjcronin/gohttp/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) apiCheckLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.QueryHashedPassword(r.Context(), params.Email)
	if err != nil {
		log.Printf("Failed to look up hashed pass: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to look up hashed pass", err)
		return
	}

	err = auth.CheckPassword(user.HashedPassword, params.Password)
	if err != nil {
		log.Printf("Failed to match hashed password: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	expires := int(cfg.expires.Seconds())
	if params.ExpiresInSeconds != 0 && params.ExpiresInSeconds < expires {
		expires = params.ExpiresInSeconds
	}

	issuedToken, err := auth.MakeJWT(user.ID, cfg.secret, (time.Duration(expires) * time.Second))
	if err != nil {
		log.Printf("Failed to issue token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Failed to issue token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Id:        user.ID,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Email:     user.Email,
		Token:     issuedToken,
	})
}

func (cfg *apiConfig) apiAddUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password.", err)
		return
	}

	args := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	u, err := cfg.db.CreateUser(r.Context(), args)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user.", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		Id:        u.ID,
		CreatedAt: u.CreatedAt.String(),
		UpdatedAt: u.UpdatedAt.String(),
		Email:     u.Email,
	})

}
