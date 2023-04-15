package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sashabaranov/go-openai"
)

type Babbler struct {
	redisClient    *redis.Client
	openaiClient   *openai.Client
	channelContext *ChannelContext
}

func NewBabbler(redisClient *redis.Client) *Babbler {

	channelContext := NewChannelContext(redisClient)

	return &Babbler{
		redisClient:    redisClient,
		openaiClient:   openaiClient,
		channelContext: channelContext,
	}
}

func (b *Babbler) HandleWeatherCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Implement the HandleWeatherCommand logic
}

func (b *Babbler) HandleAskCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Implement the HandleAskCommand logic
}

func isMentioned(botUsername, message string) bool {
	return strings.Contains(message, "@"+botUsername)
}
