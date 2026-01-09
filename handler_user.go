package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Amarnaath548/rssagg/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Name string `json:"name"`
	}
	decoder:=json.NewDecoder(r.Body)

	params := parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON: %s",err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	respondWithJSON(w,201,user)
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter,r *http.Request, user database.User ) {
	respondWithJSON(w, 200, user)
}