package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/headdetect/its-a-twitter/api/model"
)

var (
	ACCEPTABLE_MIME_TYPES = []string{"image/jpeg", "image/jpg", "image/png", "image/gif"}
)

var Sessions map[string]model.User = make(map[string]model.User) // [authToken] = user

// TODO Make these helper funcs private

func GetCurrentUser(request *http.Request) (model.User, error) {
	authToken := request.Header.Get("Authtoken")

	// The auth username is a sort of 'public key'
	// so people can't just brute-force a token.
	// they'll also have to supply the username
	// that is associated with the token
	authUsername := request.Header.Get("Username")

	if user, ok := Sessions[authToken]; ok {
		if authUsername == user.Username {
			return user, nil
		}
	}

	return model.User{}, errors.New("No user logged in")
}

func GetPathValue(request *http.Request, index int) (string, bool) {
	fields, ok := request.Context().Value(ContextKeys).([]string)

	if index >= len(fields) || !ok {
		return "", false
	}

	return fields[index], true
}

func JsonResponse(writer http.ResponseWriter, response []byte) {
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

func BadRequestResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusBadRequest)
	JsonResponse(writer, []byte(`{ message: "Invalid request" }`))
}

func UnauthorizedResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusUnauthorized)
	JsonResponse(writer, []byte(`{ message: "Invalid auth token provided. Please log in" }`))
}

func ForbiddenResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusForbidden)
	JsonResponse(writer, []byte(`{ "message": "Nice try. You can't do that."}`))
}

func NotFoundResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusNotFound)
	JsonResponse(writer, []byte(`{ "message": "Specified resource was not found" }`))
}

func ConflictRequestResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusConflict)
	JsonResponse(writer, []byte(`{ "message": "Specified resource already exists" }`))
}

func ErrorResponse(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusInternalServerError)
	JsonResponse(writer, []byte(fmt.Sprintf(`{ "message": "We messed up somehow. Strange. We never mess up", "error": "%k" }`, err)))
}
