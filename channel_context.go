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

func (c *ChannelContext) AddMessage(channelID, message string) error {
	ctx := context.Background()
	key := fmt.Sprintf("channel_context:%s", channelID)
	err := c.redisClient.LPush(ctx, key, message).Err()
	if err != nil {
		return err
	}
	return c.redisClient.Expire(ctx, key, 15*time.Minute).Err()
}
