package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) apiPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if !strings.EqualFold(params.Event, "user.upgraded") {
		respondWithCodeOnly(w, http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse UUID", err)
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithCodeOnly(w, http.StatusNotFound)
	}
	respondWithCodeOnly(w, http.StatusNoContent)
}
