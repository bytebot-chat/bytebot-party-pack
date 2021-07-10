package main

import (
	"testing"
)

func TestDadjokeTrigger(t *testing.T) {
	for i := 0; i < 10; i++ {
		joke := jokeTrigger()
		if joke == "Joke machine broke." {
			t.Errorf("Got %s, want dad joke", joke)
		}
	}
}
