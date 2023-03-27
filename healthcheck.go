package main

import (
	"context"
	"net/http"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func healthCheck(ctx context.Context, rdb *redis.Client) {
	// Setup a health check endpoint
	checker := health.NewChecker(
		health.WithCacheDuration(1*time.Second),
		health.WithTimeout(10*time.Second),
		// Check the redis connection with a ping
		health.WithPeriodicCheck(
			15*time.Second,
			3*time.Second,
			health.Check{
				Name: "redis",
				Check: func(ctx context.Context) error {
					_, err := rdb.Ping(ctx).Result()
					return err
				},
			}),

		// Test the redis pubsub connection by subscribing and unsubscribing
		health.WithPeriodicCheck(
			15*time.Second,
			3*time.Second,
			health.Check{
				Name: "redis-pubsub",
				Check: func(ctx context.Context) error {
					pubsub := rdb.Subscribe(ctx, "test")
					_, err := pubsub.Receive(ctx)
					if err != nil {
						return err
					}
					err = pubsub.Close()
					return err
				},
			}),
	)

	// Register the health check endpoint
	http.Handle("/health", health.NewHandler(checker))

	// Start the health check server
	log.Fatal().Err(http.ListenAndServe(":8080", nil)).Msg("Error starting health check server")
}
