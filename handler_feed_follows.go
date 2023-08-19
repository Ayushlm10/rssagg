package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ayushlm10/rssAgg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apicfg *apiCfg) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type Parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode the request body json: %v ", err))
		return
	}

	feedfollow, err := apicfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create the user: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedfollow))
}

func (apicfg *apiCfg) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedfollows, err := apicfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create the user: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowsToFeedFollows(feedfollows))
}

func (apicfg *apiCfg) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse the feed id: %v", err))
		return
	}
	err = apicfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't delete the user: %v", err))
		return
	}
	respondWithJSON(w, 204, struct{}{})
}
