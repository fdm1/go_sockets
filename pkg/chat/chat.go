package chat

import (
	"fmt"
	"log"
	"strings"

	"github.com/fdm1/go_sockets/pkg/server"
)

type Chat struct {
	channels         map[string]map[uint]struct{}
	socketChannelMap map[uint]string
	in               chan server.Message
	out              chan server.Message
}

func NewChat(in chan server.Message, out chan server.Message) {
	newChat := Chat{
		map[string]map[uint]struct{}{},
		map[uint]string{},
		in,
		out,
	}

	go func() {
		for msg := range newChat.in {
			log.Printf("chat recv: %s", msg.Message)
			newChat.ProcessMessage(msg)
		}
	}()
}

func (c *Chat) LeaveChannel(id uint) {
	if channel, ok := c.socketChannelMap[id]; ok {
		delete(c.socketChannelMap, id)
		delete(c.channels[channel], id)
		if len(c.channels[channel]) == 0 {
			delete(c.channels, channel)
		}
	}
}

func (c *Chat) JoinChannel(id uint, channel string) {
	c.LeaveChannel(id)
	_, ok := c.channels[channel]
	if !ok {
		c.channels[channel] = map[uint]struct{}{}
	}
	c.channels[channel][id] = struct{}{}
	c.socketChannelMap[id] = channel
}

func (c *Chat) ProcessMessage(msg server.Message) {
	outgoingMessage := ""
	socketIds := map[uint]struct{}{}
	socketIds[msg.SocketId] = struct{}{}
	if strings.HasPrefix(msg.Message, "!join ") {
		splitString := strings.Split(msg.Message, " ")
		if len(splitString) > 1 {
			c.JoinChannel(msg.SocketId, splitString[1])
			outgoingMessage = fmt.Sprintf("You have joined %v\n", splitString[1])
		}
	} else if strings.HasPrefix(msg.Message, "!leave") {
		c.LeaveChannel(msg.SocketId)
		outgoingMessage = fmt.Sprintf("You have left your channel\n")
	} else {
		if channel, ok := c.socketChannelMap[msg.SocketId]; ok {
			outgoingMessage = msg.Message
			socketIds = c.channels[channel]
		}
		if outgoingMessage == "" {
			outgoingMessage = "You must join a channel first via ':join <channel>'"
		}
	}
	c.SendMessages(socketIds, outgoingMessage)
}

func (c *Chat) SendMessages(socketIds map[uint]struct{}, message string) {
	log.Printf("Sending message from chat to server: '%v' to sockets: %v\n", message, socketIds)
	for socketId := range socketIds {
		log.Printf("sending to %v\n", socketId)
		c.out <- server.NewMessage(socketId, message)
		log.Printf("sent to %v\n", socketId)
	}
	log.Print("done sending from chat")
}
