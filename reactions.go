package main

import (
	"strings"
)

func reactions(message Message) (string, bool) {
	// Passing message and not only content to allow
	// for some funky stuff later like redirections and all.

	reactionContent := ""
	switchToken := strings.Split(message.Content, " ")[0]

	if *modules == "all" || strings.Contains(*modules, switchToken[1:]) {
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
			reactionContent = epeen(message.From)
		case "!ipinfo":
			reactionContent = ipinfo(message.Content)
		case "!roll":
			reactionContent = diceTrigger(message.From, message.Content)
		case "!choose":
			reactionContent = decisions(message.From, message.Content)
		case "!dadjoke":
			reactionContent = jokeTrigger()
		}
	}

	return reactionContent, reactionContent != ""
}
