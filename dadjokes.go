package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type dadJoke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status string `json:"status"`
}

func (j *dadJoke) Unmarshal(b []byte) error {
	if err := json.Unmarshal(b, j); err != nil {
		return err
	}
	return nil
}

func jokeTrigger() string {
	j, jerr := getJoke()
	if jerr != nil {
		return "Joke machine broke."
	}

	return j.Joke
}

func getJoke() (dadJoke, error) {
	// instantiate Zerolog sublogger
	sublogger := log.With().
		Str("trigger", "dadjoke").
		Logger()

	joke := new(dadJoke)

	// define dadjoke url
	url := "https://icanhazdadjoke.com/"

	req, err := http.NewRequest("GET", url, nil) // crafting HTTP request with url and userAgent
	if err != nil {                              // if err, loggging error via Zerolog and returning error to chat
		sublogger.Warn().Err(err).Msg("HTTP request improperly formatted.")
		return *joke, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)

	if err != nil {
		sublogger.Warn().Err(err).Msg("Error getting HTTP request from server.")
		return *joke, err
	}

	defer resp.Body.Close() // closing the request body as required by net/http

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		err = errors.New(fmt.Sprintf("Non-OK HTTP status: %d", resp.StatusCode))
		sublogger.Warn().Err(err)
		return *joke, err
	}

	r, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		sublogger.Warn().Err(err).Msg("Error reading joke response body.")
		return *joke, err
	}

	err = joke.Unmarshal(r)
	if err != nil {
		sublogger.Warn().Err(err).Msg("Error unmarshaling the joke response body.")
		return *joke, err
	}

	return *joke, err
}
