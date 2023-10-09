package postsApi

import (
	"net/http"
	"peargram/posts"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GETToggleLike(c *gin.Context) {
	session := sessions.Default(c)

	username := session.Get("Username").(string)
	postID, err := strconv.Atoi(c.Query("postID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected post ID."})
		return
	}

	err = posts.ToggleLike(username, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while (un)liking the post."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": posts.GetLikes(postID)})
}

func GetToggleBookmark(c *gin.Context) {
	session := sessions.Default(c)

	username := session.Get("Username").(string)
	postID, err := strconv.Atoi(c.Query("postID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected post ID."})
		return
	}

	err = posts.ToggleBookmark(username, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while (un)bookmarking the post."})
		return
	}

	c.Status(http.StatusOK)
}
