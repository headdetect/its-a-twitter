package controller

import (
	"encoding/json"
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
	User model.User
}

type OwnUserResponse struct {
	User model.User
}

type UserResponse struct {
	User model.User
	Followers []model.User
	Following []model.User
	Tweets []model.Tweet
}

func HandleUserFollowUser(writer http.ResponseWriter, request *http.Request) {
	JsonResponse(writer, []byte(`{"message": "TODO"}`))
}

func HandleOwnUser(writer http.ResponseWriter, request *http.Request) {
	currentUser := GetCurrentUser(request)

	response := OwnUserResponse {
		User: *currentUser,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(err, writer)
		return
	}

	JsonResponse(writer, jsonResponse)
}

func HandleUser(writer http.ResponseWriter, request *http.Request) {
	username, exists := GetPathValue(request, 0)

	if !exists {
		BadRequestResponse(writer)
		return
	}

	user, err := model.GetUserByUsername(username)

	if err != nil {
		NotFoundResponse(writer)
		return
	}

	following, err := user.GetFollowing()

	if err != nil {
		ErrorResponse(err, writer)
		return
	}

	followers, err := user.GetFollowers()

	if err != nil {
		ErrorResponse(err, writer)
		return
	}

	tweets, err := user.GetTweets()

	if err != nil {
		ErrorResponse(err, writer)
		return
	}

	response := UserResponse {
		User: user,
		Followers: followers,
		Following: following,
		Tweets: tweets,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(err, writer)
		return
	}

	JsonResponse(writer, jsonResponse)
}

func HandleUserRegister(writer http.ResponseWriter, request *http.Request) {
	JsonResponse(writer, []byte(`{"message": "TODO"}`))
}

func HandleUserLogin(writer http.ResponseWriter, request *http.Request) {
	var loginRequest LoginRequest

	err := json.NewDecoder(request.Body).Decode(&loginRequest)

	if err != nil {
		BadRequestResponse(writer)
		return
	}

	user, hashedPassword, err := model.GetUserWithPassByUsername(loginRequest.Username)

	if err != nil || !utils.CheckPasswordHash(loginRequest.Password, hashedPassword) {
		UnauthorizedResponse(writer)
		return
	}

	for authToken, val := range sessions {
		if val.Id == user.Id {
			delete(sessions, authToken)
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
		ErrorResponse(err, writer)
		return
	}

	JsonResponse(writer, jsonResponse)

	sessions[authToken] = user
}
