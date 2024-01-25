package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*

Notes about development:
- All functions that use an topic and message should have the topic passed as
  immediately before the message. This makes it easier to reason about
  not having to remember which is which.

*/

func main() {
	var (
		topic = "discord.inbound.*"
	)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msg("Hello, world!")

	rdb := redisConnect(os.Getenv("REDIS_URL"), context.Background())

	// Create a new pubsub client that listens to the topic from BYTEBOT_TOPIC or "discord.inbound.*"
	log.Info().Msg("Creating new pubsub client")
	if os.Getenv("BYTEBOT_TOPIC") != "" {
		topic = os.Getenv("BYTEBOT_TOPIC")
	} else {
		log.Info().Msg("BYTEBOT_TOPIC not set, using default topic: " + topic)
	}

	pubsub := rdb.PSubscribe(context.Background(), topic)
	defer pubsub.Close()
	log.Info().Msgf("Subscribed to topic: %s", topic)

	ch := pubsub.Channel()

	// Set a shared context for all lambdas
	ctx := context.Background()

	var lambdas = []lambda{
		pingPongLambda,
	}

	// Pass the message router to a goroutine that will listen for messages
	go func() {
		for {
			// Read message from channel
			msg := <-ch

			/*
				log.Debug().
					Str("func", "main").
					Str("channel", msg.Channel).
					Str("payload", msg.Payload).
					Msg("Received message")
			*/

			// Unmarshal the message into a discordgo.Message
			dgoMessage, err := unmarshalDiscordMessage(msg.Payload)
			if err != nil {
				log.Error().Err(err).Msg("Error unmarshalling message")
			}

			// Convert the topic string into a pubsubDiscordTopicAddr struct
			topicAddr, err := newPubsubDiscordTopicAddr(msg.Channel)
			if err != nil {
				log.Error().
					Err(err).
					Str("topic", msg.Channel).
					Msg("Error converting topic string to pubsubDiscordTopicAddr struct")
			}

			for _, lambda := range lambdas {
				go lambda(ctx, rdb, topicAddr, dgoMessage)
			}
		}
	}()

	// Add a healthcheck endpoint on port 8080
	log.Info().Msg("Registering healthcheck endpoint")
	http.Handle("/health", health.NewHandler(newHealthChecker()))
	log.Info().Msg("Starting http server on port 8080")
	http.ListenAndServe(":8080", nil)

}

func newHealthChecker() health.Checker {
	return health.NewChecker(

		health.WithCacheDuration(1*time.Second),

		health.WithTimeout(10*time.Second),

		health.WithCheck(
			health.Check{
				Name:    "redis",
				Timeout: 2 * time.Second,
				Check: func(ctx context.Context) error {
					log.Info().Msg("Running redis check")
					return nil
				},
			},
		),

		// Set a status listener that will be invoked when the health status changes.
		// More powerful hooks are also available (see docs).
		health.WithStatusListener(func(ctx context.Context, state health.CheckerState) {
			log.Info().Msgf("Health status changed: %s", state.Status)
		}),
	)
}
