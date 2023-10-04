package messagesApi

import (
	"net/http"
	"peargram/messages"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GETMessages(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	username := c.Param("username")

	chat := messages.GetChat(currentUsername, username)

	c.JSON(http.StatusOK, chat)
}

func GETSendMessage(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	target := c.Param("target")
	content := c.Param("content")

	messages.SendMessage(currentUsername, target, content)

}
