package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	redisCtx  = context.Background()
	inbound   = flag.String("inbound", "discord-inbound", "inbound channel (from discord)")
	outbound  = flag.String("outbound", "discord-outbound", "outbound channel (to discord)")
	redisAddr = flag.String("redis", "", "redis address")
)

func init() {
	// Configure zerolog to log to stdout in JSON
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

}

func main() {

	flag.Parse()
	*redisAddr = os.Getenv("REDIS_URL")
	log.Info().
		Str("inbound", *inbound).
		Str("outbound", *outbound).
		Str("redis", *redisAddr).
		Msg("Starting Party-Pack")

	// Connect to redis
	rdb := redisConnect(*redisAddr, "", "", 0, redisCtx) // Not really concerned about the username and password here because we are using a URL provided by Fly.io
	if rdb == nil {
		log.Fatal().Msg("Unable to connect to redis")
		return
	}

	// Configure and run health checks
	go healthCheck(redisCtx, rdb)

	// Start the bot
	log.Info().Msg("Listening for messages on " + *inbound)
	log.Info().Msg("Sending messages to " + *outbound)

	// Subscribe to the discord channel
	var wg sync.WaitGroup // We don't need the waitgroup here, but I don't want to remove it from the function signature
	// Subscribe to the inbound channel
	// This will block until the context is cancelled
	subscribeDiscord(redisCtx, &wg, rdb, *inbound, []string{*outbound})
}
