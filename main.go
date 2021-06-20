package main

import (
	"context"
	"flag"
	"fmt"
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

	log.Info().Msg("Bytebot Party Pack starting up!")
	log.Info().Msg("Bytebot Party Pack starting with the following configs:")
	log.Info().Msg("Redis address: " + *addr)
	log.Info().Msg("Inbound queue: " + *inbound)
	log.Info().Msg("Outbound queue: " + *outbound)

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: *addr,
		DB:   0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		err := rdb.Ping(ctx).Err()
		if err != nil {
			panic(err)
		}
	}

	topic := rdb.Subscribe(ctx, *inbound)
	channel := topic.Channel()
	for msg := range channel {
		log.Print(msg.Payload)
		m := &model.Message{}
		err := m.Unmarshal([]byte(msg.Payload))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", m)
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
