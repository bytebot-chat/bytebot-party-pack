package main

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// lambda is an interface that defines a function that takes a string and a
// discordgo.Message, meant to be used as a lambda function that is called on
// every message that comes in from the gateway
//
// Replies to messages should be sent back to the gateway via Redis pubsub
type lambda interface {
	handle(string, discordgo.Message)
}

// messageRouter is a struct that contains a slice of lambdas
//
// It implements the lambda interface and can be used as a lambda function that
// calls all of the lambdas in its slice. It is intended to be used as the primary
// entrypoint for ensuring all lambdas are handled on every message.

// Registering new lambda is done by calling the registerLambda method on the
// messageRouter struct
type messageRouter struct {
	lambdas []func(string, discordgo.Message)
}

// handle takes a string and a discordgo.Message and iterates over the slice of
// lambdas, calling each lambda as a goroutine in turn until all lambdas have
// been called
func (r messageRouter) handle(topic, message string) {
	discordMessage, err := unmarshalDiscordMessage(message)
	if err != nil {
		log.Error().
			Err(err).
			Str("topic", topic).
			Str("message", message).
			Msg("Error unmarshalling discord message")
		return
	}

	for _, lambda := range r.lambdas {

		go lambda(topic, discordMessage)
	}
}

// registerLambda takes one or more lambdas and adds them to the slice of
// lambdas in the messageRouter struct
func (r *messageRouter) registerLambda(lambdas ...func(string, discordgo.Message)) {
	r.lambdas = append(r.lambdas, lambdas...)
}

// unmsrshalDiscordMessage takes a msg.Payload from redis that is expected to
// be a JSON representation of a Discord message and returns a discord.Message
func unmarshalDiscordMessage(message string) (discordgo.Message, error) {
	var m discordgo.Message
	err := json.Unmarshal([]byte(message), &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// marshalDiscordMessage takes a discord.Message and returns a JSON
// representation of it
func marshalDiscordMessage(message discordgo.Message) (string, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// routeInboundMessage takes a discord.Message and iterates over a slice of
// handlers, calling each handler in turn until all handlers have been called
// Handlers are expected to return a discord.Message and an error, which is
// then passed to the sending function that returns the message to be sent back
// to the gateway
func routeInboundMessage(m discordgo.Message, handlers []func(discordgo.Message)) {
	for _, handler := range handlers {
		go handler(m)
	}
}
