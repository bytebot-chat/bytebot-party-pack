package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func messageLogger(topic string, m discordgo.Message) {
	log.Debug().
		Str("topic", topic).
		Str("author", m.Author.Username).
		Str("content", m.Content).
		Msg("Received message")
}

/*
func simpleHandler(m model.Message) *model.MessageSend {
	// Send the message back to the channel it came from

	var (
		app           string = "party-pack"
		content       string
		shouldReply   bool
		shouldMention bool
	)

	if strings.HasPrefix(m.Content, "!weather") {
		city := strings.TrimSpace(strings.TrimPrefix(m.Content, "!weather"))
		weather, err := getWeather(city)
		if err != nil {
			log.Warn().Err(err).Str("city", city).Msg("Error fetching weather data")
			content = fmt.Sprintf("Error fetching weather data for %s: %v", city, err)
		} else {
			content = fmt.Sprintf("The weather in %s is %s", city, weather)
		}
		return m.RespondToChannelOrThread(app, content, true, false)
	}

	if strings.HasPrefix(strings.ToLower(m.Content), fmt.Sprintf("hey <@%s", os.Getenv("BOT_DISCORD_ID"))) {
		prefix := fmt.Sprintf("hey <@%s", os.Getenv("BOT_DISCORD_ID"))
		question := strings.TrimSpace(strings.TrimPrefix(m.Content, prefix))
		answer, err := handleAskCommand(question)
		if err != nil {
			content = "Error: " + err.Error()
		} else {
			content = answer
		}
		shouldReply = true
		return m.RespondToChannelOrThread(app, content, true, false)
	}

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
	case "hello":
		content = fmt.Sprintf("hey, %s", m.Author.Username)
		shouldReply = true
	case "!epeen":
		content = epeen(m)
		shouldReply = true
	}

	return m.RespondToChannelOrThread(app, content, shouldReply, shouldMention)
}

func epeen(m model.Message) string {
	peepeeSize := 20
	peepeeCrc := crc64.Checksum([]byte(m.Author.Username+time.Now().Format("2006-01-02")), crc64.MakeTable(crc64.ECMA))
	peepeeRnd := rand.New(rand.NewSource(int64(peepeeCrc)))
	return "8" + strings.Repeat("=", peepeeRnd.Intn(peepeeSize)) + "D"
}

func getWeather(city string) (string, error) {
	apiURL := fmt.Sprintf("https://wttr.in/%s?format=%%C|%%t|%%w", url.QueryEscape(city))
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
*/
