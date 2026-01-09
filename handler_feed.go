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

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct{
		Name string `json:"name"`
		URL string `json:"url"`
	}
	decoder:=json.NewDecoder(r.Body)

	params := parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON: %s",err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(),database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		Name: params.Name,
		Url: params.URL,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Couldn't create feed: %v", err))
		return
	}

	respondWithJSON(w,201,feed)
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request){

	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Couldn't get feed: %v", err))
		return
	}

	respondWithJSON(w,201,feeds)
}

