package server

type Message struct {
	SocketId uint
	Message  string
}

func NewMessage(socketId uint, text string) Message {
	return Message{socketId, text}
}

func (message *Message) FromMessage(text string) Message {
	return NewMessage(message.SocketId, text)
}
