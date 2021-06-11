package main

import (
	"hash/crc64"
	"math/rand"
	"strings"
	"time"
)

func epeen(nick string) string {
	peepeeSize := 20
	peepeeCrc := crc64.Checksum([]byte(nick+time.Now().Format("2006-01-02")), crc64.MakeTable(crc64.ECMA))
	peepeeRnd := rand.New(rand.NewSource(int64(peepeeCrc)))
	peepee := "8" + strings.Repeat("=", peepeeRnd.Intn(peepeeSize)) + "D"
	return nick + "'s peepee: " + peepee
}
