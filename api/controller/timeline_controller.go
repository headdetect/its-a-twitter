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

	var tweets []model.TimelineTweet
	var users map[int]model.User

	if err != nil {
		// The user is unauthenticated //
		tweets, users, err = model.GetFeatured()
	} else {
		tweets, users, err = model.GetTimeline(currentUser.Id)

		if err == nil && len(tweets) == 0 {
			// Show the featured if we don't have anything on our timeline //
			tweets, users, err = model.GetFeatured()
		}
	}

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
