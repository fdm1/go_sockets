package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Socket struct {
	id  uint
	out chan Message
}

func NewSocket(id uint, out chan Message) *Socket {
	return &Socket{
		id,
		out,
	}
}

func (socket *Socket) ListenForMessages(conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error reading message:", err)
			break
		}
		if mt == websocket.TextMessage {
			socket.out <- NewMessage(socket.id, string(message))
		}
	}
}
