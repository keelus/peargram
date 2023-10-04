package usersApi

import (
	"net/http"
	"peargram/users"

	"github.com/gin-gonic/gin"
)

func GETAvatar(c *gin.Context) {
	username := c.Param("username")
	avatar := users.GetAvatar(username)
	c.Data(http.StatusOK, "image/jpeg", []byte(avatar))
}
