package notificationsApi

import (
	"net/http"
	"peargram/notifications"

	"github.com/gin-gonic/gin"
)

func GETNotifications(c *gin.Context) {
	username := c.Param("username")
	userNotis := notifications.GetNotifications(username)
	c.JSON(http.StatusOK, userNotis)
}
