package ws

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func SendWSPayload(toSocketId string, payload interface{}) error {
	ws, exists := GetManager().GetConnection(toSocketId)
	if !exists {
		return errors.New(shared.ErrSocketNotFound)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return errors.New(shared.ErrInvalidJSON)
	}

	if err := ws.WriteMessage(websocket.TextMessage, payloadJSON); err != nil {
		log.Println("Error writing to websocket:", err)
		return errors.New(shared.ErrInternal)
	}

	return nil
}
