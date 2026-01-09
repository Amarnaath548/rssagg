package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Amarnaath548/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User){
	type parameters struct{
		FeedTD uuid.UUID `json:"feed_id"`
	}
	decoder:=json.NewDecoder(r.Body)

	params := parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON: %s",err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		UserID: user.ID,
		FeedID: params.FeedTD,
	})
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Couldn't create feed follow: %v", err))
		return
	}

	respondWithJSON(w,201,feedFollow)
}


func (apiCfg *apiConfig) handlerGetFeedFollow(w http.ResponseWriter,r *http.Request, user database.User){
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(),user.ID)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Could'n get feed follows: %v",err))
		return
	}
	respondWithJSON(w,200,feedFollows)
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter,r *http.Request, user database.User){
	feedFollowIDStr :=chi.URLParam(r, "feedFollowID")
	feedFollowID, err :=uuid.Parse(feedFollowIDStr)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Could'n parse feed follow id: %v",err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		ID: feedFollowID,
		UserID: user.ID,
	})
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Could'n delete feed follow: %v",err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}