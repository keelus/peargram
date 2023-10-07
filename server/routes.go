package server

import (
	"net/http"
	"peargram/posts"
	"peargram/server/handlers/apiHandler"
	"peargram/server/handlers/apiHandler/messagesApi"
	"peargram/server/handlers/apiHandler/notificationsApi"
	"peargram/server/handlers/apiHandler/postsApi"
	"peargram/server/handlers/apiHandler/usersApi"
	"peargram/server/handlers/paneHandler"
	"peargram/server/websocketing/handler"
	"peargram/users"
	"peargram/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	_ "modernc.org/sqlite"
)

func SetupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	cookieAge := 60 * 60 * 24 * 10 // 10 days

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path: "/",
		// Secure:   true,
		Secure: false,
		// HttpOnly: true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieAge,
	})
	r.Use(sessions.Sessions("PeargramSession", store))

	r.HTMLRender = Renderer()
	r.Use(gin.Recovery())

	r.Static("/assets", "web/assets")

	r.GET("/ws", handler.WebSocketConnect)

	r.NoRoute(
		func(c *gin.Context) {
			utils.DOMError(c, http.StatusNotFound)
		},
	)

	mainGroup := r.Group("/")
	{
		mainGroup.Use(RequireAccount)
		mainGroup.GET("/", paneHandler.IndexPane)
		mainGroup.GET("/search", paneHandler.SearchPane)
		mainGroup.GET("/messages", paneHandler.MessagesPane)
		mainGroup.GET("/messages/:username", paneHandler.MessagesPane)
		mainGroup.GET("/notifications", paneHandler.NotificationsPane)
		mainGroup.GET("/profile", paneHandler.ProfilePane)
		mainGroup.GET("/profile/:username", paneHandler.ProfilePane)

		mainGroup.GET("/post/:id", paneHandler.Post)

		mainGroup.GET("/settings", paneHandler.SettingsPane)
		mainGroup.GET("/activity", paneHandler.ActivityPane)
		mainGroup.GET("/saved", paneHandler.SavedPane)
	}

	authGroup := r.Group("/auth")
	{
		authGroup.GET("/signin", RequireNotLoggedIn, paneHandler.Signin)
		authGroup.GET("/signup", RequireNotLoggedIn, paneHandler.Signup)
		authGroup.GET("/endSignup", RequireLoggedIn, RequireNotUsernameSet, paneHandler.EndSignup)
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/postPreview/:id", APIRequireAccount, posts.GetPostPreview)
		apiGroup.GET("/searchUsers/:username", APIRequireAccount, apiHandler.GETSearchUsers)
		apiGroup.GET("toggleFollow/:username", APIRequireAccount, apiHandler.GETToggleFollow)

		apiGroup.POST("signinEndpoint", RequireNotLoggedIn, apiHandler.POSTSigninEndpoint)
		apiGroup.POST("signupEndpoint", RequireNotLoggedIn, apiHandler.POSTSignupEndpoint)

		apiGroup.POST("endSignup", RequireLoggedIn, RequireNotUsernameSet, apiHandler.POSTEndSignup)
		apiGroup.GET("/logout", apiHandler.Logout)

		apiMessagesGroup := apiGroup.Group("/messages")
		{
			apiMessagesGroup.Use(APIRequireAccount)
			apiMessagesGroup.GET("/getMessages", messagesApi.GETMessages)
			apiMessagesGroup.POST("/sendMessage", messagesApi.POSTSendMessage)
		}

		apiNotificationsGroup := apiGroup.Group("/notifications")
		{
			apiNotificationsGroup.Use(APIRequireAccount)
			apiNotificationsGroup.GET("/getNotifications/:username", notificationsApi.GETNotifications)
		}

		apiUsersGroup := apiGroup.Group("/users")
		{
			apiUsersGroup.Use(APIRequireAccount)
			apiUsersGroup.GET("/getUserDetails")
			apiUsersGroup.GET("/getAvatar/:username", usersApi.GETAvatar)
		}

		apiLikesGroup := apiGroup.Group("/posts")
		{
			apiLikesGroup.Use(APIRequireAccount)
			apiLikesGroup.GET("/toggleLike", postsApi.GETToggleLike)
			apiLikesGroup.GET("/toggleBookmark", postsApi.GetToggleBookmark)
		}
	}
	return r
}

func RequireLoggedIn(c *gin.Context) {
	if !users.IsUserLoggedIn(c) {
		c.Redirect(http.StatusFound, "/auth/signin")
		c.Abort()
	} else {
		c.Next()
	}
}
func RequireNotLoggedIn(c *gin.Context) {
	if users.IsUserLoggedIn(c) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.Abort()
	} else {
		c.Next()
	}
}
func RequireUsernameSet(c *gin.Context) {
	if !users.IsUsernameSet(c) {
		c.Redirect(http.StatusFound, "/auth/endSignup")
		c.Abort()
	} else {
		c.Next()
	}
}
func RequireNotUsernameSet(c *gin.Context) {
	if users.IsUsernameSet(c) {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	} else {
		c.Next()
	}
}

func RequireAccount(c *gin.Context) {
	usernameSet := users.IsUsernameSet(c)
	userLoggedIn := users.IsUserLoggedIn(c)
	if userLoggedIn && usernameSet {
		c.Next()
	} else {
		if !userLoggedIn {
			c.Redirect(http.StatusFound, "/auth/signin")
		} else {
			c.Redirect(http.StatusFound, "/auth/endSignup")
		}
		c.Abort()
	}
}
func APIRequireAccount(c *gin.Context) {
	if users.IsUserLoggedIn(c) && users.IsUsernameSet(c) {
		c.Next()
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "You're not logged in!"})
		c.Abort()
	}
}
