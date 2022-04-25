package controller_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func TestHandleGetTweet(t *testing.T) {
	response, _ := makeTestRequest(t, http.MethodGet, "/tweet/1", nil)

	var actualResponse controller.SingleTweetResponse
	parseTestResponse(t, response, &actualResponse)

	if actualResponse.TimelineTweet.Tweet.Id != 1 {
		t.Errorf(
			"expected id:1 got id:%d\n",
			actualResponse.TimelineTweet.Tweet.Id,
		)
	}

	if actualResponse.TimelineTweet.Tweet.Text != "First tweet ever" {
		t.Errorf(
			"expected 'First tweet ever' got '%s'\n",
			actualResponse.TimelineTweet.Tweet.Text,
		)
	}

	count, ok := actualResponse.TimelineTweet.ReactionCount["party"]

	if !ok {
		t.Errorf("expected 'party' to exist\n")
	}

	if count != 2 {
		t.Errorf("expected 2 got %d\n", count)
	}

	if actualResponse.TimelineTweet.RetweetCount != 2 {
		t.Errorf("expected 2 got %d\n", actualResponse.TimelineTweet.RetweetCount)
	}
}

func TestHandlePostTweet(t *testing.T) {
	loginRequest := controller.LoginRequest{
		Username: "lurker",
		Password: "password",
	}

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/tweet",
		strings.NewReader(fmt.Sprintf("text=%s", url.QueryEscape(`Tweet Tweet`))),
	)

	err := AuthenticatedRequest(loginRequest, request)

	if err != nil {
		t.Errorf("Error authenticating. %k\n", err)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")

	ControllerServe(writer, request)

	response := writer.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", response.Status, response.StatusCode)
	}

	var parsedResponse controller.SingleTweetResponse
	parseTestResponse(t, response, &parsedResponse)

	if parsedResponse.TimelineTweet.Tweet.Text != "Tweet Tweet" {
		t.Errorf(
			"expected 'Tweet Tweet' got %s.",
			parsedResponse.TimelineTweet.Tweet.Text,
		)
	}
}

func TestHandleRetweet(t *testing.T) {
	var actualResponse controller.SingleTweetResponse

	response, _ := makeRequest(http.MethodGet, "/tweet/1", nil)
	parseTestResponse(t, response, &actualResponse)

	count := actualResponse.TimelineTweet.RetweetCount

	// Verify can put //
	makeAuthenticatedTestRequest(t, "lily", http.MethodPut, "/tweet/1/retweet", nil)

	// Verify updated //
	response, _ = makeTestRequest(t, http.MethodGet, "/tweet/1", nil)
	parseTestResponse(t, response, &actualResponse)

	if actualResponse.TimelineTweet.RetweetCount != count+1 {
		t.Errorf("Expecting %d got %d\n", count+1, actualResponse.TimelineTweet.RetweetCount)
	}

	// Verify can delete //
	makeAuthenticatedTestRequest(t, "lily", http.MethodDelete, "/tweet/1/retweet", nil)

	// Verify updated //
	response, _ = makeTestRequest(t, http.MethodGet, "/tweet/1", nil)
	parseTestResponse(t, response, &actualResponse)

	if actualResponse.TimelineTweet.RetweetCount != count {
		t.Errorf("Expecting %d got %d\n", count, actualResponse.TimelineTweet.RetweetCount)
	}
}

func TestHandleReactTweet(t *testing.T) {
	var initialResponse controller.SingleTweetResponse
	response, _ := makeTestRequest(t, http.MethodGet, "/tweet/7", nil)
	parseTestResponse(t, response, &initialResponse)

	count := initialResponse.TimelineTweet.ReactionCount["party"]

	// Verify can put //
	makeAuthenticatedTestRequest(
		t,
		"lily",
		http.MethodPut,
		"/tweet/7/react/party",
		nil,
	)

	// Verify updated //
	var updatedResponse controller.SingleTweetResponse
	response, _ = makeTestRequest(t, http.MethodGet, "/tweet/7", nil)
	parseTestResponse(t, response, &updatedResponse)

	if updatedResponse.TimelineTweet.ReactionCount["party"] != count+1 {
		t.Errorf("Expecting %d got %d\n", count+1, updatedResponse.TimelineTweet.ReactionCount["party"])
	}

	// Verify can delete //
	makeAuthenticatedTestRequest(
		t,
		"lily",
		http.MethodDelete,
		"/tweet/7/react/party",
		nil,
	)

	// Verify updated //
	var deletedResponse controller.SingleTweetResponse
	response, _ = makeTestRequest(t, http.MethodGet, "/tweet/7", nil)
	parseTestResponse(t, response, &deletedResponse)

	if deletedResponse.TimelineTweet.ReactionCount["party"] != count {
		t.Errorf("Expecting %d got %d\n", count, deletedResponse.TimelineTweet.ReactionCount["party"])
	}
}

// Non happy path //
func TestHandleCannotRetweetOwnTweet(t *testing.T) {
	response, _, err := makeAuthenticatedRequest("admin", http.MethodPut, "/tweet/1/retweet", nil)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("expected Bad Request (400) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotReactOwnTweet(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"admin",
		http.MethodPut,
		"/tweet/1/react/🎉",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("expected Bad Request (400) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotReactInvalid(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"lily",
		http.MethodPut,
		"/tweet/1/react/not-a-real-reaction",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("expected Bad Request (400) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotDeleteOtherTweet(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"lily",
		http.MethodDelete,
		"/tweet/1",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusForbidden {
		t.Errorf("expected Forbidden (403) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotDeleteInvalidTweet(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"admin",
		http.MethodDelete,
		"/tweet/100",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected Not Found (404) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotReactInvalidTweet(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"admin",
		http.MethodPut,
		"/tweet/asdfasdf/react/party",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected Not Found (404) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotRetweetInvalidTweet(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"admin",
		http.MethodDelete,
		"/tweet/-1/retweet",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected Not Found (404) got %s (%d)", response.Status, response.StatusCode)
	}
}
