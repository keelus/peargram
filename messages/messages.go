package messages

import (
	"fmt"
	"peargram/database"
	"peargram/models"
	wsHandler "peargram/server/websocketing/handler"
	"peargram/server/websocketing/ws"
	"peargram/users"
	"sort"
	"time"
)

func SendMessage(actor string, target string, content string) error {
	DB := database.ConnectDB()

	date := time.Now().Unix()

	// TODO: Check user exists
	if !users.UserExists(target) {
		return fmt.Errorf("User does not exist.")
	}
	// TODO IN DB: CHECK INT FILTER
	_, err := DB.Query("INSERT INTO messages (actor, target, content, date) VALUES (?, ?, ?, ?)", actor, target, content, date)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Message sent")

	var sendingMessage models.Message
	sendingMessage = models.Message{
		ID:        99,
		Actor:     actor,
		Target:    target,
		Content:   content,
		Reactions: nil,
		Date:      0,
	}

	wsClients := wsHandler.GetClients()
	for username := range wsClients {
		if username == target {
			fmt.Println("Message via Websocket sent")
			ws.NewChatMessage(wsHandler.GetClients(), wsHandler.GetClients()[target], sendingMessage)
			return nil
		}
	}

	return nil

}

func GetChat(username1 string, username2 string) models.Chat {
	var messages []models.Message

	DB := database.ConnectDB()
	err := DB.Select(&messages, "SELECT * FROM messages WHERE (target=? OR actor=?) AND (target=? OR actor=?) ORDER BY date DESC;", username1, username1, username2, username2)
	if err != nil {
		fmt.Println(err)
		return models.Chat{}
	}

	chat := models.Chat{Participants: [2]string{username1, username2}, Messages: messages}
	return chat
}

func GetChats(username string) []models.Chat {
	var interactedWith []string
	var chats []models.Chat

	DB := database.ConnectDB()
	err := DB.Select(&interactedWith, "SELECT target as withUsers FROM messages WHERE actor=? UNION SELECT actor as withUsers FROM messages WHERE target=?", username, username)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Now we get the chats

	for _, user := range interactedWith {
		chat := GetChat(username, user)
		chats = append(chats, chat)
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[i].Messages[0].Date > chats[j].Messages[0].Date
	})

	return chats
}
