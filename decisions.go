package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func decisions(nick, msg string) string {
	choice := choose(msg)

	// if choose returns null, return "Choose what?" and command usage
	if choice == "" {
		choice = "Choose what?"
		return nick + ": Choose what? Usage: !choose choice1 or choice2"
	}

	return fmt.Sprintf("%s: The powers that be have chosen %s", nick, choice) // returning the choice
}

// the actual choose function
func choose(msg string) string {
	msg = strings.TrimPrefix(msg, "!choose ") // remove !choose trigger from the msg string
	choices := strings.Split(msg, "or")       // split on " or "

	fmt.Println(choices)

	// if the array has less than 2 elements (choices), return null
	if len(choices) < 2 {
		return ""
	}

	r := choices[rand.Intn(len(choices))] // using rand to select a random choice
	return strings.TrimSpace(r)           // trim spaces on the edges
}
