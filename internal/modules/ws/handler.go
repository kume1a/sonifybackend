package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kume1a/sonifybackend/internal/shared"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func HandleWsUpgrade(w http.ResponseWriter, r *http.Request) {
	authPayload, err := shared.GetAuthPayload(r)
	if err != nil {
		log.Println("Error getting auth payload: ", err)
		shared.ResUnauthorized(w, shared.ErrUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		return
	}

	GetManager().addConnection(authPayload.UserID.String(), conn)
}
