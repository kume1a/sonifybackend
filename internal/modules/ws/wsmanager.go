package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	manager = &WebSocketManager{
		connections: make(map[string]*websocket.Conn),
	}
	lock sync.RWMutex
)

type WebSocketManager struct {
	connections map[string]*websocket.Conn
}

func GetManager() *WebSocketManager {
	return manager
}

func (m *WebSocketManager) addConnection(key string, conn *websocket.Conn) {
	lock.Lock()
	defer lock.Unlock()
	m.connections[key] = conn
}

func (m *WebSocketManager) GetConnection(key string) (*websocket.Conn, bool) {
	lock.RLock()
	defer lock.RUnlock()
	conn, exists := m.connections[key]
	return conn, exists
}
