package controller

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/headdetect/its-a-twitter/api/model"
	"github.com/headdetect/its-a-twitter/api/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthToken string     `json:"authToken"`
	User      model.User `json:"user"`
}

type OwnUserResponse struct {
	User      model.User    `json:"user"`
	Followers []model.User  `json:"followers"`
	Following []model.User  `json:"following"`
	Tweets    []model.Tweet `json:"tweets"`
}

type UserResponse struct {
	User           model.User    `json:"user"`
	FollowerCount  int           `json:"followerCount"`
	FollowingCount int           `json:"followingCount"`
	Tweets         []model.Tweet `json:"tweets"`
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`

	// [Scaling]
	// In a production-grade app, this should be hashed from the
	// client using the email or username as the salt.
	// From the API side, it would then hash it again using a private salt
	// when storing
	Password string `json:"password"`
}

func getUser(request *http.Request) (model.User, bool) {
	requestedUserName, exists := GetPathValue(request, 0)

	if !exists {
		return model.User{}, false
	}

	user, _, err := model.GetUserByUsername(requestedUserName)

	return user, err == nil
}

func HandleUserFollowUser(writer http.ResponseWriter, request *http.Request) {
	JsonResponse(writer, []byte(`{"message": "TODO"}`))
}

func HandleOwnUser(writer http.ResponseWriter, request *http.Request) {
	currentUser, err := GetCurrentUser(request)

	if err != nil {
		// Returning an error response because this shouldn't be possible //
		ErrorResponse(writer, err)
		return
	}

	following, err := currentUser.GetFollowing()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	followers, err := currentUser.GetFollowers()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	tweets, err := currentUser.GetTweets()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	response := OwnUserResponse{
		User:      currentUser,
		Followers: followers,
		Following: following,
		Tweets:    tweets,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	JsonResponse(writer, jsonResponse)
}

func HandleUser(writer http.ResponseWriter, request *http.Request) {
	requestedUser, exists := getUser(request)

	if !exists {
		NotFoundResponse(writer)
		return
	}

	following, err := requestedUser.GetFollowing()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	followers, err := requestedUser.GetFollowers()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	tweets, err := requestedUser.GetTweets()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	response := UserResponse{
		User:           requestedUser,
		FollowerCount:  len(followers),
		FollowingCount: len(following),
		Tweets:         tweets,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	JsonResponse(writer, jsonResponse)
}

func HandleUserRegister(writer http.ResponseWriter, request *http.Request) {
	var registerUserRequest RegisterUserRequest

	err := json.NewDecoder(request.Body).Decode(&registerUserRequest)

	if err != nil {
		BadRequestResponse(writer)
		return
	}

	if match, err := regexp.MatchString("^[a-z0-9_-]*$", registerUserRequest.Username); !match || err != nil {
		// This should be handled from client side. 
		// No need to get specific on the error

		BadRequestResponse(writer)
		return
	}

	hashedPassword, err := utils.HashPassword(registerUserRequest.Password)

	if err != nil {
		BadRequestResponse(writer)
		return
	}

	if _, err := model.MakeUser(registerUserRequest.Email, registerUserRequest.Username, hashedPassword); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {

			// These are traditionally reserved for PUT requests
			// but in this case it makes sense too
			ConflictRequestResponse(writer)
			return
		}
		
		BadRequestResponse(writer)
	}
}

func HandleUserLogin(writer http.ResponseWriter, request *http.Request) {
	var loginRequest LoginRequest

	err := json.NewDecoder(request.Body).Decode(&loginRequest)

	if err != nil {
		BadRequestResponse(writer)
		return
	}

	user, hashedPassword, email, err := model.GetUserByUsernameWithPass(loginRequest.Username)

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

	// Email should be explicitly assigned so it's not accidentally leaked //
	user.Email = email

	response := LoginResponse{
		AuthToken: authToken,
		User:      user,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	JsonResponse(writer, jsonResponse)

	sessions[authToken] = user
}

func HandleFollowUser(writer http.ResponseWriter, request *http.Request) {
	requestedUser, exists := getUser(request)

	if !exists {
		NotFoundResponse(writer)
		return
	}

	currentUser, err := GetCurrentUser(request)

	if err != nil {
		// Returning an error response because this shouldn't be possible w/getUser //
		ErrorResponse(writer, err)
		return
	}

	if requestedUser.Id == currentUser.Id {
		BadRequestResponse(writer)
		return
	}

	if err = currentUser.FollowUser(requestedUser.Id); err != nil {
		ErrorResponse(writer, err)
	}
}

func HandleUnFollowUser(writer http.ResponseWriter, request *http.Request) {
	requestedUser, exists := getUser(request)

	if !exists {
		NotFoundResponse(writer)
		return
	}

	currentUser, err := GetCurrentUser(request)

	if err != nil {
		// Returning an error response because this shouldn't be possible //
		ErrorResponse(writer, err)
		return
	}

	// TODO: Validate they can remove retweet

	if err = currentUser.UnFollowUser(requestedUser.Id); err != nil {
		ErrorResponse(writer, err)
	}
}
