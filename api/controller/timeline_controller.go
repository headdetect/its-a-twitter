package controller

import (
	"encoding/json"
	"net/http"

	"github.com/headdetect/its-a-twitter/api/model"
)

type TimelineResponse struct {
	Tweets []model.TimelineTweet `json:"tweets"`

	// A map of users (id : user) that are mentioned in the timeline //
	Users map[int]model.User `json:"users"`
}

func HandleTimeline(writer http.ResponseWriter, request *http.Request) {
	currentUser, err := GetCurrentUser(request)

	if err != nil {
		// Returning an error response because this shouldn't be possible //
		ErrorResponse(writer, err)
		return
	}

	tweets, users, err := model.GetTimeline(currentUser.Id, 25)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	response := TimelineResponse{
		Tweets: tweets,
		Users:  users,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		ErrorResponse(writer, err)
		return
	}

	JsonResponse(writer, jsonResponse)
}
