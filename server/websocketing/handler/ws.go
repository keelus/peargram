package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"peargram/server/websocketing/client"
	wshandler "peargram/server/websocketing/ws"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients client.WebSocketClientsPool
	m       sync.Mutex
)

func WebSocketConnect(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	go startClient(c, ws)
}

func startClient(c *gin.Context, ws *websocket.Conn) {
	if clients == nil {
		clients = make(client.WebSocketClientsPool)
	}

	session := sessions.Default(c)

	if session.Get("Username") == nil {
		fmt.Println("User is not logged in!")
		return
	}

	username := session.Get("Username").(string)

	fmt.Printf("NEW WEBSOCKET SESSION REQUEST | USER: %s \n", username)

	cl := client.NewWebSocketClient(ws)

	m.Lock()
	clients[username] = cl
	m.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Printf("error: %v", err)
		}

		m.Lock()
		defer m.Unlock()
		// ws session closed
		fmt.Printf("WEBSOCKET SESSION DISCONNECTED | USER : %s\n", username)
		delete(clients, username)
		cl.Close()
	}()

	ctx := context.Context(c)
	cl.Launch(ctx)
	wshandler.MemberJoin(clients, cl)

	for {
		select {
		case msg, ok := <-cl.Listen():
			if !ok {
				return
			} else {
				fmt.Println(clients)
				switch msg.Type {
				case wshandler.MESSAGE:
					// wshandler.NewMessage(clients, cl, msg.Content["message"])
				default:
					log.Printf("unknown message type: %s", msg.Type)
					return
				}
			}
		case err := <-cl.Error():
			log.Printf("web socket error: %v", err)
		case <-cl.Done():
			wshandler.MemberLeave(clients, cl)
			return
		}
	}
}

func GetClients() client.WebSocketClientsPool {
	return clients
}
