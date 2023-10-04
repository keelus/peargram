package posts

import (
	"fmt"
	"net/http"
	"peargram/database"
	"peargram/follows"
	"peargram/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func GetPost(id int) models.Post {
	var post models.Post

	DB := database.ConnectDB()
	DB.QueryRowx("SELECT * FROM posts WHERE id=?", id).StructScan(&post)

	// post.UserDetails = users.GetUserDetails(post.Username)

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

func GetPosts(username string) []models.Post {
	var posts []models.Post

	DB := database.ConnectDB()
	DB.Select(&posts, "SELECT * FROM posts WHERE username=?", username)

	// for i := 0; i < len(posts); i++ {
	// 	posts[i].UserDetails = users.GetUserDetails(posts[i].Username)
	// }

	return posts
}

func GetFeedPosts(username string) []models.Post { // Posts of prople this user follows
	var posts []models.Post

	following := follows.GetFollowings(username)

	// TODO: Improve this
	DB := database.ConnectDB()
	query, args, err := sqlx.In("SELECT * FROM posts WHERE username IN (?) ORDER BY date DESC", following)
	query = DB.Rebind(query)
	rows, err := DB.Query(query, args...)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.Username, &post.Content, &post.Likes, &post.Comments, &post.Date)
		if err != nil {
			fmt.Println(err)
			return posts
		}
		posts = append(posts, post)
	}

	// for i := 0; i < len(posts); i++ {
	// 	posts[i].UserDetails = users.GetUserDetails(posts[i].Username)
	// }

	return posts
}

func GetPostAmount(username string) int {
	var amount int

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT COUNT(*) FROM posts WHERE username=?", username).Scan(&amount)
	if err != nil {
		fmt.Println(err)
	}

	return amount
}
