package controller_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func TestHandleUser(t *testing.T) {
	writer := httptest.NewRecorder()
	request, err := AuthenticatedRequest(http.MethodGet, "/user", nil)

	if err != nil {
		t.Errorf("Error while attempting to authenticate the request %k", err)
	}

	controller.HandleUser(writer, request)

	response := writer.Result()
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var expectedResponse controller.UserResponse

	if err != nil {
		log.Fatal("Error caught from response body")
	}

	err = json.Unmarshal(body, &expectedResponse)

	if err != nil {
		log.Fatalf("Error parsing TestHandleUser json response. Got: %s\n", string(body))
	}

	if expectedResponse.User.Username != controller.CurrentUser.Username {
		t.Errorf(
			"expected %s got %s", 
			expectedResponse.User.Username, 
			controller.CurrentUser.Username,
		)
	}
}

func TestHandleUserLogin(t *testing.T) {
	expectedResponse, err := Login()

	if err != nil {
		t.Errorf("Error logging in %k\n", err)
	}

	if expectedResponse.User.Username != "admin" {
		t.Errorf("expected admin got %v", string(expectedResponse.User.Username))
	}
}