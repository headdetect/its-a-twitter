package controller

import "net/http"

func HandleRoot(writer http.ResponseWriter, request *http.Request) {
	// TODO: Use newer responses formats
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"hello": "world"}`))
}
