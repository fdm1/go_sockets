package main

import (
	"log"

	"github.com/fdm1/go_sockets/pkg/client"
)

func main() {
	log.SetFlags(0)

	client.RunClient()
}
