package handlers

import (
	"net/http"

	"github.com/headdetect/its-a-twitter/api/store"
)

func HasValidAuth(writer http.ResponseWriter, request *http.Request) bool {
  authToken := request.Header.Get("AuthToken")
  authUsername := request.Header.Get("Username")

	if val, ok := store.Sessions[authToken]; ok {
		if authUsername == val.Username {
			return true
		}
	}

	writer.WriteHeader(http.StatusUnauthorized)
  JsonResponse(writer, []byte(`{ message: "Invalid auth token provided. Please log in" }`))
	return false
}

func JsonResponse(writer http.ResponseWriter, response []byte) {
  writer.Header().Set("Content-Type", "application/json")
  writer.Header().Set("Access-Control-Allow-Origin", "*")

  writer.WriteHeader(http.StatusOK)
  writer.Write(response)
}