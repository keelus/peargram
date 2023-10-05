package ws

import (
	"peargram/models"
	"peargram/server/websocketing/client"
	"peargram/server/websocketing/model"
)

func MemberJoin(clients client.WebSocketClientsPool, cl client.WebSocketClient) {
	broadcast(clients, cl, model.WebSocketMessage{
		Type: MEMBER_JOIN,
		Content: map[string]string{
			"id": cl.Id(),
		},
	})
}

func MemberLeave(clients client.WebSocketClientsPool, cl client.WebSocketClient) {
	broadcast(clients, cl, model.WebSocketMessage{
		Type: MEMBER_LEAVE,
		Content: map[string]string{
			"id": cl.Id(),
		},
	})
}

func NewChatMessage(clients client.WebSocketClientsPool, cl client.WebSocketClient, message models.Message) {
	broadcast(clients, cl, model.WebSocketMessage{
		Type:    MESSAGE,
		Content: message,
	})
}
func NewMessage(clients client.WebSocketClientsPool, cl client.WebSocketClient, message string) {
	broadcast(clients, cl, model.WebSocketMessage{
		Type: MESSAGE,
		Content: map[string]string{
			"id":      cl.Id(),
			"message": message,
		},
	})
}

func broadcast(clients client.WebSocketClientsPool, cl client.WebSocketClient, msg model.WebSocketMessage) {
	// var clientList []string
	// for username, _ := range clients {
	// 	clientList = append(clientList, username)
	// }
	// for _, username := range clientList {
	// 	clients[username].Write(msg)
	// }
	cl.Write(msg)
	// array.ForEach(
	// 	array.Except(clients, func(item client.WebSocketClient) bool { return item.Id() == cl.Id() }),
	// 	func(item client.WebSocketClient) { item.Write(msg) },
	// )
}
