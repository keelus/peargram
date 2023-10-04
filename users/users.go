package users

import (
	"fmt"
	"os"
	"peargram/database"
	"peargram/follows"
	"peargram/models"
	"peargram/posts"
	"peargram/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetAvatar(username string) string {
	DB := database.ConnectDB()
	var avatar string
	DB.QueryRow("SELECT avatar FROM userDetails WHERE username=?", username).Scan(&avatar)

	if avatar == "" {
		data, err := os.ReadFile(utils.GetAbsolutePath() + "/web/assets/media/images/defaultAvatar.webp")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return string(data)
	}

	return avatar
}

func UserExists(username string) bool {
	DB := database.ConnectDB()
	exists := false
	err := DB.QueryRowx("SELECT EXISTS(SELECT * FROM users WHERE username=?)", username).StructScan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func GetUserDetails(username string) (models.UserDetails, bool) {
	var userDetails models.UserDetails

	DB := database.ConnectDB()
	err := DB.QueryRowx("SELECT * FROM userDetails WHERE username=?", username).StructScan(&userDetails)
	if err != nil {
		return userDetails, false
	}

	userDetails.PostAmount = posts.GetPostAmount(username)
	userDetails.FollowerAmount = follows.GetFollowerAmount(username)
	userDetails.FollowingAmount = follows.GetFollowingAmount(username)

	return userDetails, true
}

func IsUserLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	if session.Get("Username") == nil {
		// fmt.Println("User not logged in")
		return false
	}
	if session.Get("Username").(string) == "" {
		// fmt.Println("User not logged in")
		return false
	}

	// fmt.Println("User logged in")
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
