package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ayushlm10/rssAgg/internal/database"
	"github.com/google/uuid"
)

func (apicfg *apiCfg) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type Parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode the request body json: %v ", err))
		return
	}

	feed, err := apicfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create the user: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apicfg *apiCfg) handlerGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apicfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't fetch feeds: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedToFeeds(feeds))
}
