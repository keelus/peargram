package messages

import (
	"fmt"
	"peargram/database"
	"peargram/models"
	"peargram/server/internal/users"
	wsHandler "peargram/server/websocketing/handler"
	"peargram/server/websocketing/ws"
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
			ws.NewChatMessage(wsClients, wsClients[target], sendingMessage)
			return nil
		}
	}

	return nil

}

func GetChatMessageCount(username1 string, username2 string) uint {
	var messageCount int

	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT COUNT(*) FROM messages WHERE (target=? OR actor=?) AND (target=? OR actor=?);", username1, username1, username2, username2).Scan(&messageCount)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return uint(messageCount)
}

func GetChatPreview(username1 string, username2 string) models.Chat {
	var lastMessage models.Message

	DB := database.ConnectDB()
	err := DB.QueryRowx("SELECT * FROM messages WHERE (target=? OR actor=?) AND (target=? OR actor=?) ORDER BY date DESC, id DESC LIMIT 1", username1, username1, username2, username2).StructScan(&lastMessage)
	if err != nil {
		fmt.Println(err)
		return models.Chat{}
	}

	chat := models.Chat{Participants: [2]string{username1, username2}, Messages: []models.Message{lastMessage}}
	return chat
}

func GetChat(username1 string, username2 string, offset uint) models.Chat {
	var messages []models.Message
	messages = make([]models.Message, 0)

	DB := database.ConnectDB()
	err := DB.Select(&messages, "SELECT * FROM messages WHERE (target=? OR actor=?) AND (target=? OR actor=?) ORDER BY date DESC, id DESC LIMIT 30 OFFSET ?;", username1, username1, username2, username2, offset)
	if err != nil {
		fmt.Println(err)
		return models.Chat{}
	}

	messageCount := GetChatMessageCount(username1, username2)
	chat := models.Chat{Participants: [2]string{username1, username2}, Messages: messages, MessageTotalCount: messageCount}
	return chat
}

func GetChatPreviews(username string) []models.Chat {
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
		chat := GetChatPreview(username, user)
		chats = append(chats, chat)
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[i].Messages[0].Date > chats[j].Messages[0].Date
	})

	return chats
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
		chat := GetChat(username, user, 0)
		chats = append(chats, chat)
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[i].Messages[0].Date > chats[j].Messages[0].Date
	})

	return chats
}
