package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/bytebot-chat/gateway-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
)

func replyIRC(ctx context.Context, m Message, rdb *redis.Client, topic, reply string) {
	returnMsg := &model.Message{}

	if !strings.HasPrefix(m.To, "#") { // DMs go back to source, channel goes back to channel
		returnMsg.To = m.From
	}
	returnMsg.From = ""
	returnMsg.Metadata.Dest = m.Metadata.Source
	returnMsg.Metadata.Source = "party-pack"
	returnMsg.Content = reply
	returnMsg.Metadata.ID = uuid.Must(uuid.NewV4(), *new(error))
	stringReply, _ := json.Marshal(returnMsg)
	rdb.Publish(ctx, topic, stringReply)

	log.Debug().
		RawJSON("message", stringReply).
		Msg("Reply")

	return
}

func replyDiscord(ctx context.Context, m Message, rdb *redis.Client, topic, reply string) {
	metadata := Metadata{
		Dest:   m.Metadata.Source,
		Source: "party-pack",
		ID:     uuid.Must(uuid.NewV4(), *new(error)),
	}

	returnMsg := &Message{
		From:      "",
		ChannelID: m.From,
		Metadata:  metadata,
	}

	stringReply, _ := json.Marshal(returnMsg)
	rdb.Publish(ctx, topic, stringReply)
	log.Debug().
		RawJSON("message", stringReply).
		Msg("Reply")

	return
}

type stringArrayFlags []string

func (i *stringArrayFlags) String() string {
	return "String array flag"
}

func (i *stringArrayFlags) Set(s string) error {
	*i = append(*i, s)
	return nil
}

type Message struct {
	From      string
	To        string
	Source    string
	Content   string
	ChannelID string
	Raw       interface{}
	Metadata  Metadata
}

type Metadata struct {
	Source string
	Dest   string
	ID     uuid.UUID
}

func (m *Message) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, m); err != nil {
		return err
	}
	return nil
}
