package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/headdetect/its-a-twitter/api/controller"
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

// TODO: Handle auth checking here
// TODO: Check for HTTP Method here

func StartRouter() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)
	middle := middleware(logger)
	port, _ := utils.GetIntOrDefault("API_PORT", 5555)
	routeMux := http.NewServeMux()

	routeMux.HandleFunc("/", middle(controller.HandleRoot))

	// TODO: Change the paths once the function registration is done
	routeMux.HandleFunc("/user", middle(controller.HandleUser))
	routeMux.HandleFunc("/user/login", middle(controller.HandleUserLogin))
	routeMux.HandleFunc("/user/new", middle(controller.HandleUserRegister))

	routeMux.HandleFunc("/timeline", middle(controller.HandleTimeline))
	routeMux.HandleFunc("/tweet", middle(controller.HandlePostTweet))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), routeMux))
}
