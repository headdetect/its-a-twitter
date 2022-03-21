package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/headdetect/its-a-twitter/api/handlers"
	"github.com/headdetect/its-a-twitter/api/utils"
)

func middleware(logger *log.Logger) func(http.HandlerFunc) http.HandlerFunc  {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			f(w, r)
		}
	}
}

func StartRouter() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)
	middle := middleware(logger)
	port, _ := utils.GetIntOrDefault("API_PORT", 5555)
	routeMux := http.NewServeMux()

	routeMux.HandleFunc("/", middle(handlers.HandleRoot))
	routeMux.HandleFunc("/user", middle(handlers.HandleUser))
	routeMux.HandleFunc("/user/login", middle(handlers.HandleUserLogin))
	routeMux.HandleFunc("/user/new", middle(handlers.HandleUserRegister))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routeMux))
}
