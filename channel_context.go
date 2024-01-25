package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type ChannelContext struct {
	redisClient *redis.Client
}

func NewChannelContext(redisClient *redis.Client) *ChannelContext {
	return &ChannelContext{
		redisClient: redisClient,
	}
}

// Keeps track of the last 100 messages in a channel or the last 15 minutes, whichever is shorter
// This is used to help provide context to the handlers if they need it
// It's particularly useful for ChatGPT because it's context sensitive
func (c *ChannelContext) AddMessage(channelID, message string) error {
	ctx := context.Background()
	key := fmt.Sprintf("channel_context:%s", channelID)
	err := c.redisClient.LPush(ctx, key, message).Err()
	if err != nil {
		return err
	}
	err = c.redisClient.LTrim(ctx, key, 0, 99).Err()
	if err != nil {
		return err
	}
	return c.redisClient.Expire(ctx, key, 15*time.Minute).Err()
}
