package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := gin.Default()
	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/",
		Index,
	},

	{
		"V1PollsGet",
		http.MethodGet,
		"/v1/polls",
		V1PollsGet,
	},

	{
		"V1PollsPost",
		http.MethodPost,
		"/v1/polls",
		V1PollsPost,
	},

	{
		"V1PollsQuestionGet",
		http.MethodGet,
		"/v1/polls/question",
		V1PollsQuestionGet,
	},

	{
		"V1PollsQuestionPost",
		http.MethodPost,
		"/v1/polls/question",
		V1PollsQuestionPost,
	},

	{
		"V1UserBlacklistPost",
		http.MethodPost,
		"/v1/user/blacklist",
		V1UserBlacklistPost,
	},

	{
		"V1UserFollowerPost",
		http.MethodPost,
		"/v1/user/follower",
		V1UserFollowerPost,
	},

	{
		"V1UserInfoGet",
		http.MethodGet,
		"/v1/user/info",
		V1UserInfoGet,
	},

	{
		"V1UserLoginGet",
		http.MethodGet,
		"/v1/user/login",
		V1UserLoginGet,
	},

	{
		"V1BignosecatGet",
		http.MethodGet,
		"/v1/bignosecat",
		V1BignosecatGet,
	},
}
