package main

import (
	"flag"
	"os"
)

// This file contains environnement variables parsing related methods,
// for configuration purpose.

// parseEnv parses configuration environnement variables.
func parseEnv() {
	if !isFlagSet("inbound") {
		*inbound = parseStringFromEnv("INBOUND", "discord-inbound")
	}

	if !isFlagSet("outbound") {
		*outbound = parseStringFromEnv("OUTBOUND", "discord-outbound")
	}

	if !isFlagSet("redis") {
		*redisAddr = parseStringFromEnv("REDIS", "localhost:6379")
	}
}

// Parses a string from an env variable and returns it.
func parseStringFromEnv(varName, defaultVal string) string {
	val, set := os.LookupEnv(varName)
	if set {
		return val
	}
	return defaultVal
}

// This is used to check if a flag was set
// Must be called after flag.Parse()
func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
