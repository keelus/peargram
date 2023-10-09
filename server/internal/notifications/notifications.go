package notifications

import (
	"fmt"
	"peargram/database"
	"peargram/models"
	"time"
)

func SendNotification(notificationType int, actor string, target string, postID int) {
	DB := database.ConnectDB()

	date := time.Now().Unix()

	// TODO IN DB: CHECK INT FILTER
	if notificationType == models.NOTIFICATION_FOLLOW || notificationType == models.NOTIFICATION_FOLLOW_REQUEST {
		_, err := DB.Query("INSERT INTO notifications (actor, target, type, date) VALUES (?, ?, ?, ?)", actor, target, notificationType, date)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		_, err := DB.Query("INSERT INTO notifications (actor, target, type, date, post) VALUES (?, ?, ?, ?, ?)", actor, target, notificationType, date, postID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Notification registered")
}

func GetNotifications(username string) []models.Notification {
	var notifications []models.Notification

	DB := database.ConnectDB()
	err := DB.Select(&notifications, "SELECT * FROM notifications WHERE target=? ORDER BY date DESC;", username)
	if err != nil {
		fmt.Println(err)
	}

	return notifications
}
