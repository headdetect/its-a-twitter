package controller_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

// As test. UserId = 2

// Follows. UserId = 1, 3
// Tweets from followed. TweetId = 1, 2, 4, 5, 6
// Retweets from followed. TweetId = 1 (dup), 3 (own), 7

// Timeline should be:
// tweetId, retweetUserId
// 7, 3
// 6, null
// 5, null
// 4, null
// 2, null
// 1, null
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
		if tweets[i].Tweet.Id != tweetId {
			t.Errorf("Expected Tweet ID to be %d. Got %d\n", tweetId, tweets[i].Tweet.Id)
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
	validTweetIds = append([]int{8}, validTweetIds...) // Add the new retweet //

	if len(tweets) != len(validTweetIds) {
		t.Errorf(
			"Expected %d tweets on timeline. Got %d\n",
			len(validTweetIds),
			len(actualResponse.Tweets),
		)
	}

	for i, tweetId := range validTweetIds {
		if tweets[i].Tweet.Id != tweetId {
			t.Errorf("Expected Tweet ID to be %d. Got %d.\n", tweetId, tweets[i].Tweet.Id)
		}
	}
}
