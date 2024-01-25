package main

import (
	"context"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
)

type lambda func(context.Context, *redis.Client, *pubsubDiscordTopicAddr, discordgo.Message)

// unmashalDiscordMessage is a function that takes a string and returns a
// discordgo.Message
func unmarshalDiscordMessage(msg string) (discordgo.Message, error) {
	var dgoMessage discordgo.Message
	err := json.Unmarshal([]byte(msg), &dgoMessage)
	return dgoMessage, err
}

// publishDiscordMessage is a function that takes a string and
// writes it to destination pubsub channel. It requires a context and a redis
// client.
func publishDiscordMessage(ctx context.Context, rdb *redis.Client, destination, content string) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	return rdb.Publish(ctx, destination, string(b)).Err()
}
