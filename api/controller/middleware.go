package controller

import (
	"log"
	"net/http"
)

func AuthMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user := GetCurrentUser(request)

			if user != nil {
				next.ServeHTTP(writer, request)
				return
			}

			UnauthorizedResponse(writer)
		})
	}
}

func LogMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			log.Println(request.Method, request.URL.Path, request.RemoteAddr, request.UserAgent())
			next.ServeHTTP(writer, request)
		})
	}
}

func CorsMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			next.ServeHTTP(writer, request)
		})
	}
}