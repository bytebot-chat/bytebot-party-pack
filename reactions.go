package main

import (
	"strings"

	"github.com/bytebot-chat/gateway-discord/model"
)

func reactions(message model.Message) (string, bool) {
	// Passing message and not only content to allow
	// for some funky stuff later like redirections and all.

	reactionContent := ""
	switchToken := strings.Split(message.Message.Content, " ")[0]
	switch switchToken {
	case "!shrug":
		reactionContent = "¯\\_(ツ)_/¯"
	case "!lenny":
		reactionContent = "( ͡° ͜ʖ ͡°)"
	case "!tableflip":
		reactionContent = "(╯°□°)╯︵ ┻━┻"
	case "!tablefix":
		reactionContent = "┬─┬ノ( º _ ºノ)"
	case "!8ball":
		reactionContent = make8BallAnswer()
	case "!epeen":
		reactionContent = epeen(message.Author.Username)
	case "!ipinfo":
		reactionContent = ipinfo(message.Message.Content)
	case "!roll":
		reactionContent = diceTrigger(message.Author.Username, message.Message.Content)
	case "!choose":
		reactionContent = decisions(message.Author.Username, message.Message.Content)
	case "!dadjoke":
		reactionContent = jokeTrigger()

	}

	return reactionContent, reactionContent != ""
}
