package likes

import (
	"fmt"
	"peargram/database"
	"peargram/models"
	"peargram/users"
	"time"
)

func GetLikes(postID int) uint {
	postLikes := uint(0)

	DB := database.ConnectDB()
	DB.QueryRow("SELECT COUNT(*) FROM likes WHERE postID=?", postID).Scan(&postLikes)

	return postLikes
}

func GetLikedPosts(username string) []models.Post {
	likedPosts := make([]models.Post, 0)
	if !users.UserExists(username) {
		return likedPosts
	}

	DB := database.ConnectDB()
	DB.Select(&likedPosts, "SELECT * FROM posts WHERE id IN (SELECT postID FROM likes WHERE actor=?)", username)

	return likedPosts
}

func HasLiked(username string, postID int) bool {
	if !users.UserExists(username) {
		return false
	}

	hasLiked := false

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT EXISTS(SELECT * FROM likes WHERE actor=? AND postID=?)", username, postID).Scan(&hasLiked)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return hasLiked
}

func ToggleLike(username string, postID int) error { // TODO: Check if user exists
	if !users.UserExists(username) {
		return fmt.Errorf("User does not exist")
	}

	DB := database.ConnectDB()

	if HasLiked(username, postID) {
		_, err := DB.Exec("DELETE FROM likes WHERE actor=? AND postID=?", username, postID)
		if err != nil {
			return err
		}
	} else {
		date := time.Now().Unix()
		_, err := DB.Exec("INSERT INTO likes (actor, postID, date) VALUES(?, ?, ?)", username, postID, date)
		if err != nil {
			return err
		}
	}

	return nil
}
