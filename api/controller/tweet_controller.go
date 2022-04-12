package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/headdetect/its-a-twitter/api/model"
	"github.com/headdetect/its-a-twitter/api/utils"
)

type SingleTweetResponse struct {
	Tweet model.Tweet `json:"tweet"`
	ReactionCount map[string]int `json:"reactionCount"` // A reaction & count map //
	RetweetCount int `json:"retweetCount"`
}

type CreateTweetRequest struct {
	Text string `json:"text"`
}

type TweetReactionRequest struct {
	Reaction string `json:"reaction"`
}

func getTweet(request *http.Request) (model.Tweet, bool) {
	var tweet model.Tweet
	requestedTweetId, exists := GetPathValue(request, 0)

	if !exists {
		return tweet, false
	}

	tweetId, err := strconv.Atoi(requestedTweetId)

	if err != nil {
		return tweet, false 
	}

	tweet, err = model.GetTweetById(tweetId)

	if err != nil {
		return tweet, false
	}

	user, _, err := model.GetUserByTweetId(tweetId)

	tweet.User = &user

	return tweet, true
}

func HandleGetTweet(writer http.ResponseWriter, request *http.Request) {
	tweet, exists := getTweet(request)
	
	if !exists {
		NotFoundResponse(writer)
		return
	}

	retweetCount, err := tweet.RetweetCount()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	reactionCount, err := tweet.ReactionCount()

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	response := SingleTweetResponse {
		Tweet: tweet,
		RetweetCount: retweetCount,
		ReactionCount: reactionCount,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	JsonResponse(writer, jsonResponse)
}

func HandlePostTweet(writer http.ResponseWriter, request *http.Request) {
	var createTweetRequest CreateTweetRequest
	if err := json.NewDecoder(request.Body).Decode(&createTweetRequest); err != nil {
		BadRequestResponse(writer)
		return
	}

	currentUser, err := GetCurrentUser(request)

	if err != nil {
		// Returning an error response because this shouldn't be possible //
		ErrorResponse(writer, err)
		return
	}

	var mediaPath string
	
	// Receive file upload if there is one //
	file, _, err := request.FormFile("file")
	
	// The file exists. Proceed with uploading //
	if err == nil {
		defer file.Close()

		// [Scaleability]
		// This would be processed to reduce filesize as much as
		// possible using some lossless or a low lossy compression
		name := fmt.Sprintf("t-%s.jpg", utils.RandomString(64))

		// Copy to disk //
		mediaPath, _ := utils.GetStringOrDefault("MEDIA_PATH", "./assets/media")
		fullFilePath := fmt.Sprintf("%s/%s", mediaPath, name)
		diskFile, err := os.OpenFile(fullFilePath, os.O_CREATE | os.O_WRONLY, 0644)
		defer diskFile.Close()

		if err != nil {
			ErrorResponse(writer, err)
			return
		}

		io.Copy(diskFile, file)

		mediaPath = fullFilePath 
	}

	tweet, err := model.MakeTweet(currentUser, createTweetRequest.Text, mediaPath)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	response := SingleTweetResponse {
		Tweet: tweet,
	}

	if jsonResponse, err := json.Marshal(response); err == nil {
		JsonResponse(writer, jsonResponse)
	} else {
		ErrorResponse(writer, err)
	}
}

func HandleDeleteTweet(writer http.ResponseWriter, request *http.Request) {
	tweet, exists := getTweet(request)

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

  if tweet.User.Id != currentUser.Id {
		ForbiddenResponse(writer)
		return
	}

	if err = tweet.DeleteTweet(); err != nil {
		BadRequestResponse(writer)
	}
}

func HandleRetweet(writer http.ResponseWriter, request *http.Request) {
	tweet, exists := getTweet(request)

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

	if tweet.User.Id == currentUser.Id {
		// Not allowed to retweet self //
		BadRequestResponse(writer)
		return
	}

	if err = tweet.MakeRetweet(currentUser.Id); err != nil {
		BadRequestResponse(writer)
	}
}

func HandleRemoveRetweet(writer http.ResponseWriter, request *http.Request) {
	tweet, exists := getTweet(request)

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

	if err = tweet.DeleteRetweet(currentUser.Id); err != nil {
		BadRequestResponse(writer)
	}
}

func HandleReactTweet(writer http.ResponseWriter, request *http.Request) {
	tweet, exists := getTweet(request)

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

	var tweetReactionRequest TweetReactionRequest

	err = json.NewDecoder(request.Body).Decode(&tweetReactionRequest)

	if err != nil {
		BadRequestResponse(writer)
		return
	}

	if tweet.User.Id == currentUser.Id {
		// Not allowed to react to own tweet //
		BadRequestResponse(writer)
		return
	}

	if tweetReactionRequest.Reaction == "" {
		BadRequestResponse(writer)
	}

	// [scaling]
	// The valid emotes should be in some dynamic data store like a database
	// so a new app deploy wouldn't be required to add or remove.
	valid := false
	for _, reaction := range []string{ "👏", "🎉", "😔", "❤️", "👍", "👎" } {
		if reaction == tweetReactionRequest.Reaction {
			valid = true
			break
		}
	}

	if !valid {
		BadRequestResponse(writer)
	}

	if err = tweet.MakeReaction(currentUser.Id, tweetReactionRequest.Reaction); err != nil {
		BadRequestResponse(writer)
	}
}

func HandleRemoveReactTweet(writer http.ResponseWriter, request *http.Request) {
	tweet, exists := getTweet(request)

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
	
	if err = tweet.DeleteReaction(currentUser.Id); err != nil {
		BadRequestResponse(writer)
	}
}