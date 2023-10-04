package model

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content interface{}
}
