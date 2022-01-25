package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/fdm1/go_sockets/pkg/chat"
	"github.com/fdm1/go_sockets/pkg/server"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	log.SetFlags(0)

	flag.Parse()
	log.Printf("Starting server at %v", *addr)
	s := server.NewServer()
	chat.NewChat(s.OutToChat, s.InFromChat)
	http.HandleFunc("/", s.HandleConnection)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
