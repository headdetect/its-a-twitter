package controller_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

// As test. UserId = 2

// Follows. UserId = 1, 3
// Tweets from followed. TweetId = 1, 2, 4, 5, 6
// Retweets from followed. TweetId = 1 (dup), 3, 7

// Timeline should be:
// tweetId, retweetUserId
// 7, 3
// 6, null
// 5, null
// 4, null
// 3, null
// 2, null
// 1, null
func TestHandleTimeline(t *testing.T) {
	var actualResponse controller.TimelineResponse
	response, _ := makeAuthenticatedTestRequest(t, "test", http.MethodGet, "/timeline", nil)
	parseTestResponse(t, response, &actualResponse)

	tweets := actualResponse.Tweets

	// Verify state from seeded value //
	validTweetIds := []int{7, 6, 5, 4, 3, 2, 1}

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
	makeAuthenticatedTestFormRequest(
		t,
		"admin",
		http.MethodPost,
		"/tweet",
		strings.NewReader(fmt.Sprintf("text=%s", url.QueryEscape("Tweet Tweet"))),
	)

	makeAuthenticatedTestFormRequest(
		t,
		"lurker",
		http.MethodPost,
		"/tweet",
		strings.NewReader(fmt.Sprintf("text=%s", url.QueryEscape("Tweet Tweet"))),
	)

	// Add a retweet from a followed user to a non-followed user  //
	makeAuthenticatedTestRequest(
		t,
		"basic",
		http.MethodPut,
		"/tweet/8/retweet",
		nil,
	)

	// Make our own tweet //
	makeAuthenticatedTestFormRequest(
		t,
		"test",
		http.MethodPost,
		"/tweet",
		strings.NewReader(fmt.Sprintf("text=%s", url.QueryEscape("Tweet Tweet"))),
	)

	actualResponse = controller.TimelineResponse{}
	response, _ = makeAuthenticatedTestRequest(t, "test", http.MethodGet, "/timeline", nil)
	parseTestResponse(t, response, &actualResponse)

	tweets = actualResponse.Tweets

	// Verify the timeline only shows the followed user's tweet/retweet //
	validTweetIds = append([]int{10, 8}, validTweetIds...) // Add the new retweet //

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

	// Verify our tweet is in the #1 spot //
	if tweets[0].Poster != 2 {
		t.Errorf("Expected 'test' got '%s'\n", tweets[0].Tweet.User.Username)
	}
}


func TestNewHandleTimeline(t *testing.T) {
	response, _ := makeRequest(
		http.MethodPost,
		"/user/register",
		bytes.NewReader([]byte(`{ "username": "timelineuser", "password": "password", "email": "not-real@example.com" }`)),
	)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", response.Status, response.StatusCode)
		return
	}

	var profileResponse controller.ProfileUserResponse
	parseResponse(response, &profileResponse)

	if profileResponse.User.Username != "timelineuser" {
		t.Errorf(
			"expected 'timelineuser' got '%s'",
			profileResponse.User.Username,
		)

		return
	}

	userId := profileResponse.User.Id

	var timelineResponse controller.TimelineResponse
	response, _ = makeAuthenticatedTestRequest(t, "timelineuser", http.MethodGet, "/timeline", nil)
	parseTestResponse(t, response, &timelineResponse)

	tweets := timelineResponse.Tweets

	if len(tweets) != 0 {
		t.Errorf(
			"Expected 0 tweets on timeline. Got %d\n",
			len(timelineResponse.Tweets),
		)

		return
	}

	// Add a tweet from a followed user and a non-followed user  //
	makeAuthenticatedTestFormRequest(
		t,
		"timelineuser",
		http.MethodPost,
		"/tweet",
		strings.NewReader(fmt.Sprintf("text=%s", url.QueryEscape("Tweet Tweet"))),
	)

	timelineResponse = controller.TimelineResponse{}
	response, _ = makeAuthenticatedTestRequest(t, "timelineuser", http.MethodGet, "/timeline", nil)
	parseTestResponse(t, response, &timelineResponse)

	tweets = timelineResponse.Tweets

	if len(tweets) != 1 {
		t.Errorf(
			"Expected 1 tweet on timeline. Got %d\n",
			len(tweets),
		)

		return
	}

	// Verify our tweet is in the #1 spot //
	if tweets[0].Poster != userId {
		t.Errorf("Expected '%d' got '%d'\n", userId, tweets[0].Poster)
	}

	if tweets[0].Tweet.Text != "Tweet Tweet" {
		t.Errorf("Expected 'Tweet Tweet' got '%s'\n", tweets[0].Tweet.Text)
	}
}