package main

import (
	"hash/crc64"
	"math/rand"
	"strings"
	"time"

	"github.com/bytebot-chat/gateway-discord/model"
)

func simpleHandler(m model.Message) *model.MessageSend {
	// Send the message back to the channel it came from

	var (
		app           string = "party-pack"
		content       string
		shouldReply   bool
		shouldMention bool
	)

	switch m.Content {
	case "ping":
		content = "pong"
	case "pong":
		content = "ping"
	case "!shrug":
		content = "¯\\_(ツ)_/¯"
	case "!lenny":
		content = "( ͡° ͜ʖ ͡°)"
	case "!tableflip":
		content = "(╯°□°）╯︵ ┻━┻"
	case "!tablefix":
		content = "┬─┬ノ( º _ ºノ)"
	case "!unflip":
		content = "┬─┬ノ( º _ ºノ)"
	case "!epeen":
		content = epeen(m)
	}

	return m.RespondToChannelOrThread(app, content, shouldReply, shouldMention)
}

func epeen(m model.Message) string {
	peepeeSize := 20
	peepeeCrc := crc64.Checksum([]byte(m.Author.Username+time.Now().Format("2006-01-02")), crc64.MakeTable(crc64.ECMA))
	peepeeRnd := rand.New(rand.NewSource(int64(peepeeCrc)))
	return "8" + strings.Repeat("=", peepeeRnd.Intn(peepeeSize)) + "D"
}
