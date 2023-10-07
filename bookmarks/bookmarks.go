package bookmarks

import (
	"fmt"
	"peargram/database"
	"peargram/models"
)

func GetBookmarks(username string) []models.Post {
	var bookmarkedPosts []models.Post

	DB := database.ConnectDB()
	DB.Select(&bookmarkedPosts, "SELECT * FROM posts WHERE id IN (SELECT postID FROM bookmarks WHERE actor=?)", username)

	return bookmarkedPosts
}

func HasBookmarked(username string, postID int) bool {
	var hasBookmarked bool
	hasBookmarked = false

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT EXISTS(SELECT * FROM bookmarks WHERE actor=? AND postID=?)", username, postID).Scan(&hasBookmarked)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return hasBookmarked
}
