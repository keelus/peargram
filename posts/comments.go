package posts

import (
	"peargram/database"
	"peargram/models"
)

func GetComments(postID int) []models.Comment {
	var comments []models.Comment
	// comments := make([]models.Comment, 0)

	DB := database.ConnectDB()
	DB.Select(&comments, "SELECT * FROM comments WHERE postID=?", postID)

	return comments
}

func GetCommentAmount(postID int) uint {
	commentAmount := uint(0)

	DB := database.ConnectDB()
	DB.QueryRow("SELECT COUNT(*) FROM comments WHERE postID=?", postID).Scan(&commentAmount)

	return commentAmount
}

func GetCommented(username string) []models.Comment {
	comments := make([]models.Comment, 0)

	DB := database.ConnectDB()
	DB.Select(&comments, "SELECT * FROM comments WHERE id IN (SELECT postID FROM comments WHERE actor=?)", username)

	return comments
}
