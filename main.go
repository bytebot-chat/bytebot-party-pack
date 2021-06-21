package main

import (
	"context"
	"flag"
	"time"

	"github.com/bytebot-chat/gateway-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var addr = flag.String("redis", "localhost:6379", "Redis server address")
var inbound = flag.String("inbound", "irc-inbound", "Pubsub queue to listen for new messages")
var outbound = flag.String("outbound", "irc", "Pubsub queue for sending messages outbound")

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	flag.Parse()

	log.Info().
		Str("Redis address", *addr).
		Str("Inbound queue", *inbound).
		Str("Outbound queue", *outbound).
		Msg("Bytebot Party Pack starting up!")

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: *addr,
		DB:   0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		log.Warn().Msg("Ping timeout, trying to connect to redis again...")
		time.Sleep(3 * time.Second)
		err := rdb.Ping(ctx).Err()
		if err != nil {
			log.Fatal().Err(err).
				Msg("Couldn't connect to redis server")
		}
	}

	topic := rdb.Subscribe(ctx, *inbound)
	channel := topic.Channel()
	for msg := range channel {
		m := &model.Message{}
		err := m.Unmarshal([]byte(msg.Payload))
		if err != nil {
			log.Error().
				Str("message payload", msg.Payload).
				Err(err)
		}
		log.Debug().
			RawJSON("Received message", []byte(msg.Payload)).
			Msg("Received message")

		if m.Content == "!epeen" {
			reply(ctx, *m, rdb, epeen(m.From))
		} else {
			// Trigger doing it's own treatment of the message
			answer, activated := reactions(*m)
			if activated {
				reply(ctx, *m, rdb, answer)
			}
		}

	}
}
