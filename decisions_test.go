package main

import (
	"fmt"
	"testing"
)

/// Running simple tests on the choose function from the decisions trigger
func TestDecisionsChoose(t *testing.T) {
	var tests = []struct {
		input           string
		possibleOutputs []string // This makes the check efficient
	}{
		{"!choose a, b", []string{"a", "b"}},
		{"!choose a, b, c", []string{"a", "b", "c"}},
		{"!choose ", []string{""}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s", tt.input)
		t.Run(testname, func(t *testing.T) {
			ans := choose(tt.input)
			if !stringInStringSlice(ans, tt.possibleOutputs) {
				t.Errorf("got %s, want one of %v",
					ans,
					tt.possibleOutputs)
			}
		})
	}
}

func stringInStringSlice(a string, b []string) bool {
	fmt.Printf("Iterating over %s\n", b)
	for _, v := range b {
		if a == v {
			return true
		}
		fmt.Printf("%s != %s\n", a, v)
	}

	return false
}
