package handlers

import (
	"net/http"

	"github.com/headdetect/its-a-twitter/api/models"
	"github.com/headdetect/its-a-twitter/api/store"
)

func AuthUser(request *http.Request) (*models.User, bool) {
  authToken := request.Header.Get("AuthToken")
  authUsername := request.Header.Get("Username")

	if user, ok := store.Sessions[authToken]; ok {
		if authUsername == user.Username {
			return user, true
		}
	}

	return nil, false
}

func JsonResponse(writer http.ResponseWriter, response []byte) {
  writer.Header().Set("Content-Type", "application/json")
  writer.Header().Set("Access-Control-Allow-Origin", "*")

  writer.WriteHeader(http.StatusOK)
  writer.Write(response)
}

func RejectResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusUnauthorized)
  JsonResponse(writer, []byte(`{ message: "Invalid auth token provided. Please log in" }`))
}