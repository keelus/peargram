package paneHandler

import (
	"fmt"
	"net/http"
	"strconv"

	"peargram/bookmarks"
	"peargram/follows"
	"peargram/likes"
	"peargram/messages"
	"peargram/models"
	"peargram/notifications"
	"peargram/posts"
	"peargram/users"
	"peargram/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var AllowDataRemovedPage = false

func IndexPane(c *gin.Context) {
	session := sessions.Default(c)

	fmt.Println(session.Get("Username"))
	username := session.Get("Username").(string)
	feedPosts := posts.GetFeedPosts(username)

	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "index", gin.H{
			"User":      CurrentUserData(c),
			"FeedPosts": feedPosts,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane":  "index",
			"User":      CurrentUserData(c),
			"FeedPosts": feedPosts,
		})
	}
}

func SearchPane(c *gin.Context) {
	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "search", gin.H{
			"User": CurrentUserData(c),
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane": "search",
			"User":     CurrentUserData(c),
		})
	}
}

func MessagesPane(c *gin.Context) {
	var activeChat models.Chat
	session := sessions.Default(c)

	noneSelected := true

	selfUsername := session.Get("Username").(string)
	userChats := messages.GetChatPreviews(selfUsername)
	// firstChatUsername := userChats[0].Participants[1]

	if c.Param("username") != "" {
		for _, userChat := range userChats {
			if userChat.Participants[1] == c.Param("username") {
				activeChat = messages.GetChat(selfUsername, c.Param("username"), 0)
				noneSelected = false
			}
		}
	}

	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "messages", gin.H{
			"User":         CurrentUserData(c),
			"Chats":        userChats,
			"ActiveChat":   activeChat,
			"NoneSelected": noneSelected,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane":     "messages",
			"User":         CurrentUserData(c),
			"Chats":        userChats,
			"ActiveChat":   activeChat,
			"NoneSelected": noneSelected,
		})
	}
}
func NotificationsPane(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("Username").(string)
	notis := notifications.GetNotifications(username)

	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "notifications", gin.H{
			"User":          CurrentUserData(c),
			"Notifications": notis,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane":      "notifications",
			"User":          CurrentUserData(c),
			"Notifications": notis,
		})
	}
}
func ProfilePane(c *gin.Context) {
	username := c.Param("username")

	ProfilePosts := posts.GetPosts(username)
	selfProfile := false

	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	if username == "" {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/profile/%s", currentUsername))
	}

	if username == currentUsername {
		selfProfile = true
	}

	var profileData models.UserDetails
	profileData, exists := users.GetUserDetails(username)
	if !exists {
		utils.DOMError(c, http.StatusNotFound)
		return
	}

	profileData.PostAmount = posts.GetPostAmount(username)
	profileData.FollowingAmount = follows.GetFollowingAmount(username)
	profileData.FollowerAmount = follows.GetFollowerAmount(username)

	following := follows.IsFollowing(currentUsername, username)

	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "profile", gin.H{
			"User":         CurrentUserData(c),
			"Profile":      profileData,
			"SelfProfile":  selfProfile,
			"ProfilePosts": ProfilePosts,
			"Following":    following,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane":     "profile",
			"User":         CurrentUserData(c),
			"Profile":      profileData,
			"SelfProfile":  selfProfile,
			"ProfilePosts": ProfilePosts,
			"Following":    following,
		})
	}
}

func Signin(c *gin.Context) {
	endpoint := utils.GetEnv("SIGNIN_ENDPOINT")
	googleClientID := utils.GetEnv("GOOGLE_CLIENT_ID")

	errorType := ""
	errorType = c.Query("error")
	c.HTML(http.StatusOK, "signin", gin.H{
		"ErrorType":      errorType,
		"Endpoint":       endpoint,
		"GoogleClientID": googleClientID,
	})
}
func Signup(c *gin.Context) {
	endpoint := utils.GetEnv("SIGNUP_ENDPOINT")
	googleClientID := utils.GetEnv("GOOGLE_CLIENT_ID")

	errorType := ""
	errorType = c.Query("error")
	c.HTML(http.StatusOK, "signup", gin.H{
		"ErrorType":      errorType,
		"Endpoint":       endpoint,
		"GoogleClientID": googleClientID,
	})
}

func Post(c *gin.Context) {
	var post models.Post
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	post = posts.GetPost(idInt)

	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "post", gin.H{
			"User": CurrentUserData(c),
			"Post": post,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane": "post",
			"User":     CurrentUserData(c),
			"Post":     post,
		})
	}
}

func CurrentUserData(c *gin.Context) models.UserDetails {
	var userDetails models.UserDetails
	session := sessions.Default(c)

	userDetails.Username = session.Get("Username").(string)
	return userDetails
}

func EndSignup(c *gin.Context) {

	c.HTML(http.StatusOK, "endSignup", gin.H{})
}

func SettingsPane(c *gin.Context) {
	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "settings", gin.H{
			"User": CurrentUserData(c),
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane": "settings",
			"User":     CurrentUserData(c),
		})
	}
}

func ActivityPane(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	likedPosts := likes.GetLikes(currentUsername)

	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "activity", gin.H{
			"User":       CurrentUserData(c),
			"LikedPosts": likedPosts,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane":   "activity",
			"User":       CurrentUserData(c),
			"LikedPosts": likedPosts,
		})
	}
}
func SavedPane(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	savedPosts := bookmarks.GetBookmarks(currentUsername)

	if c.Query("type") == "short" {
		c.HTML(http.StatusOK, "saved", gin.H{
			"User":       CurrentUserData(c),
			"SavedPosts": savedPosts,
		})
	} else {
		c.HTML(http.StatusOK, "base", gin.H{
			"LoadPane":   "saved",
			"User":       CurrentUserData(c),
			"SavedPosts": savedPosts,
		})
	}
}
