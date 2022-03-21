package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/headdetect/its-a-twitter/api/models"
	"github.com/headdetect/its-a-twitter/api/store"
	"github.com/headdetect/its-a-twitter/api/utils"
)

type loginRequest struct {
	Username string
	Password string
}

type loginResponse struct {
	AuthToken string `json:"authToken"`
	User *models.User `json:"user"`
}


func HandleUser(writer http.ResponseWriter, request *http.Request) {
	if !HasValidAuth(writer, request) {
		return
	}

	JsonResponse(writer, []byte(`{"message": "TODO"}`))
}

func HandleUserRegister(writer http.ResponseWriter, request *http.Request) {
	JsonResponse(writer, []byte(`{"message": "TODO"}`))
}

func HandleUserLogin(writer http.ResponseWriter, request *http.Request) {
	var loginRequest loginRequest

	err := json.NewDecoder(request.Body).Decode(&loginRequest)

	if err != nil || request.Method != "POST" {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := store.Users[loginRequest.Username]

	if !ok || !utils.CheckPasswordHash(loginRequest.Password, user.HashedPassword) {
		writer.WriteHeader(http.StatusUnauthorized)
		JsonResponse(writer, []byte(`{ message: "Invalid username or password" }`))
		return
	}
	
	// Persist through the session //
	authToken := utils.RandomString(32)	

	response := loginResponse {
		AuthToken: authToken,
		User: &user,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		JsonResponse(writer, []byte(`{ message: "We messed up somehow. Strange. We never mess up" }`))
		return
	}

	JsonResponse(writer, jsonResponse)
}