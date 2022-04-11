package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
	"github.com/headdetect/its-a-twitter/api/model"
	"github.com/headdetect/its-a-twitter/api/store"
	"github.com/joho/godotenv"
)

var currentUser *model.User

func TestMain(m *testing.M) {
	log.Println("Loading env")
	err := godotenv.Load("../.env")
	os.Setenv("APP_ENV", "test")
	os.Setenv("STORE_PATH", "../store")
	
	if err != nil {
		log.Fatalf("Error loading env. Error %k\n", err)
	}

	log.Println("Loading database")
	store.LoadDatabase(true)

	m.Run()

	os.Exit(0)
}

func login(request controller.LoginRequest) (controller.LoginResponse, *http.Response, error) {
	response, _ := makeRequest(
		http.MethodPost,
		"/user/login", 
		bytes.NewReader(
			[]byte(fmt.Sprintf(`{ "Username": "%s", "Password":"%s" }`, request.Username, request.Password)),
		),
	)
	
	defer response.Body.Close()

	var actualResponse controller.LoginResponse

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return actualResponse, response, err
	}

	err = json.Unmarshal(body, &actualResponse)

	return actualResponse, response, err
}

func AuthenticatedRequest(loginRequest controller.LoginRequest, request *http.Request) (*http.Request, error) {
	loginResponse, _, err := login(loginRequest)

	if err != nil {
		return nil, err
	}	

	request.Header.Add("AuthToken", loginResponse.AuthToken)
	request.Header.Add("Username", loginResponse.User.Username)

	return request, nil
}

func makeRequest(
	method string, 
	route string, 
	body io.Reader, 
) (*http.Response, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(method, route, body)

	controller.Serve(writer, request)

	response := writer.Result()
	
	return response, request
}

func makeAuthenticatedRequest(
	userName string,
	method string, 
	route string, 
	body io.Reader,
) (*http.Response, *http.Request, error) {
	loginRequest := controller.LoginRequest{
		Username: userName,
		Password: "password",
	}
	
	writer := httptest.NewRecorder()
	request, err := AuthenticatedRequest(loginRequest, httptest.NewRequest(method, route, body))

	if err != nil {
		return nil, nil, err
	}

	controller.Serve(writer, request)

	response := writer.Result()

	return response, request, nil
}

func parseResponse(response *http.Response, v any) ([]byte, error) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}
	
	return body, json.Unmarshal(body, &v)
}