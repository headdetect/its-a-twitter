package controller_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func TestHandleTestUser(t *testing.T) {
	response, _ := makeTestRequest(t, http.MethodGet, "/user/profile/basic", nil)
	var actualResponse controller.SmallProfileUserResponse
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

	// Verify user //
	if actualResponse.User.Username != "test" {
		t.Errorf(
			"expected test got %s",
			actualResponse.User.Username,
		)
	}

	if actualResponse.FollowerCount != 2 {
		t.Errorf(
			"expected 2 got %d",
			actualResponse.FollowerCount,
		)
	}

	if actualResponse.FollowingCount != 1 {
		t.Errorf(
			"expected 1 got %d",
			actualResponse.FollowingCount,
		)
	}

	for _, tweet := range actualResponse.Timeline.Tweets {
		if !tweet.UserRetweeted {
			t.Errorf("expected false got true")
		}
	}
}

func TestHandleTestUserLoggedIn(t *testing.T) {
	response, _ := makeAuthenticatedTestRequest(t, "admin", http.MethodGet, "/user/profile/test", nil)
	var actualResponse controller.ProfileUserResponse
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

	// Verify user //
	if actualResponse.User.Username != "test" {
		t.Errorf(
			"expected test got %s",
			actualResponse.User.Username,
		)
	}

	if len(actualResponse.Followers) != 0 {
		t.Errorf(
			"expected 0 got %d",
			len(actualResponse.Followers),
		)
	}

	if len(actualResponse.Following) != 2 {
		t.Errorf(
			"expected 2 got %d",
			len(actualResponse.Following),
		)
	}
}

func TestHandleOwnUser(t *testing.T) {
	response, _ := makeAuthenticatedTestRequest(t, "test", http.MethodGet, "/user/self", nil)

	var actualResponse controller.ProfileUserResponse
	parseTestResponse(t, response, &actualResponse)

	if actualResponse.User.Username != "test" {
		t.Errorf(
			"expected 'test' got '%s'",
			actualResponse.User.Username,
		)
	}
}

func TestHandleFollow(t *testing.T) {
	// Verify can put //
	response, _ := makeAuthenticatedTestRequest(
		t,
		"test",
		http.MethodPut,
		"/user/profile/lurker/follow",
		nil,
	)

	// Verify updated //
	response, _ = makeAuthenticatedTestRequest(t, "test", http.MethodGet, "/user/self", nil)
	var actualResponse controller.ProfileUserResponse
	parseTestResponse(t, response, &actualResponse)

	exists := false

	for _, following := range actualResponse.Following {
		if following.Username == "lurker" {
			exists = true
			break
		}
	}

	if !exists {
		t.Error("Expected 'lurker' to exist")
	}

	// Verify can delete //
	makeAuthenticatedTestRequest(
		t,
		"test",
		http.MethodDelete,
		"/user/profile/lurker/follow",
		nil,
	)

	// Verify updated //
	response, _ = makeAuthenticatedTestRequest(t, "test", http.MethodGet, "/user/self", nil)
	parseResponse(response, &actualResponse)

	exists = false

	for _, following := range actualResponse.Following {
		if following.Username == "lurker" {
			exists = true
			break
		}
	}

	if exists {
		t.Error("Expected 'lurker' to not exist")
	}
}

func TestHandleRegister(t *testing.T) {
	makeTestRequest(
		t,
		http.MethodPost,
		"/user/register",
		bytes.NewReader([]byte(`{ "username": "fake", "password": "password", "email": "fake@fake.com" }`)),
	)

	var actualResponse controller.ProfileUserResponse
	response, _ := makeAuthenticatedTestRequest(t, "fake", http.MethodGet, "/user/self", nil)
	parseResponse(response, &actualResponse)

	if actualResponse.User.Username != "fake" {
		t.Errorf(
			"expected 'fake' got '%s'",
			actualResponse.User.Username,
		)
	}
}

func TestHandleUserLogin(t *testing.T) {
	actualResponse, httpResponse, err := login(controller.LoginRequest{Username: "test", Password: "password"})

	if err != nil {
		t.Errorf("Error logging in %k\n", err)
	}

	if httpResponse.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", httpResponse.Status, httpResponse.StatusCode)
	}

	if actualResponse.User.Username != "test" {
		t.Errorf("expected test got %v", string(actualResponse.User.Username))
	}
}

// Invalids //

func TestHandleWrongUserLogin(t *testing.T) {
	_, httpResponse, err := login(controller.LoginRequest{Username: "sdfsdf", Password: "asdfasdfasdf"})

	if httpResponse.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected Unauthorized (401) got %s (%d)", httpResponse.Status, httpResponse.StatusCode)
	}

	if err != nil && httpResponse.StatusCode != http.StatusUnauthorized {
		t.Errorf("Error Attempting to log in %k\n", err)
	}
}

func TestHandleInvalidOwnUser(t *testing.T) {
	// Unauthenticated //
	response, _ := makeRequest(http.MethodGet, "/user/self", nil)

	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected Unauthorized (401) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleInvalidUser(t *testing.T) {
	response, _ := makeRequest(http.MethodGet, "/user/profile/asdfasdf", nil)

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected Not Found (404) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleInvalidUser2(t *testing.T) {
	response, _ := makeRequest(http.MethodGet, "/user/profile/", nil)

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected Not Found (404) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleInvalidUserLogin(t *testing.T) {
	response, _ := makeRequest(
		http.MethodPost,
		"/user/login",
		bytes.NewReader([]byte(`{ "username": "not-a-real-user", "password":"not-a-real-user" }`)),
	)

	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected Unauthorized (401) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotFollowSelf(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"test",
		http.MethodPut,
		"/user/profile/test/follow",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("expected Bad Request (400) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotFollowInvalidUser(t *testing.T) {
	response, _, err := makeAuthenticatedRequest(
		"test",
		http.MethodPut,
		"/user/profile/notarealuser/follow",
		nil,
	)

	if err != nil {
		t.Errorf("Error making authenticated request. %k\n", err)
	}

	if response.StatusCode != http.StatusNotFound {
		t.Errorf("expected Not Found (404) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotRegisterSameUsername(t *testing.T) {
	response, _ := makeRequest(
		http.MethodPost,
		"/user/register",
		bytes.NewReader([]byte(`{ "username": "test", "password": "password", "email": "not-test@example.com" }`)),
	)

	if response.StatusCode != http.StatusConflict {
		t.Errorf("expected Conflict (409) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotRegisterHaveInvalidUsername(t *testing.T) {
	response, _ := makeRequest(
		http.MethodPost,
		"/user/register",
		bytes.NewReader([]byte(`{ "username": "a user with spaces", "password": "password", "email": "not-test@example.com" }`)),
	)

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("expected Bad Request (400) got %s (%d)", response.Status, response.StatusCode)
	}
}
