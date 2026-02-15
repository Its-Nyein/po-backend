package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"po-backend/configs"
	"po-backend/helper"
	"po-backend/utilities"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 30 * time.Second
	maxMsgSize = 1024
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Token          string `json:"token"`
	Event          string `json:"event"`
	ConversationID uint   `json:"conversationId,omitempty"`
	RecipientID    uint   `json:"recipientId,omitempty"`
}

type WSTypingEvent struct {
	Event          string `json:"event"`
	ConversationID uint   `json:"conversationId"`
	UserID         uint   `json:"userId"`
}

func writePump(client *configs.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := client.Conn.Close(); err != nil {
			log.Printf("Error closing websocket in writePump: %v", err)
		}
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				// Channel closed
				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := client.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func readPump(client *configs.Client) {
	defer func() {
		configs.WSManager.UnregisterClient(client.UserID, client)
	}()

	client.Conn.SetReadLimit(maxMsgSize)
	if err := client.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}
	client.Conn.SetPongHandler(func(string) error {
		return client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	authenticated := false

	for {
		_, rawMsg, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			continue
		}

		// Handle authentication
		if msg.Token != "" && !authenticated {
			claims, err := helper.ParseToken(msg.Token)
			if err != nil {
				log.Println("WS: Invalid token")
				continue
			}
			client.UserID = claims.UserID
			configs.WSManager.AddClient(client.UserID, client)
			authenticated = true
			log.Printf("WS: User %d connected\n", client.UserID)
			continue
		}

		if !authenticated {
			continue
		}

		// Handle typing events
		if msg.Event == "typing" {
			log.Printf("WS: Typing event from user %d, convID=%d, recipientID=%d\n", client.UserID, msg.ConversationID, msg.RecipientID)
			if msg.ConversationID > 0 && msg.RecipientID > 0 {
				typingEvent := WSTypingEvent{
					Event:          "typing",
					ConversationID: msg.ConversationID,
					UserID:         client.UserID,
				}
				utilities.SendWebSocketEvent(msg.RecipientID, typingEvent)
				log.Printf("WS: Typing relayed to user %d\n", msg.RecipientID)
			}
		}
	}
}

func HandleWebSocket(c *echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &configs.Client{
		Conn: ws,
		Send: make(chan []byte, 256),
	}

	go writePump(client)
	readPump(client)

	return nil
}
