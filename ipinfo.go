package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type geoData struct {
	IP       string
	City     string
	Region   string
	Country  string
	Loc      string
	Org      string
	Postal   string
	Timezone string
}

func ipinfo(arg string) string {
	arg = strings.Split(arg, " ")[1]
	IP := net.ParseIP(arg)

	if IP == nil { // what is provided isn't an IP
		ips, err := net.LookupIP(arg)
		if err != nil {
			return "IP or domain not found"
		}

		IP = ips[0]
	}

	return ipLookup(IP.String())
}

func ipLookup(ip string) string {
	sublogger := log.With().
		Str("trigger", "ipinfo").
		Logger()

	sublogger.Info().Str("IP", ip).Msg("Looking up an IP")

	HTTPClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://ipinfo.io/%s", ip), nil) //nolint:noctx
	res, err := HTTPClient.Do(req)

	if err != nil {
		sublogger.Warn().Err(err).Msg("Couldn't query ipinfo.io")
	}

	jsonData, err := ioutil.ReadAll(res.Body)

	if err != nil {
		sublogger.Warn().Err(err).Msg("Couldn't read ipinfo.io answer")
	}

	err = res.Body.Close()
	if err != nil {
		sublogger.Warn().Err(err).Msg("IPInfo trigger, request's body was closed improperly")
	}
	return decodeJSON(jsonData, sublogger)
}

func decodeJSON(encodedJSON []byte, logger zerolog.Logger) string {
	var (
		ipinfo geoData
		reply  string
	)

	err := json.Unmarshal(encodedJSON, &ipinfo)
	if err != nil {
		logger.Warn().Err(err).Msg("IPinfo trigger, couldn't decode JSON")
	}

	logger.Debug().RawJSON("message", encodedJSON)

	if ipinfo.IP == "" {
		reply = "We are being rate limited, try again later or use ipinfo.io yourself."
	} else {
		reply = fmt.Sprintf("\u000312\u001f%s\u000f (%s): in %s, %s, %s (\u000312%s\u000f) postal code: %s, TZ: %s",
			ipinfo.IP,
			ipinfo.Org,
			ipinfo.City,
			ipinfo.Region,
			ipinfo.Country,
			ipinfo.Loc,
			ipinfo.Postal,
			ipinfo.Timezone)
	}
	return reply
}
