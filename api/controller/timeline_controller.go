package controller

import (
	"encoding/json"
	"net/http"

	"github.com/headdetect/its-a-twitter/api/model"
)

type TimelineResponse struct {
	Tweets []model.Tweet
}

func HandleTimeline(writer http.ResponseWriter, request *http.Request) {
	currentUser := GetCurrentUser(request)

	tweets, err := currentUser.GetTimeline(25)

	response := TimelineResponse {
		Tweets: tweets,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(err, writer)
		return
	}

	JsonResponse(writer, jsonResponse)
}