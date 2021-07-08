package main

import (
	"math/rand"
	"strings"
)

func choose(nick, msg string) string {

	choices := strings.Split(msg, "or") // declaring the choices array and splitting the msg content into parts

	choice := choices[rand.Intn(len(choices)-1)] // using rand to select a random choice

	return nick + ": The Powers that Be have chosen: " + choice // returning the choice
}
