package handlers

import "net/http"

func HandleRoot(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"hello": "world"}`))
}