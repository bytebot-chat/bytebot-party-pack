package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytebot-chat/gateway-discord/model"
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

func (c *ChannelContext) AddMessage(msg model.Message) error {
	ctx := context.Background()
	key := fmt.Sprintf("channel_context:%s", msg.ChannelID)
	message := fmt.Sprintf("%s: %s", msg.Author.Username, msg.Content)
	err := c.redisClient.LPush(ctx, key, message).Err()
	if err != nil {
		return err
	}
	return c.redisClient.Expire(ctx, key, 15*time.Minute).Err()
}

// GetContext returns all messages in the channel context
func (c *ChannelContext) GetContext(channelID string) ([]string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("channel_context:%s", channelID)
	return c.redisClient.LRange(ctx, key, 0, -1).Result()
}

// FormatContext formats the channel context into a string separated by newlines for sending to OpenAI
func (c *ChannelContext) FormatContext(channelID string) (string, error) {
	ctx, err := c.GetContext(channelID)
	if err != nil {
		return "", err
	}
	return strings.Join(ctx, "\n"), nil
}
