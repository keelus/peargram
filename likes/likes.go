package likes

import (
	"fmt"
	"peargram/database"
	"peargram/models"
)

func GetLikes(username string) []models.Post {
	var likedPosts []models.Post

	DB := database.ConnectDB()
	DB.Select(&likedPosts, "SELECT * FROM posts WHERE id IN (SELECT postID FROM likes WHERE actor=?)", username)

	return likedPosts
}

func HasLiked(username string, postID int) bool {
	var hasLiked bool
	hasLiked = false

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT EXISTS(SELECT * FROM likes WHERE actor=? AND postID=?)", username, postID).Scan(&hasLiked)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return hasLiked
}
