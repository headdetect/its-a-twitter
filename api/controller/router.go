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

type routeInfo struct {
	method  string
	path    string // Will be compiled to regex //
	handler http.HandlerFunc

	// The logger middleware will apply to all of them by default //
	middlewares []MiddlewareFunc
}

// A regex compiled version of the struct above //
type Route struct {
	method  string
	path    *regexp.Regexp
	handler http.HandlerFunc

	middlewares []MiddlewareFunc
}

func ServeWithRoutes(routes []Route) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()

		// Default middleware that is always enabled //
		logMiddle := LogMiddleware()
		corsMiddle := CorsMiddleware()

		for _, route := range routes {
			matches := route.path.FindStringSubmatch(request.URL.Path)

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
}

func MakeRoutes() (routes []Route) {
	routeInfos := []routeInfo{
		// Root //
		{method: "GET", path: "/", handler: HandleRoot},

		// Timeline //
		{method: "GET", path: "/timeline", handler: HandleTimeline, middlewares: []MiddlewareFunc{AuthMiddleware()}},

		// Users //
		{method: "GET", path: "/user/self", handler: HandleOwnUser, middlewares: []MiddlewareFunc{AuthMiddleware()}},
		{method: "POST", path: "/user/login", handler: HandleUserLogin},
		{method: "POST", path: "/user/register", handler: HandleUserRegister},
		{method: "GET", path: "/user/profile/([^/]+)", handler: HandleUser},
		{method: "PUT", path: "/user/profile/([^/]+)/follow", handler: HandleFollowUser, middlewares: []MiddlewareFunc{AuthMiddleware()}},
		{method: "DELETE", path: "/user/profile/([^/]+)/follow", handler: HandleUnFollowUser, middlewares: []MiddlewareFunc{AuthMiddleware()}},

		// Tweets //
		{method: "GET", path: "/tweet/([\\d^/]+)", handler: HandleGetTweet}, // No auth required //
		{method: "POST", path: "/tweet", handler: HandlePostTweet, middlewares: []MiddlewareFunc{AuthMiddleware()}},
		{method: "DELETE", path: "/tweet/([\\d^/]+)", handler: HandleDeleteTweet, middlewares: []MiddlewareFunc{AuthMiddleware()}},
		{method: "PUT", path: "/tweet/([\\d^/]+)/retweet", handler: HandleRetweet, middlewares: []MiddlewareFunc{AuthMiddleware()}},
		{method: "DELETE", path: "/tweet/([\\d^/]+)/retweet", handler: HandleRemoveRetweet, middlewares: []MiddlewareFunc{AuthMiddleware()}},
		{method: "PUT", path: "/tweet/([\\d^/]+)/react", handler: HandleReactTweet, middlewares: []MiddlewareFunc{AuthMiddleware()}},
		{method: "DELETE", path: "/tweet/([\\d^/]+)/react", handler: HandleRemoveReactTweet, middlewares: []MiddlewareFunc{AuthMiddleware()}},
	}

	routes = make([]Route, len(routeInfos))

	for i, info := range routeInfos {
		routes[i] = Route{
			method:      info.method,
			path:        regexp.MustCompile(fmt.Sprintf("^%s$", info.path)),
			handler:     info.handler,
			middlewares: info.middlewares,
		}
	}

	return
}

func StartRouter() {
	routes := MakeRoutes()
	routerServe := ServeWithRoutes(routes)
	port, _ := utils.GetIntOrDefault("API_PORT", 5555)
	handler := http.HandlerFunc(routerServe)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}
