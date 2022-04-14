package controller

import "net/http"

func HandleRoot(writer http.ResponseWriter, request *http.Request) {
	JsonResponse(writer, []byte(`{"github": "https://github.com/headdetect/its-a-twitter"}`))
}
