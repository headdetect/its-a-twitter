package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/headdetect/its-a-twitter/api/model"
	"github.com/headdetect/its-a-twitter/api/utils"
)

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	AuthToken string
	User *model.User
}

type UserResponse struct {
	User *model.User
}

func HandleUserFollowUser(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	JsonResponse(writer, []byte(`{"message": "TODO"}`))
}

func HandleUser(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	if user, ok := AuthUser(request); ok {
		response := UserResponse {
			User: user,
		}

		jsonResponse, err := json.Marshal(response)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			JsonResponse(writer, []byte(`{ "message": "We messed up somehow. Strange. We never mess up" }`))
			return
		}

		JsonResponse(writer, jsonResponse)
		return
	}

	RejectResponse(writer)
}

func HandleUserRegister(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	JsonResponse(writer, []byte(`{"message": "TODO"}`))
}



func HandleUserLogin(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var loginRequest LoginRequest

	err := json.NewDecoder(request.Body).Decode(&loginRequest)

	if err != nil {
		log.Printf("%k\n", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Method != "POST" {
		http.Error(writer, "Not found", http.StatusNotFound)
		return
	}

	user, hashedPassword, err := model.GetUserWithPassByUsername(loginRequest.Username)

	if err != nil || !utils.CheckPasswordHash(loginRequest.Password, hashedPassword) {
		log.Printf("%k\n", err)
		writer.WriteHeader(http.StatusUnauthorized)
		JsonResponse(writer, []byte(`{ message: "Invalid username or password" }`))
		return
	}

	for authToken, val := range model.Sessions {
		if val.Id == user.Id {
			delete(model.Sessions, authToken)
			break
		}
	}
	
	// Persist through the session //
	authToken := utils.RandomString(32)	

	response := LoginResponse {
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

	model.Sessions[authToken] = user
}