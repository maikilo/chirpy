package main

import (
	"net/http"
	"time"

	"github.com/maikilo/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokenRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find token", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token is not valid", err)
		return
	}

	expirationTime := time.Hour

	accessToken, err := auth.MakeJWT(
		refreshToken.UserID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) handlerTokenRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find token in request header", err)
		return
	}

	_, err = cfg.db.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find refresh token", err)
		return
	}
	
	_, err = cfg.db.UpdateRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
