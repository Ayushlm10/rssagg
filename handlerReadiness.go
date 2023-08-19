package main

import (
	"log"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	log.Println("Came here")
	respondWithJSON(w, 200, struct{}{})
}
