package controllers

import (
	"log"
	"net/http"

	"po-backend/configs"
	"po-backend/helper"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Token string `json:"token"`
	Event string `json:"event"`
}

func HandleWebSocket(c *echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := ws.Close(); err != nil {
			log.Printf("Error closing websocket: %v", err)
		}
	}()

	var userID uint

	for {
		var msg WSMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			if userID != 0 {
				configs.WSManager.RemoveClient(userID)
			}
			break
		}

		if msg.Token != "" && userID == 0 {
			claims, err := helper.ParseToken(msg.Token)
			if err != nil {
				log.Println("WS: Invalid token")
				continue
			}
			userID = claims.UserID
			configs.WSManager.AddClient(userID, ws)
			log.Printf("WS: User %d connected\n", userID)
		}
	}

	return nil
}
