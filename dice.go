package main

import (
	"fmt"

	"github.com/justinian/dice"
)

func diceTrigger(nick, r string) string {
	result := rollDice(r) // instantiate result variable and store result of roll

	return nick + " rolled a " + result
}

// function invoking dice library
func rollDice(r string) string {
	const DICE_USAGE = "Usage: [num dice]d[sides](+/-num) (opt: if fudging)"

	res, _, err := dice.Roll(r)

	// if there's an error, return the DICE_USAGE const and log the error
	if err != nil {
		return DICE_USAGE
	}

	return fmt.Sprintf("%v", res.Int())
}
