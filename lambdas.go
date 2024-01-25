package main

import (
	"context"
	"fmt"
	"hash/crc64"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func pingPongLambda(ctx context.Context, rdb *redis.Client, topic *pubsubDiscordTopicAddr, m discordgo.Message) {
	if m.Content == "ping" {
		log.Info().Msg("Pong")

		err := publishDiscordMessage(ctx, rdb, topic.getReplyTopic(), "pong")
		if err != nil {
			log.Error().
				Err(err).
				Str("topic", topic.getReplyTopic()).
				Msg("Error publishing message")
		}
	}
}

func weatherLambda(ctx context.Context, rdb *redis.Client, topic *pubsubDiscordTopicAddr, m discordgo.Message) {
	var content string
	if strings.HasPrefix(m.Content, "!weather") {
		city := strings.TrimSpace(strings.TrimPrefix(m.Content, "!weather"))
		weather, err := getWeather(city)
		if err != nil {
			log.Warn().Err(err).Str("city", city).Msg("Error fetching weather data")
			content = fmt.Sprintf("Error fetching weather data for %s: %v", city, err)
		} else {
			content = weather
		}
	}

	err := publishDiscordMessage(ctx, rdb, topic.getReplyTopic(), content)
	if err != nil {
		log.Error().
			Str("func", "weatherLambda").
			Err(err).
			Str("topic", topic.getReplyTopic()).
			Msg("Error publishing message")
	}
}

func reactionsLambda(ctx context.Context, rdb *redis.Client, topic *pubsubDiscordTopicAddr, m discordgo.Message) {
	var content string
	switch m.Content {
	case "ping":
		content = "pong"
	case "pong":
		content = "ping"
	case "!shrug":
		content = "¯\\_(ツ)_/¯"
	case "!lenny":
		content = "( ͡° ͜ʖ ͡°)"
	case "!tableflip":
		content = "(╯°□°）╯︵ ┻━┻"
	case "!tablefix":
		content = "┬─┬ノ( º _ ºノ)"
	case "!unflip":
		content = "┬─┬ノ( º _ ºノ)"
	case "hello", "hi", "hey", "howdy", "sup", "yo", "hiya", "anyong", "bonjour", "salut", "hallo", "moin":
		content = sayHello(m.Author.Username)
	case "!epeen":
		content = epeen(m)
	}

	err := publishDiscordMessage(ctx, rdb, topic.getReplyTopic(), content)
	if err != nil {
		log.Error().
			Str("func", "reactionsLambda").
			Err(err).
			Str("topic", topic.getReplyTopic()).
			Msg("Error publishing message")
	}
}

func sayHello(username string) string {
	greeting := []string{
		"hi",
		"hello",
		"hey",
		"howdy",
		"sup",
		"yo",
		"hiya",
		"anyong",
		"bonjour",
		"salut",
		"hallo",
		"moin",
	}
	return fmt.Sprintf("%s %s", greeting[rand.Intn(len(greeting))], username)
}

func epeen(m discordgo.Message) string {
	peepeeSize := 20
	peepeeCrc := crc64.Checksum([]byte(m.Author.Username+time.Now().Format("2006-01-02")), crc64.MakeTable(crc64.ECMA))
	peepeeRnd := rand.New(rand.NewSource(int64(peepeeCrc)))
	return "8" + strings.Repeat("=", peepeeRnd.Intn(peepeeSize)) + "D"
}

func getWeather(city string) (string, error) {
	apiURL := fmt.Sprintf("https://wttr.in/%s?format=2", url.QueryEscape(city)) // format=2 is a preconfigured format for a single line of weather data. 3 and 4 include the city name with ugly + signs as delimiters
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	weather := string(body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("weather API returned an error: %s", weather)
	}

	return weather, nil
}
