package messagesApi

import (
	"fmt"
	"net/http"
	"peargram/messages"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GETMessages(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	username := c.Query("username")
	offsetStr := c.Query("offset")

	if username == "" || offsetStr == "" {
		fmt.Println("Parameter missing!")
	}

	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		fmt.Println("Unexpected offset type")
	}

	fmt.Println(username)
	fmt.Println(offsetInt)
	fmt.Println(uint(offsetInt))

	chat := messages.GetChat(currentUsername, username, uint(offsetInt))

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

	err := messages.SendMessage(currentUsername, target, content)
	if err != nil {
		fmt.Printf("ERROR SENDING MESSAGE. Info: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}
