package main

import (
	"math/rand"
	"strings"
)

func troll(nick, msg string) string {
	const TROLL_USAGE = "Usage: !troll <nick>"
	
	msg = strings.Trim(msg, "!troll ")
	if (len(msg) > 1 || len(msg) < 1) {
		return TROLL_USAGE
	}

	target := string(msg)
	numTrolls, dmg, dmgType := launchTrolls()

	switch dmg {
		case "":
			return "The troll launcher has malfunctioned."
		case "miss":
			return "Wha?! The trolls missed! That, like, never happens!"
	}

	return nick + " fires " + numTrolls + " at " + target + ", dealing " + dmg + " points of " + dmgType + " damage!"
}

func launchTrolls() (numTrolls, dmg, dmgType string) {
	damage_type := [13]string{"bludgeoning", "piercing", "slashing", "cold", "fire", "acid", "poison",
	"psychic", "necrotic", "radiant", "lightning", "thunder", "force"}

	trolls := rand.Intn(10)
	if trolls == 0 {
		return _, "", _
	}
	
	dmg = trollDamage(trolls)

	return string(trolls), dmg, damage_type[rand.Intn(12)]
}

func trollDamage (trolls int) string {
	i := 0
	trollDmg := 0
	for i < trolls {
		trollDmg += rand.Intn(20)
	}

	if trollDmg == 0 {
		return "miss"
	}

	return string(trollDmg)
}