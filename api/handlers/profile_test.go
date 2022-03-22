package handlers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/headdetect/its-a-twitter/api/handlers"
	"github.com/headdetect/its-a-twitter/api/store"
	"github.com/joho/godotenv"
)

var currentLoginResponse handlers.LoginResponse


func TestMain(m *testing.M) {
	log.Println("Loading env")
	err := godotenv.Load("../.env")
	
	if err != nil {
		log.Fatalf("Error loading env. Error %k\n", err)
	}

	log.Println("Loading database")
	store.LoadDatabaseFromFile("../store/store.db", "../store/initial.sql")

	m.Run()

	os.Exit(0)
}

func TestHandleUser(t *testing.T) {
	if currentLoginResponse == (handlers.LoginResponse{}) {
		TestHandleUserLogin(t)
	}

	if currentLoginResponse == (handlers.LoginResponse{}) {
		t.Fatal("Failed to get auth token from login")
	}

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	request.Header.Add("AuthToken", currentLoginResponse.AuthToken)
	request.Header.Add("Username", currentLoginResponse.User.Username)
	
	handlers.HandleUser(writer, request)

	response := writer.Result()
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var expectedResponse handlers.UserResponse

	if err != nil {
		log.Fatal("Error caught from response body")
	}

	err = json.Unmarshal(body, &expectedResponse)

	if err != nil {
		log.Fatalf("Error parsing TestHandleUser json response. Got: %s\n", string(body))
	}

	if expectedResponse.User.Username != currentLoginResponse.User.Username {
		t.Errorf(
			"expected %s got %s", 
			expectedResponse.User.Username, 
			currentLoginResponse.User.Username,
		)
	}

	if expectedResponse.User.DisplayName != currentLoginResponse.User.DisplayName {
		t.Errorf(
			"expected %s got %s", 
			expectedResponse.User.DisplayName, 
			currentLoginResponse.User.DisplayName,
		)
	}
}

func TestHandleUserLogin(t *testing.T) {
	requestBody := []byte(`{ "Username": "admin", "Password":"password" }`)

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewReader(requestBody))
	
	handlers.HandleUserLogin(writer, request)

	response := writer.Result()
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	var expectedResponse handlers.LoginResponse

	if err != nil {
		log.Fatal("Error caught from response body")
	}

	err = json.Unmarshal(body, &expectedResponse)

	if err != nil {
		log.Fatalf("Error parsing TestHandleUserLogin json response. Got: %s\n", string(body))
	}

	if expectedResponse.User.Username != "admin" {
		t.Errorf("expected admin got %v", string(expectedResponse.User.Username))
	}

	currentLoginResponse = expectedResponse
}