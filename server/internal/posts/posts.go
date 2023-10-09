package posts

import (
	"fmt"
	"net/http"
	"peargram/database"
	"peargram/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPost(id int) models.Post {
	var post models.Post

	DB := database.ConnectDB()
	DB.QueryRowx("SELECT * FROM posts WHERE id=?", id).Scan(&post.ID, &post.Username, &post.Content, &post.Date)

	likes := GetLikes(id)
	post.Likes = likes

	comments := GetComments(id)
	post.Comments = comments

	post.CommentAmount = uint(len(comments))

	return post
}

func GetPostPreview(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	post := GetPost(idInt)

	c.Data(http.StatusOK, "image/jpeg", []byte(post.Content))
}

func GetUserPosts(username string) []models.Post {
	var posts []models.Post

	DB := database.ConnectDB()
	DB.Select(&posts, "SELECT * FROM posts WHERE username=? ORDER BY date DESC;", username)

	for i := 0; i < len(posts); i++ {
		posts[i].Likes = GetLikes(posts[i].ID)
		posts[i].CommentAmount = GetCommentAmount(posts[i].ID)
	}

	return posts
}

func GetFeedPosts(username string) []models.Post { // Posts of prople this user follows
	var posts []models.Post

	// TODO: Improve this
	DB := database.ConnectDB()

	DB.Select(&posts, "SELECT * FROM posts WHERE username IN (SELECT target FROM follows WHERE actor=?) ORDER BY date DESC", username)

	for i := 0; i < len(posts); i++ {
		posts[i].Likes = GetLikes(posts[i].ID)
		posts[i].CommentAmount = GetCommentAmount(posts[i].ID)
	}

	return posts
}

func GetUserPostAmount(username string) int {
	var amount int

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE username=?", username).Scan(&amount)
	if err != nil {
		fmt.Println(err)
	}

	return amount
}
