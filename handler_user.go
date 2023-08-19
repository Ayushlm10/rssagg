package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ayushlm10/rssAgg/internal/database"
	"github.com/google/uuid"
)

func (apicfg *apiCfg) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't decode the request body json: %v ", err))
		return
	}

	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create the user: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseUserToUser(user))
}

func (apicfg *apiCfg) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}

func (apicfg *apiCfg) handlerGetPostForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apicfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts for user: %v", err))
		return
	}
	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
