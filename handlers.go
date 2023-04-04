package main

import (
	"github.com/bytebot-chat/gateway-discord/model"
)

func echoHandler(m model.Message) *model.MessageSend {
	// Send the message back to the channel it came from

	reply := m.RespondToChannelOrThread("party-pack", m.Content, true, true)

	switch m.Content {
	case "ping":
		reply.Content = "pong"
	case "pong":
		reply.Content = "ping"
	default:
		reply = nil
	}

	return reply
}
