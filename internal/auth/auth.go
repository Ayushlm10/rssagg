package auth

import (
	"errors"
	"net/http"
	"strings"
)

//Extracts the api key from the request headers if present
// should be of the form "Authorization": "ApiKey <api_key>"

func GetApiKey(headers http.Header) (apikey string, err error) {
	auth_string := headers.Get("Authorization")
	if auth_string == "" {
		return "", errors.New("auth header does not exist")
	}
	authVals := strings.Split(auth_string, " ")
	if len(authVals) != 2 {
		return "", errors.New("incorrect auth header format")
	}

	if authVals[0] != "ApiKey" {
		return "", errors.New("auth header is in the wrong format")
	}

	if authVals[1] == "" {
		return "", errors.New("api key is required")
	}

	return authVals[1], nil
}
