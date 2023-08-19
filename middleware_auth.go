package main

import (
	"fmt"
	"net/http"

	"github.com/Ayushlm10/rssAgg/internal/auth"
	"github.com/Ayushlm10/rssAgg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apicfg *apiCfg) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetApiKey(r.Header)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get API key: %v", err))
			return
		}

		if apikey == "" {
			respondWithError(w, 400, fmt.Sprintf("Wrong auth info: %v", err))
			return
		}

		user, err := apicfg.DB.GetUserByApikey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
