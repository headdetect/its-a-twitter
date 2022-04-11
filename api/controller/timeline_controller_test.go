package controller_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func TestHandleTimeline(t *testing.T) {
	response, _, err := makeAuthenticatedRequest("test", http.MethodGet, "/timeline", nil)

	if err != nil {
		t.Errorf("Error authenticating. %k\n", err)
	}
	
	var actualResponse controller.TimelineResponse 
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

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
	
	response, _, err = makeAuthenticatedRequest(
		"admin",
		http.MethodPost,
		"/tweet", 
		bytes.NewReader([]byte(`{ "Text": "Tweet Tweet" }`)),
	)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", response.Status, response.StatusCode)
	}
	
	response, _, err = makeAuthenticatedRequest(
		"lurker", 
		http.MethodPost, 
		"/tweet", 
		bytes.NewReader([]byte(`{ "Text": "Tweet Tweet" }`)),
	)

	if err != nil {
		t.Fatalf("Got error while making auth request. %k\n", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", response.Status, response.StatusCode)
	}

	// Add a retweet from a followed user to a non-followed user  //
	response, _, err = makeAuthenticatedRequest(
		"basic", 
		http.MethodPut, 
		"/tweet/8/retweet", 
		nil,
	)

	if err != nil {
		t.Fatalf("Got error while making auth request. %k\n", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", response.Status, response.StatusCode)
	}

  response, _, err = makeAuthenticatedRequest("basic", http.MethodGet, "/timeline", nil)

	if err != nil {
		t.Errorf("Error authenticating. %k\n", err)
	}
	
	actualResponse = controller.TimelineResponse{}
	body, err = parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

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
