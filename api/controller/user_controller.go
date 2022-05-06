package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

type profileTimeline struct {
	Tweets []model.TimelineTweet `json:"tweets"`
	Users  map[int]model.User    `json:"users"`
}

type ProfileUserResponse struct {
	User      model.User      `json:"user"`
	Followers []model.User    `json:"followers"`
	Following []model.User    `json:"following"`
	Timeline  profileTimeline `json:"timeline"`
}

type SmallProfileUserResponse struct {
	User           model.User      `json:"user"`
	FollowerCount  int             `json:"followerCount"`
	FollowingCount int             `json:"followingCount"`
	Timeline       profileTimeline `json:"timeline"`
}

type ProfileImageChangedResponse struct {
	ProfilePicPath string `json:"profilePicPath"`
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

const (
	MAX_USERNAME_CHARS = 20
)

func getUser(request *http.Request) (model.User, bool) {
	requestedUserName, exists := GetPathValue(request, 0)

	if !exists {
		return model.User{}, false
	}

	user, _, err := model.GetUserByUsername(requestedUserName)

	return user, err == nil
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

	tweets, usersFromTweets, err := model.GetTimelineByUser(currentUser, nil)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	response := ProfileUserResponse{
		User:      currentUser,
		Followers: followers,
		Following: following,
		Timeline: profileTimeline{
			Tweets: tweets,
			Users:  usersFromTweets,
		},
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	JsonResponse(writer, jsonResponse)
}

func HandleUpdateUserAvatar(writer http.ResponseWriter, request *http.Request) {
	// Receive file upload if there is one //
	file, fileHeader, err := request.FormFile("file")

	if err != nil {
		BadRequestResponse(writer)
	}

	defer file.Close()

	fileMime := fileHeader.Header.Get("Content-Type")
	types := strings.Split(fileMime, "/")

	if strings.TrimSpace(fileMime) == "" || len(types) != 2 {
		BadRequestResponse(writer)
		return
	}

	valid := false

	for _, mime := range ACCEPTABLE_MIME_TYPES {
		if mime == fileMime {
			valid = true
			break
		}
	}

	if !valid {
		BadRequestResponse(writer)
		return
	}

	extension := types[1]

	// [Scaleability]
	// This would be processed to reduce filesize as much as
	// possible using some lossless or a low lossy compression
	//
	// We'd also want to check for filename conflicts.
	name := fmt.Sprintf("u-%s.%s", utils.RandomHex(8), extension)

	// Copy to disk //
	path, _ := utils.GetStringOrDefault("MEDIA_PATH", "./assets/media")
	fullFilePath := fmt.Sprintf("%s/%s", path, name)
	diskFile, err := os.OpenFile(fullFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	defer diskFile.Close()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	io.Copy(diskFile, file)
	
	currentUser, err := GetCurrentUser(request)

	if err != nil {
		// Returning an error response because this shouldn't be possible //
		ErrorResponse(writer, err)
		return
	}

	err = currentUser.UpdateProfilePicPath(name)

	if err != nil {
		// Returning an error response because this shouldn't be possible //
		ErrorResponse(writer, err)
		return
	}

	// Update logged in user //
	for auth, val := range Sessions {
		if val.Id == currentUser.Id {
			val.ProfilePicPath = name // Update //

			Sessions[auth] = val // Copy back //
		}
	}

	response := ProfileImageChangedResponse{
		ProfilePicPath: name,
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

	currentUser, err := GetCurrentUser(request)

	var response any

	if err == nil {
		tweets, usersFromTweets, err := model.GetTimelineByUser(requestedUser, &currentUser)

		if err != nil {
			ErrorResponse(writer, err)
			return
		}

		response = ProfileUserResponse{
			User:      requestedUser,
			Followers: followers,
			Following: following,
			Timeline: profileTimeline{
				Tweets: tweets,
				Users:  usersFromTweets,
			},
		}
	} else {
		tweets, usersFromTweets, err := model.GetTimelineByUser(requestedUser, nil)

		if err != nil {
			ErrorResponse(writer, err)
			return
		}

		response = SmallProfileUserResponse{
			User:           requestedUser,
			FollowerCount:  len(followers),
			FollowingCount: len(following),
			Timeline: profileTimeline{
				Tweets: tweets,
				Users:  usersFromTweets,
			},
		}
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

	registerUserRequest.Username = strings.ToLower(registerUserRequest.Username)

	match, err := regexp.MatchString("(?i)^[a-z0-9_-]*$", registerUserRequest.Username)

	if !match || err != nil || len(registerUserRequest.Username) > MAX_USERNAME_CHARS {
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

	user, err := model.MakeUser(registerUserRequest.Email, registerUserRequest.Username, hashedPassword)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {

			// This status is traditionally reserved for PUT requests
			// but in this case it makes sense too
			ConflictRequestResponse(writer)
			return
		}

		BadRequestResponse(writer)
	}

	// Sign them in now //
	authToken := utils.RandomString(32)

	// Email should be explicitly assigned so it's not accidentally leaked //
	user.Email = registerUserRequest.Email

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

	Sessions[authToken] = user
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

	for authToken, val := range Sessions {
		if val.Id == user.Id {
			delete(Sessions, authToken)
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

	Sessions[authToken] = user
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
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			// We're attempting to follow a user that we already follow. Ignore //
			return
		}

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

	if err = currentUser.UnFollowUser(requestedUser.Id); err != nil {
		ErrorResponse(writer, err)
	}
}
