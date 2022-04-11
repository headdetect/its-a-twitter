package controller_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func TestHandleTestUser(t *testing.T) {
	response, _ := makeRequest(http.MethodGet, "/user/profile/test", nil)
  var actualResponse controller.UserResponse
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

	// Verify tweets //

	// Verify followers //

	// Verify followed //
}

func TestHandleLurkerUser(t *testing.T) {
	response, _ := makeRequest(http.MethodGet, "/user/profile/lurker", nil)
  var actualResponse controller.UserResponse
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

	// Verify user //
	if actualResponse.User.Username != "lurker" {
		t.Errorf(
			"expected lurker got %s", 
			actualResponse.User.Username, 
		)
	}

	// Verify tweets //

	// Verify followers //

	// Verify followed //
}

func TestHandleAdminUser(t *testing.T) {
	response, _ := makeRequest(http.MethodGet, "/user/profile/admin", nil)
  var actualResponse controller.UserResponse
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

	// Verify user //
	if actualResponse.User.Username != "admin" {
		t.Errorf(
			"expected admin got %s", 
			actualResponse.User.Username, 
		)
	}

	// Verify tweets //

	// Verify followers //

	// Verify followed //
}

func TestHandleOwnUser(t *testing.T) {
	response, _, err := makeAuthenticatedRequest("test", http.MethodGet, "/user/self", nil)

	if err != nil {
		t.Errorf("Error authenticating. %k\n", err)
	}

	var actualResponse controller.UserResponse
	body, err := parseResponse(response, &actualResponse)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}

	if actualResponse.User.Username != "test" {
		t.Errorf(
			"expected 'test' got '%s'",
			actualResponse.User.Username, 
		)
	}
}

func TestHandleFollow(t *testing.T) {
	
}

func TestHandleRegister(t *testing.T) {
	
}

func TestHandleUserLogin(t *testing.T) {
	actualResponse, httpResponse, err := login(controller.LoginRequest{ Username: "test", Password: "password" })

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
	_, httpResponse, err := login(controller.LoginRequest{ Username: "sdfsdf", Password: "asdfasdfasdf" })

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
		bytes.NewReader([]byte(`{ "Username": "fake", "Password":"fake" }`)),
	)

	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected Unauthorized (401) got %s (%d)", response.Status, response.StatusCode)
	}
}

func TestHandleCannotFollowSelf(t *testing.T) {
	
}

func TestHandleCannotFollowInvalidUser(t *testing.T) {
	
}

func TestHandleCannotRegister(t *testing.T) {
	
}