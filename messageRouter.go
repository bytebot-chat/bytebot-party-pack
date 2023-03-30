package main

import (
	"context"
	"encoding/json"

	"github.com/bytebot-chat/gateway-discord/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func messageRouter(rdb *redis.Client, m model.Message) {
	// Use a for loop to send the message to all handlers

	var handlers []func(model.Message) *model.MessageSend

	handlers = append(handlers, echoHandler)
	//handlers = append(handlers, helloHandler)

	for _, handler := range handlers {
		reply := handler(m)
		if reply != nil {
			send(rdb, *reply)
		}
	}
}

// this function is ugly, but it works. don't you dare touch it.
func send(rdb *redis.Client, reply model.MessageSend) {
	// Set the message metadata
	meta := model.Metadata{
		ID:     uuid.NewV4(),
		Source: "party-pack",
		Dest:   reply.Metadata.Source,
	}

	// Set the message metadata
	reply.Metadata = meta

	// Convert the message to JSON
	jsonReply, err := json.Marshal(reply)
	if err != nil {
		log.Err(err).
			Str("func", "send").
			Str("msg", reply.Metadata.ID.String()).
			Msg("Unable to marshal message")
		return
	}

	// Send the message to redis
	err = rdb.Publish(context.Background(), "discord-outbound", jsonReply).Err()
	if err != nil {
		log.Err(err).
			Str("func", "send").
			Str("msg", reply.Metadata.ID.String()).
			Msg("Unable to publish message")
		return
	}

	log.Info().
		Str("func", "send").
		Str("msg", reply.Metadata.ID.String()).
		Msg("Message sent")
}
