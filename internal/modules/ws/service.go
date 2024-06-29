package ws

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
	"github.com/kume1a/sonifybackend/internal/shared"
)

type wsPayloadDTO struct {
	MessageType string      `json:"messageType"`
	Payload     interface{} `json:"payload"`
}

type SendWSPayloadInput struct {
	ToSocketId  string
	MessageType string
	Payload     interface{}
}

func SendWSPayload(input SendWSPayloadInput) error {
	ws, exists := GetManager().GetConnection(input.ToSocketId)
	if !exists {
		return errors.New(shared.ErrSocketNotFound)
	}

	fullPayload := wsPayloadDTO{
		MessageType: input.MessageType,
		Payload:     input.Payload,
	}

	payloadJSON, err := json.Marshal(fullPayload)
	if err != nil {
		return errors.New(shared.ErrInvalidJSON)
	}

	if err := ws.WriteMessage(websocket.TextMessage, payloadJSON); err != nil {
		log.Println("Error writing to websocket:", err)
		return errors.New(shared.ErrInternal)
	}

	return nil
}
