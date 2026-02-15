package configs

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	UserID uint
}

type WebSocketManager struct {
	Clients map[uint]*Client
	Mu      sync.RWMutex
}

var WSManager = &WebSocketManager{
	Clients: make(map[uint]*Client),
}

func (m *WebSocketManager) AddClient(userID uint, client *Client) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if existing, ok := m.Clients[userID]; ok {
		close(existing.Send)
	}
	m.Clients[userID] = client
}

func (m *WebSocketManager) RemoveClient(userID uint) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if client, ok := m.Clients[userID]; ok {
		close(client.Send)
		if err := client.Conn.Close(); err != nil {
			log.Printf("Error closing websocket: %v", err)
		}
		delete(m.Clients, userID)
	}
}

func (m *WebSocketManager) UnregisterClient(userID uint, client *Client) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	current, ok := m.Clients[userID]
	if !ok || current != client {
		return
	}
	close(client.Send)
	if err := client.Conn.Close(); err != nil {
		log.Printf("Error closing websocket: %v", err)
	}
	delete(m.Clients, userID)
}

func (m *WebSocketManager) GetClient(userID uint) (*Client, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	client, ok := m.Clients[userID]
	return client, ok
}
