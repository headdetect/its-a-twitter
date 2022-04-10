package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/headdetect/its-a-twitter/api/utils"
)

type MiddlewareFunc func(http.Handler) http.Handler

var ContextKeys struct{}

type route struct {
	method string
	path string // Will be compiled to regex //
	handler http.HandlerFunc

	// The logger middleware will apply to all of them by default //
	middlewares []MiddlewareFunc
}

var routes = []route{
	// Root //
	{ method: "GET", path: "/", handler: HandleRoot },

	// Timeline //
	{ method: "GET", path: "/timeline", handler: HandleTimeline, middlewares: []MiddlewareFunc{ AuthMiddleware() } },

	// Users //
	{ method: "GET", path: "/user/profile/([^/]+)", handler: HandleUser },
	{ method: "GET", path: "/user/self", handler: HandleOwnUser, middlewares: []MiddlewareFunc{ AuthMiddleware() } },
	{ method: "POST", path: "/user/login", handler: HandleUserLogin },
	{ method: "POST", path: "/user/register", handler: HandleUserRegister },

	// Tweets //
	{ method: "GET", path: "/tweet/([^/]+)", handler: HandleGetTweet }, // No auth required //
	{ method: "POST", path: "/tweet", handler: HandlePostTweet, middlewares: []MiddlewareFunc{ AuthMiddleware() } },
	{ method: "DELETE", path: "/tweet/([^/]+)", handler: HandleDeleteTweet, middlewares: []MiddlewareFunc{ AuthMiddleware() } },
	{ method: "PUT", path: "/tweet/([^/]+)/retweet", handler: HandleRetweet, middlewares: []MiddlewareFunc{ AuthMiddleware() } },
	{ method: "DELETE", path: "/tweet/([^/]+)/retweet", handler: HandleRemoveRetweet, middlewares: []MiddlewareFunc{ AuthMiddleware() } },
	{ method: "PUT", path: "/tweet/([^/]+)/reactions", handler: HandleReactTweet, middlewares: []MiddlewareFunc{ AuthMiddleware() } },
	{ method: "DELETE", path: "/tweet/([^/]+)/reactions", handler: HandleRemoveReactTweet, middlewares: []MiddlewareFunc{ AuthMiddleware() } },
}

func Serve(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	// Default middleware that is always enabled //
	logMiddle := LogMiddleware()
	corsMiddle := CorsMiddleware()

	for _, route := range routes {
		pattern := regexp.MustCompile(fmt.Sprintf("^%s$", route.path))
		matches := pattern.FindStringSubmatch(request.URL.Path)

		if len(matches) > 0 {
			if request.Method != route.method {
				continue
			}

			log.Println(request.Context(), ContextKeys, matches[1:])

			ctx := context.WithValue(request.Context(), ContextKeys, matches[1:])
			newContextRequest := request.WithContext(ctx)

			var next http.Handler = route.handler

			for i := len(route.middlewares) - 1; i >= 0; i-- {
				ware := route.middlewares[i]
				next = ware(next)
			}

			// Default middlewares //
			wares := logMiddle(corsMiddle(next))

			wares.ServeHTTP(writer, newContextRequest)

			return
		}
	}
	
	http.NotFound(writer, request)
}

func StartRouter() {
	port, _ := utils.GetIntOrDefault("API_PORT", 5555)
	handler := http.HandlerFunc(Serve)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
