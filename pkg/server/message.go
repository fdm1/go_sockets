package server

type Message struct {
	socketId uint
	message  string
}

func NewMessage(socketId uint, message string) Message {
	return Message{socketId, message}
}
