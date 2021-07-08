package main

import (
	"fmt"
	"strings"

	"github.com/justinian/dice"
)

const DICE_USAGE = "Usage: [num dice]d[sides](+/-num) (opt: if fudging)"

func dice(nick, r string) string {
	result := roll(r) // instantiate result variable and store result of roll

	return nick + " rolled a " + result
}

// function invoking dice library
func roll(r) string {
	res, _, err := dice.Roll(r)

	// if there's an error, return the DICE_USAGE const and log the error
	if err != nil {
		return DICE_USAGE
	}

	return fmt.Sprintf("%v", rec.Int())
}
