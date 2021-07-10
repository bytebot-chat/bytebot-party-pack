package main

import "testing"

func TestDiceTrigger(t *testing.T) {
	const DICE_USAGE = "Usage: [num dice]d[sides](+/-num) (opt: if fudging)"

	roll := diceTrigger("parsec", "1d20")
	if roll == DICE_USAGE {
		t.Errorf("Got DICE_USAGE const, expected int.")
	}
}
