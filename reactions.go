package main

import (
	"github.com/bytebot-chat/gateway-irc/model"
)

func reactions(message model.Message) (string, bool) {
	// Passing message and not only content to allow
	// for some funky stuff later like redirections and all.
	reactionContent := ""
	switch message.Content {
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
		reactionContent = epeen(message.From)
	}

	return reactionContent, reactionContent != ""
}
