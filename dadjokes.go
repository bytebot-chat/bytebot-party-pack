package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

func dadjoke() string {
	return httpRequest()
}

func httpRequest() string {

	// instantiate Zerolog sublogger
	sublogger := log.With().
		Str("trigger", "dadjoke").
		Logger()

	// define dadjoke url and userAgent
	url := "https://icanhazdadjoke.com/"
	userAgent := "bytebot-party-pack, an IRC microservice bot funpack thingy. https://github.com/bytebot-chat/bytebot-party-pack/"

	req, err := http.NewRequest("GET", url, nil) // crafting HTTP request with url and userAgent
	if err != nil {                              // if err, loggging error via Zerolog and returning error to chat
		sublogger.Warn().Err(err).Msg("HTTP request formatting error")
		return fmt.Sprintf("An error occured: %v", err)
	}

	req.Header.Set("User-Agent", userAgent)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sublogger.Warn().Err(err).Msg("HTTP request error")
		return fmt.Sprintf("An error occured: %v", err)
	}
	defer resp.Body.Close() // closing the request body as required by net/http

	if resp.StatusCode != 200 {
		statusCode := strconv.Itoa(resp.StatusCode)
		err = errors.New("resp.StatusCode: " + statusCode)
		sublogger.Warn().Err(err).Msg("HTTP status code")
		return "Error: HTTP Status Code " + statusCode
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		sublogger.Warn().Err(err).Msg("ioutil.ReadAll error")
		return fmt.Sprintf("An error occured: %v", err)
	}
	return string(r)
}
