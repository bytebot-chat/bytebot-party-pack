package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/bytebot-chat/gateway-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
)

func reply(ctx context.Context, m model.Message, rdb *redis.Client, reply string) {
	if !strings.HasPrefix(m.To, "#") { // DMs go back to source, channel goes back to channel
		m.To = m.From
	}
	m.From = ""
	m.Metadata.Dest = m.Metadata.Source
	m.Metadata.Source = "hello-world"
	m.Content = reply
	m.Metadata.ID = uuid.Must(uuid.NewV4(), *new(error))
	stringMsg, _ := json.Marshal(m)
	rdb.Publish(ctx, *outbound, stringMsg)
	return
}
