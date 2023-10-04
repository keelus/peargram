package messagesApi

import (
	"fmt"
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

type SendMessageBody struct {
	Target  string
	Content string
}

func POSTSendMessage(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	var requestBody SendMessageBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		fmt.Println(err)
		return
	}

	target := requestBody.Target
	content := requestBody.Content

	target = "user39399393"
	err := messages.SendMessage(currentUsername, target, content)
	if err != nil {
		fmt.Printf("ERROR SENDING MESSAGE. Info: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}
