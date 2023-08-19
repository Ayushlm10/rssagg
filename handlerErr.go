package main

import (
	"log"
	"net/http"
)

func handlerErr(w http.ResponseWriter, r *http.Request) {
	log.Println("Came here")
	respondWithError(w, 200, "Something went wrong")
}
