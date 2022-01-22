package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO: close all sockets on interrupt

type Server struct {
	sockets   map[uint]*Socket
	currentId uint
	in        chan Message
	out       chan Message
}

var upgrader = websocket.Upgrader{} // use default options

func NewServer() *Server {
	server := &Server{
		map[uint]*Socket{},
		0,
		make(chan Message),
		make(chan Message),
	}

	go func() {
		for msg := range server.in {
			log.Printf("recv: %s", msg.message)
		}
	}()

	return server
}

func (server *Server) HandleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	socket := NewSocket(server.currentId, server.in)
	server.sockets[server.currentId] = socket
	server.currentId += 1
	socket.ListenForMessages(c)
}
