package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Socket struct {
	id           uint
	OutToServer  chan Message
	InFromServer chan Message
}

func NewSocket(s *Server, conn *websocket.Conn) {
	newSocket := &Socket{
		s.currentId,
		s.InFromSockets,
		make(chan Message),
	}

	s.sockets[s.currentId] = newSocket
	s.currentId += 1
	log.Printf("Socket %v connected\n", newSocket.id)
	go newSocket.ListenForMessages(s, conn)
}

// TODO: listen from server

func (socket *Socket) ListenForMessages(s *Server, conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Socket %v disconnected (%v) \n", socket.id, err)
			// socket.OutToServer <- NewMessage(socket.id, "!leave")
			delete(s.sockets, socket.id)
			break
		}
		if mt == websocket.TextMessage {
			socket.OutToServer <- NewMessage(socket.id, string(message))
		}
	}
}
