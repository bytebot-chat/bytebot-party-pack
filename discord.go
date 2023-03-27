package main

import (
	"context"
	"sync"

	"github.com/bytebot-chat/gateway-discord/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func subscribeDiscord(ctx context.Context, wg *sync.WaitGroup, rdb *redis.Client, topic string, outbound []string) {
	// Unused, but I don't want to remove it from the function signature
	// this will be used when we have multiple inbound channels again
	defer wg.Done()

	log.Info().Msg("Subscribing to " + topic)
	sub := rdb.Subscribe(ctx, topic)
	log.Info().Msg("Subscribed!")
	channel := sub.Channel()

	// This will block until the context is cancelled
	for inboundMessage := range channel {
		msg := model.Message{}
		if err := msg.UnmarshalJSON([]byte(inboundMessage.Payload)); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal message")
			continue
		}

		log.Info().
			Str("channel", msg.Message.ChannelID).
			Str("user", msg.Message.Author.Username).
			Str("message", msg.Message.Content).
			Msg("Received message")

		answer, ok := reactions(msg)
		if ok {
			replyDiscord(ctx, msg, rdb, outbound[0], answer, true, false) // Reply but do not mention for now
		}
	}
}

// replyDiscord uses the discord-gateway's struct definitions and methods to construct a message
// and send it to the outbound channel
func replyDiscord(ctx context.Context, msg model.Message, rdb *redis.Client, topic string, answer string, shouldReply bool, shouldMention bool) {
	log.Info().
		Str("channel", msg.Message.ChannelID).
		Str("user", msg.Message.Author.Username).
		Str("message", msg.Message.Content).
		Msg("Sending message")

	// Construct the message
	outboundMessage := msg.RespondToChannelOrThread(msg.Metadata.Source, answer, shouldReply, shouldMention)
	stringMessage, err := outboundMessage.MarshalJSON()
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return
	}

	// Send the message
	if err := rdb.Publish(ctx, topic, stringMessage).Err(); err != nil {
		log.Error().Err(err).Msg("Failed to publish message")
	}
}
