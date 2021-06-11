package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bbriggs/bytebot-irc/model"
	"github.com/go-redis/redis/v8"
	"github.com/satori/go.uuid"
)

var addr = flag.String("redis", "localhost:6379", "Redis server address")
var inbound = flag.String("inbound", "irc-inbound", "Pubsub queue to listen for new messages")
var outbound = flag.String("outbound", "irc", "Pubsub queue for sending messages outbound")

func main() {
	flag.Parse()
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: *addr,
		DB:   0,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		err := rdb.Ping(ctx).Err()
		if err != nil {
			panic(err)
		}
	}

	topic := rdb.Subscribe(ctx, *inbound)
	channel := topic.Channel()
	for msg := range channel {
		m := &model.Message{}
		err := m.Unmarshal([]byte(msg.Payload))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%+v\n", m)
		if m.Content == "!epeen" {
			reply(ctx, *m, rdb, epeen(m.From))
		}
	}
}

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
