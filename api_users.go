package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
		Id           uuid.UUID `json:"id"`
		CreatedAt    string    `json:"created_at"`
		UpdatedAt    string    `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
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
		respondWithError(w, http.StatusInternalServerError, "Failed to look up hashed pass", err)
		return
	}

	err = auth.CheckPassword(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	issuedToken, err := auth.MakeJWT(user.ID, cfg.secret, cfg.expires)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to issue token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to issue refresh token", err)
		return
	}

	var refreshRequest = database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token in db", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt.String(),
		UpdatedAt:    user.UpdatedAt.String(),
		Email:        user.Email,
		Token:        issuedToken,
		RefreshToken: refreshToken,
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password.", err)
		return
	}

	args := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	u, err := cfg.db.CreateUser(r.Context(), args)
	if err != nil {
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

func (cfg *apiConfig) apiRefreshToken(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "No authorization header found", nil)
	}

	if !strings.Contains(authHeader, "Bearer") {
		respondWithError(w, http.StatusUnauthorized, "No bearer field found in authorization header", nil)
	}

	splitString := strings.Fields(authHeader)

	if splitString[1] == "" {
		respondWithError(w, http.StatusUnauthorized, "No token found in bearer field", nil)
	}

	refToken, err := cfg.db.RetrieveSelectRefreshToken(r.Context(), splitString[1])
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No matching token found in DB", nil)
		return
	}

	if refToken.ExpiresAt.Before(time.Now()) || refToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	newToken, err := auth.MakeJWT(refToken.UserID, cfg.secret, cfg.expires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refreshed token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Token: newToken,
	})
}

func (cfg *apiConfig) apiRevokeToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "No authorization header found", nil)
	}

	if !strings.Contains(authHeader, "Bearer") {
		respondWithError(w, http.StatusUnauthorized, "No bearer field found in authorization header", nil)
	}

	splitString := strings.Fields(authHeader)

	if splitString[1] == "" {
		respondWithError(w, http.StatusUnauthorized, "No token found in bearer field", nil)
	}

	err := cfg.db.RevokeRefreshToken(r.Context(), splitString[1])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token", err)
		return
	}

	respondWithCodeOnly(w, http.StatusNoContent)
}

func (cfg *apiConfig) apidUpdateUser(w http.ResponseWriter, r *http.Request) {
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

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	updatedRecord, err := cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             userID,
		HashedPassword: hashedPassword,
		Email:          params.Email,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update record", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Id:        updatedRecord.ID,
		CreatedAt: updatedRecord.CreatedAt.String(),
		UpdatedAt: updatedRecord.UpdatedAt.String(),
		Email:     updatedRecord.Email,
	})

}
