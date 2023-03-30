package main

import (
	"github.com/bytebot-chat/gateway-discord/model"
)

func echoHandler(m model.Message) *model.MessageSend {
	// Send the message back to the channel it came from

	reply := m.RespondToChannelOrThread("party-pack", m.Content, true, true)

	return reply
}
