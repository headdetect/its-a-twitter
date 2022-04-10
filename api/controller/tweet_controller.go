package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/headdetect/its-a-twitter/api/model"
	"github.com/headdetect/its-a-twitter/api/utils"
)

type SingleTweetResponse struct {
	Tweet model.Tweet
}

type CreateTweetRequest struct {
	Text string
}

func getTweet(request *http.Request) (model.Tweet, error) {
	return model.Tweet{}, nil
}

func HandleGetTweet(writer http.ResponseWriter, request *http.Request) {
	
}

func HandlePostTweet(writer http.ResponseWriter, request *http.Request) {
	var createTweetRequest CreateTweetRequest
	if err := json.NewDecoder(request.Body).Decode(&createTweetRequest); err != nil {
		BadRequestResponse(writer)
		return
	}

	currentUser := GetCurrentUser(request)

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
			ErrorResponse(err, writer)
			return
		}

		io.Copy(diskFile, file)

		mediaPath = fullFilePath 
	}

	tweet, err := model.MakeTweet(currentUser.Id, createTweetRequest.Text, mediaPath)

	response := SingleTweetResponse {
		Tweet: tweet,
	}

	if jsonResponse, err := json.Marshal(response); err == nil {
		JsonResponse(writer, jsonResponse)
	} else {
		ErrorResponse(err, writer)
	}
}

func HandleDeleteTweet(writer http.ResponseWriter, request *http.Request) {
	
}

func HandleRetweet(writer http.ResponseWriter, request *http.Request) {
	
}

func HandleRemoveRetweet(writer http.ResponseWriter, request *http.Request) {
	
}

func HandleReactTweet(writer http.ResponseWriter, request *http.Request) {
	
}

func HandleRemoveReactTweet(writer http.ResponseWriter, request *http.Request) {
	
}