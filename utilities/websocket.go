package utilities

import (
	"log"
	"po-backend/configs"
)

type WSResponse struct {
	Event string `json:"event"`
}

func SendWebSocketMessage(userID uint, event string) {
	conn, ok := configs.WSManager.GetClient(userID)
	if !ok {
		return
	}

	msg := WSResponse{Event: event}
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("WS: Error sending message to user %d: %v\n", userID, err)
		configs.WSManager.RemoveClient(userID)
	}
}
