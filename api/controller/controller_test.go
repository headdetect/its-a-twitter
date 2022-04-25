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

var ControllerServe http.HandlerFunc

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

	routes := controller.MakeRoutes()
	ControllerServe = controller.ServeWithRoutes(routes)

	m.Run()

	os.Exit(0)
}

func TestOptionRequest(t *testing.T) {
	response, _ := makeRequest(
		http.MethodOptions,
		"/user/self",
		nil,
	)

	if response.Header.Get("Allow") != "GET,PUT" {
		t.Errorf("expected 'GET,PUT' got %s", response.Header.Get("Allow"))
	}
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

func AuthenticatedRequest(loginRequest controller.LoginRequest, request *http.Request) error {
	loginResponse, _, err := login(loginRequest)

	if err != nil {
		return err
	}

	request.Header.Add("AuthToken", loginResponse.AuthToken)
	request.Header.Add("Username", loginResponse.User.Username)

	return nil
}

// Functions to help with testing happy paths/valid cases //

func makeTestRequest(
	t testing.TB,
	method string,
	route string,
	body io.Reader,
) (*http.Response, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(method, route, body)

	ControllerServe(writer, request)

	response := writer.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", response.Status, response.StatusCode)
	}

	return response, request
}

func makeAuthenticatedTestRequest(
	t testing.TB,
	userName string,
	method string,
	route string,
	body io.Reader,
) (*http.Response, *http.Request) {
	loginRequest := controller.LoginRequest{
		Username: userName,
		Password: "password",
	}

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(method, route, body)
	err := AuthenticatedRequest(loginRequest, request)

	if err != nil {
		t.Errorf("Error authenticating. %k\n", err)
	}

	ControllerServe(writer, request)

	response := writer.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected OK (200) got %s (%d)", response.Status, response.StatusCode)
	}

	return response, request
}

func parseTestResponse(t testing.TB, response *http.Response, v any) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		t.Fatalf("Error reading body response. %k\nBody: %s\n", err, string(body))
	}

	err = json.Unmarshal(body, &v)

	if err != nil {
		t.Fatalf("Error parsing json response. %k\nBody: %s\n", err, string(body))
	}
}

// Functions to build own handler for response //
func makeRequest(
	method string,
	route string,
	body io.Reader,
) (*http.Response, *http.Request) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(method, route, body)

	ControllerServe(writer, request)

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
	request := httptest.NewRequest(method, route, body)
	err := AuthenticatedRequest(loginRequest, request)

	if err != nil {
		return nil, nil, err
	}

	ControllerServe(writer, request)

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
