package postsApi

import (
	"net/http"
	"peargram/bookmarks"
	"peargram/likes"
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

	err = likes.ToggleLike(username, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while (un)liking the post."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes.GetLikes(postID)})
}

func GetToggleBookmark(c *gin.Context) {
	session := sessions.Default(c)

	username := session.Get("Username").(string)
	postID, err := strconv.Atoi(c.Query("postID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected post ID."})
		return
	}

	err = bookmarks.ToggleBookmark(username, postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while (un)bookmarking the post."})
		return
	}

	c.Status(http.StatusOK)
}
