package controller_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func TestHandleTimeline(t *testing.T) {
	var actualResponse controller.TimelineResponse 
	response, _ := makeAuthenticatedTestRequest(t, "test", http.MethodGet, "/timeline", nil)
	parseTestResponse(t, response, &actualResponse)

	tweets := actualResponse.Tweets

	// Verify state from seeded value //
	validTweetIds := []int{7, 6, 5, 4, 2, 1}

	if len(tweets) != len(validTweetIds) {
		t.Errorf(
			"Expected %d tweets on timeline. Got %d\n", 
			len(validTweetIds), 
			len(actualResponse.Tweets),
		)
	}

	for i, tweetId := range validTweetIds {
		if tweets[i].Id != tweetId {
			t.Errorf("Expected Tweet ID to be %d. Got %d\n", tweetId, tweets[i].Id)
		}
	}

	// Add a tweet from a followed user and a non-followed user  //
	makeAuthenticatedTestRequest(
		t,
		"admin",
		http.MethodPost,
		"/tweet", 
		bytes.NewReader([]byte(`{ "text": "Tweet Tweet" }`)),
	)
	
	makeAuthenticatedTestRequest(
		t,
		"lurker", 
		http.MethodPost, 
		"/tweet", 
		bytes.NewReader([]byte(`{ "text": "Tweet Tweet" }`)),
	)

	// Add a retweet from a followed user to a non-followed user  //
	makeAuthenticatedTestRequest(
		t,
		"basic", 
		http.MethodPut,
		"/tweet/8/retweet", 
		nil,
	)

	actualResponse = controller.TimelineResponse{}
  response, _ = makeAuthenticatedTestRequest(t, "test", http.MethodGet, "/timeline", nil)
	parseTestResponse(t, response, &actualResponse)

	tweets = actualResponse.Tweets

	// Verify the timeline only shows the followed user's tweet/retweet //
	validTweetIds = append([]int{ 8 }, validTweetIds...) // Add the new retweet //

	if len(tweets) != len(validTweetIds) {
		t.Errorf(
			"Expected %d tweets on timeline. Got %d\n", 
			len(validTweetIds), 
			len(actualResponse.Tweets),
		)
	}

	for i, tweetId := range validTweetIds {
		if tweets[i].Id != tweetId {
			t.Errorf("Expected Tweet ID to be %d. Got %d.\n", tweetId, tweets[i].Id)
		}
	}
}
