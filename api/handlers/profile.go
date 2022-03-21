package handlers

import (
	"encoding/json"
	"log"
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
	AuthToken string
	User *models.User
}

type userResponse struct {
	User *models.User
}


func HandleUser(writer http.ResponseWriter, request *http.Request) {
	if user, ok := AuthUser(request); ok {
		response := userResponse {
			User: user,
		}

		jsonResponse, err := json.Marshal(response)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			JsonResponse(writer, []byte(`{ message: "We messed up somehow. Strange. We never mess up" }`))
			return
		}

		JsonResponse(writer, jsonResponse)
		return
	}

	RejectResponse(writer)
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

	user, hashedPassword, err := store.GetUserWithPassByUsername(loginRequest.Username)

	if err != nil || !utils.CheckPasswordHash(loginRequest.Password, hashedPassword) {
		log.Printf("%k\n", err)
		writer.WriteHeader(http.StatusUnauthorized)
		JsonResponse(writer, []byte(`{ message: "Invalid username or password" }`))
		return
	}

	for authToken, val := range store.Sessions {
		if val.Id == user.Id {
			delete(store.Sessions, authToken)
			break
		}
	}
	
	// Persist through the session //
	authToken := utils.RandomString(32)	

	response := loginResponse {
		AuthToken: authToken,
		User: user,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		JsonResponse(writer, []byte(`{ message: "We messed up somehow. Strange. We never mess up" }`))
		return
	}

	JsonResponse(writer, jsonResponse)

	store.Sessions[authToken] = user
}