package utilities

import (
	"encoding/json"
	"log"
	"po-backend/configs"
)

type WSResponse struct {
	Event string `json:"event"`
}

func safeSend(client *configs.Client, data []byte) (sent bool) {
	defer func() {
		if r := recover(); r != nil {
			sent = false
		}
	}()
	select {
	case client.Send <- data:
		return true
	default:
		return false
	}
}

func SendWebSocketMessage(userID uint, event string) {
	client, ok := configs.WSManager.GetClient(userID)
	if !ok {
		return
	}

	msg := WSResponse{Event: event}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("WS: Error marshaling message: %v\n", err)
		return
	}

	if !safeSend(client, data) {
		log.Printf("WS: Send channel full or closed for user %d, dropping message\n", userID)
	}
}

func SendWebSocketEvent(userID uint, payload interface{}) {
	client, ok := configs.WSManager.GetClient(userID)
	if !ok {
		log.Printf("WS: SendWebSocketEvent - no client found for user %d\n", userID)
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("WS: Error marshaling event: %v\n", err)
		return
	}

	if !safeSend(client, data) {
		log.Printf("WS: Send channel full or closed for user %d, dropping event\n", userID)
	}
}
