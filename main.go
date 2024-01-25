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

	// Create a new message router and register lambdas
	router := messageRouter{}
	router.registerLambda(messageLogger)

	// Pass the message router to a goroutine that will listen for messages
	go func(router messageRouter) {
		for {
			// Read message from channel
			msg := <-ch
			//log.Debug().Msgf("Received message: %s", msg.Payload)

			// Call our message router with the topic and message
			router.handle(msg.Channel, msg.Payload)
		}
	}(router)

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
