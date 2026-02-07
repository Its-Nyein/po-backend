package configs

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	Clients map[uint]*websocket.Conn
	Mu      sync.RWMutex
}

var WSManager = &WebSocketManager{
	Clients: make(map[uint]*websocket.Conn),
}

func (m *WebSocketManager) AddClient(userID uint, conn *websocket.Conn) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Clients[userID] = conn
}

func (m *WebSocketManager) RemoveClient(userID uint) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if conn, ok := m.Clients[userID]; ok {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing websocket: %v", err)
		}
		delete(m.Clients, userID)
	}
}

func (m *WebSocketManager) GetClient(userID uint) (*websocket.Conn, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	conn, ok := m.Clients[userID]
	return conn, ok
}
