package users

import (
	"fmt"
	"peargram/database"
	"peargram/follows"
	"peargram/models"
	"peargram/posts"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAvatar(username string) string {
	DB := database.ConnectDB()
	var avatar string
	DB.QueryRow("SELECT avatar FROM userDetails WHERE username=?", username).Scan(&avatar)

	return avatar
}

func GetUserDetails(username string) models.UserDetails {
	var userDetails models.UserDetails

	DB := database.ConnectDB()
	DB.QueryRowx("SELECT * FROM userDetails WHERE username=?", username).StructScan(&userDetails)

	userDetails.PostAmount = posts.GetPostAmount(username)
	userDetails.FollowerAmount = follows.GetFollowerAmount(username)
	userDetails.FollowingAmount = follows.GetFollowingAmount(username)

	return userDetails
}

func IsUserLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	if session.Get("Username") == nil {
		fmt.Println("User not logged in")
		return false
	}
	if session.Get("Username").(string) == "" {
		fmt.Println("User not logged in")
		return false
	}

	fmt.Println("User logged in")
	return true
}

func IsUsernameSet(c *gin.Context) bool {
	if !IsUserLoggedIn(c) {
		return false
	}
	session := sessions.Default(c)
	if session.Get("Username").(string) == "UNDEFINED" {
		return false
	}

	return true
}
