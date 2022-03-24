package controller_test

import (
	"bytes"
	"encoding/json"
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

var CurrentUser *models.User

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


func Login() (controller.LoginResponse, error) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/user/login", 
		bytes.NewReader([]byte(`{ "Username": "admin", "Password":"password" }`)),
	)
	
	controller.HandleUserLogin(writer, request)

	response := writer.Result()
	defer response.Body.Close()

	var expectedResponse controller.LoginResponse

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return expectedResponse, err
	}

	err = json.Unmarshal(body, &expectedResponse)

	return expectedResponse, err
}

func AuthenticatedRequest(method, target string, body io.Reader) (*http.Request, error) {
	var loginResponse controller.LoginResponse

	if CurrentUser == nil {
		response, err := Login()

		if err != nil {
			return nil, err
		}

		loginResponse = response
		CurrentUser = response.User
	}

	request := httptest.NewRequest(method, target, body)
	request.Header.Add("AuthToken", loginResponse.AuthToken)
	request.Header.Add("Username", loginResponse.User.Username)

	return request, nil
}