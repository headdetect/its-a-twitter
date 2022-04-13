package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/headdetect/its-a-twitter/api/controller"
)

func BenchmarkTweets(b *testing.B) {
	loginResponse, _, err := login(controller.LoginRequest{
		Username: "lurker",
		Password: "password",
	})

	if err != nil {
		b.Errorf("Error logging in")
		return
	}

	for i := 0; i < b.N; i++ {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/tweet", bytes.NewReader([]byte(`{ "text": "Tweet Tweet" }`)))

		request.Header.Add("AuthToken", loginResponse.AuthToken)
		request.Header.Add("Username", loginResponse.User.Username)

		ControllerServe(writer, request)
	}
}
