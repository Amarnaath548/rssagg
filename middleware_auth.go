package main

import (
	"fmt"
	"net/http"

	"github.com/Amarnaath548/rssagg/internal/auth"
	"github.com/Amarnaath548/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)


func (apiCfg *apiConfig) middlewareAuth(hanler authedHandler) http.HandlerFunc {
	return func (w http.ResponseWriter,r *http.Request) {
	apiKey, err:= auth.GetAPIKey(r.Header)
	if err!=nil{
		respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}
	user, err :=apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey )
	if err!=nil{
		respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
		return
	}
	hanler(w, r, user)
}
}