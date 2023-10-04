package server

import (
	"net/http"
	"peargram/posts"
	"peargram/server/handlers/apiHandler"
	"peargram/server/handlers/apiHandler/messagesApi"
	"peargram/server/handlers/apiHandler/notificationsApi"
	"peargram/server/handlers/apiHandler/usersApi"
	"peargram/server/handlers/paneHandler"
	"peargram/server/websocketing/handler"
	"peargram/users"

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

	// r.NoRoute(func(c *gin.Context) {
	// 	utils.ForceError(http.StatusNotFound, c, "")
	// })

	mainGroup := r.Group("/")
	{
		mainGroup.Use(RequireLoggedIn, RequireUsernameSet)
		mainGroup.GET("/", paneHandler.IndexPane)
		mainGroup.GET("/search", RequireLoggedIn, paneHandler.SearchPane)
		mainGroup.GET("/messages", RequireLoggedIn, paneHandler.MessagesPane)
		mainGroup.GET("/messages/:username", RequireLoggedIn, paneHandler.MessagesPane)
		mainGroup.GET("/notifications", RequireLoggedIn, paneHandler.NotificationsPane)
		mainGroup.GET("/profile/:username", RequireLoggedIn, paneHandler.ProfilePane)

		mainGroup.GET("/post/:id", paneHandler.Post)

		mainGroup.GET("temporaryWS", paneHandler.TempWS)
	}

	authGroup := r.Group("/auth")
	{
		authGroup.GET("/signin", RequireNotLoggedIn, paneHandler.Signin)
		authGroup.GET("/signup", RequireNotLoggedIn, paneHandler.Signup)
		authGroup.GET("/endSignup", RequireNotUsernameSet, paneHandler.EndSignup)
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/logout", apiHandler.Logout)
		apiGroup.GET("/postPreview/:id", RequireLoggedIn, RequireNotUsernameSet, posts.GetPostPreview)

		apiGroup.GET("/searchUsers/:username", RequireLoggedIn, RequireNotUsernameSet, apiHandler.GETSearchUsers)

		apiGroup.GET("toggleFollow/:username", RequireLoggedIn, RequireNotUsernameSet, apiHandler.GETToggleFollow)

		apiGroup.POST("signinEndpoint", apiHandler.POSTSigninEndpoint)
		apiGroup.POST("signupEndpoint", apiHandler.POSTSignupEndpoint)

		apiGroup.POST("endSignup", apiHandler.POSTEndSignup)

		apiMessagesGroup := apiGroup.Group("/messages")
		{
			apiMessagesGroup.Use(RequireLoggedIn, RequireNotUsernameSet)
			apiMessagesGroup.GET("/getMessages/:username", messagesApi.GETMessages)
			apiMessagesGroup.GET("/sendMessage/:target/:content", messagesApi.GETSendMessage) // TODO: POST
		}

		apiNotificationsGroup := apiGroup.Group("/notifications")
		{
			apiNotificationsGroup.Use(RequireLoggedIn, RequireNotUsernameSet)
			apiNotificationsGroup.GET("/getNotifications/:username", notificationsApi.GETNotifications)
		}

		apiUsersGroup := apiGroup.Group("/users")
		{
			apiUsersGroup.Use(RequireLoggedIn, RequireUsernameSet)
			apiUsersGroup.GET("/getUserDetails")
			apiUsersGroup.GET("/getAvatar/:username", usersApi.GETAvatar)
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
