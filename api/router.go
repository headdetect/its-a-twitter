package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

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

type route struct {
	method string
	path string // Will be compiled to regex //
	handler http.HandlerFunc

	// TODO: Have middleware to support auth before
}

var routes = []route{
	{ method: "GET", path: "/", handler: controller.HandleRoot },

	{ method: "GET", path: "/user/profile/([^/]+)", handler: controller.HandleUser },
	{ method: "GET", path: "/user/self", handler: controller.HandleOwnUser },
	{ method: "POST", path: "/user/login", handler: controller.HandleUserLogin },
	{ method: "POST", path: "/user/register", handler: controller.HandleUserRegister },

	{ method: "GET", path: "/timeline", handler: controller.HandleTimeline },
	{ method: "POST", path: "/tweet", handler: controller.HandlePostTweet },
}

func serve(w http.ResponseWriter, r *http.Request) {
	for _, route := range routes {
		pattern := regexp.MustCompile(fmt.Sprintf("^%s$", route.path))
		matches := pattern.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			if r.Method != route.method {
				continue
			}
			
			ctx := context.WithValue(r.Context(), controller.ContextKeys, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	
	http.NotFound(w, r)
}

func StartRouter() {
	logger := log.New(os.Stdout, "", log.Lmicroseconds)
	middle := middleware(logger)
	port, _ := utils.GetIntOrDefault("API_PORT", 5555)
	handler := http.HandlerFunc(middle(serve))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
