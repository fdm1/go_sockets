package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// TODO: close all sockets on interrupt

type Server struct {
	sockets       map[uint]*Socket
	currentId     uint
	InFromSockets chan Message
	OutToChat     chan Message
	InFromChat    chan Message
	OutToSockets  chan Message
}

var upgrader = websocket.Upgrader{} // use default options

func NewServer() *Server {
	server := &Server{
		map[uint]*Socket{},
		0,
		make(chan Message),
		make(chan Message),
		make(chan Message),
		make(chan Message),
	}

	go func() {
		for msg := range server.InFromSockets {
			log.Printf("server recv from sockets: %s", msg.Message)
			server.OutToChat <- msg.FromMessage(string(msg.Message))
		}
	}()

	// go func() {
	// 	for msg := range server.InFromChat {
	// 		log.Printf("server recv from chat: %s", msg.Message)
  //     // server.sockets[msg.SocketId].In
	// 		// server.OutToSockets <- msg.FromMessage(string(msg.Message))
	// 	}
	// }()

	return server
}

func (s *Server) HandleConnection(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	NewSocket(s, c)
}
