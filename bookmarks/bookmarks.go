package bookmarks

import (
	"fmt"
	"peargram/database"
	"peargram/models"
	"peargram/users"
	"time"
)

func GetBookmarkedPosts(username string) []models.Post {
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

func ToggleBookmark(username string, postID int) error { // TODO: Check if user exists
	if !users.UserExists(username) {
		return fmt.Errorf("User does not exist")
	}

	DB := database.ConnectDB()

	if HasBookmarked(username, postID) {
		_, err := DB.Exec("DELETE FROM bookmarks WHERE actor=? AND postID=?", username, postID)
		if err != nil {
			return err
		}
	} else {
		date := time.Now().Unix()
		_, err := DB.Exec("INSERT INTO bookmarks (actor, postID, date) VALUES(?, ?, ?)", username, postID, date)
		if err != nil {
			return err
		}
	}

	return nil
}
