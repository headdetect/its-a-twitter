package controller_test

import (
	"net/http"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func TestHandleGetTweet(t *testing.T) {
	response, _ := makeRequest(http.MethodGet, "/tweet/1", nil)
  var actualResponse controller.SingleTweetResponse
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

	if actualResponse.Tweet.Id != 1 {
		t.Errorf(
			"expected id:1 got id:%d\n", 
			actualResponse.Tweet.Id, 
		)
	}

	if actualResponse.Tweet.Text != "First tweet ever" {
		t.Errorf(
			"expected 'First tweet ever' got '%s'\n", 
			actualResponse.Tweet.Text, 
		)
	}

	count, ok := actualResponse.ReactionCount["🎉"]

	if !ok {
		t.Errorf("expected '🎉' to exist\n")
	}

	if count != 2 {
		t.Errorf("expected 2 got %d\n", count)	
	}

	if actualResponse.RetweetCount != 2 {
		t.Errorf("expected 2 got %d\n", actualResponse.RetweetCount)	
	}
}

func TestHandlePostTweet(t *testing.T) {

}

func TestHandleRetweet(t *testing.T) { 
	response, _, err := makeAuthenticatedRequest("lily", http.MethodPut, "/tweet/1/retweet", nil)

	if err != nil {
		t.Fatalf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected Status Code (200) got (%d)\n", response.StatusCode)
	}

	response, _ = makeRequest(http.MethodGet, "/tweet/1", nil)

	var actualResponse controller.SingleTweetResponse
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected Status Code (200) got (%d)\n", response.StatusCode)
	}

	if actualResponse.RetweetCount != 3 {
		t.Errorf("Expecting 3 got %d\n", actualResponse.RetweetCount)
	}

	// TODO: Figure this out ^
}


func TestHandleReactTweet(t *testing.T) {
	
}

func TestHandleDeleteRetweet(t *testing.T) {
	
}


func TestHandleDeleteReactTweet(t *testing.T) {
	
}

// Non happy path //
func TestHandleCannotRetweetOwnTweet(t *testing.T) {
	
}

func TestHandleCannotReactOwnTweet(t *testing.T) {
	
}

func TestHandleCannotDeleteOtherTweet(t *testing.T) {
	
}

func TestHandleCannotDeleteInvalidTweet(t *testing.T) {
	
}

func TestHandleCannotReactInvalidTweet(t *testing.T) {
	
}

func TestHandleCannotRetweetInvalidTweet(t *testing.T) {
	
}