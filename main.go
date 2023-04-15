package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/bytebot-chat/gateway-discord/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	topic = "discord-inbound"
)

func main() {
	log.Info().Msg("Hello, world!")
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	rdb := redisConnect(os.Getenv("REDIS_URL"), context.Background())

	initOpenAIClient()

	// Create a new pubsub client
	log.Info().
		Str("topic", topic).
		Msg("Subscribing to topic")
	pubsub := rdb.Subscribe(context.Background(), topic)
	defer pubsub.Close()

	// Create a channel to receive messages
	log.Info().
		Str("topic", topic).
		Msg("Creating channel")
	ch := pubsub.Channel()

	// Loop forever
	log.Info().
		Str("topic", topic).
		Msg("Starting loop")
	go func() {
		for msg := range ch {
			log.Debug().
				Str("topic", topic).
				Str("msg", msg.Payload).
				Msg("Received message")

			// Unmarshal the message into a model.Message struct
			var m model.Message
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				log.Err(err).
					Str("func", "main").
					Str("msg", msg.Payload).
					Msg("Unable to unmarshal message")
				continue
			}

			// Route the message to the appropriate handler
			messageRouter(rdb, m)
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
