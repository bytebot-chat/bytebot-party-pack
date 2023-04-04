package main

import (
	"context"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

// redisConnect is used to manage the connection to redis and gracefully exit if the connection fails
// If the address is a URL, then it will be parsed and the redis.Options struct will be populated
// There is no need to pass the username and password if the address is a URL. Use "" for both instead.
func redisConnect(addr string, ctx context.Context) *redis.Client {

	var (
		redisOpts redis.Options
	)

	// This is deployed on fly, so we automatically get a redis URL
	redisOpts = redisParseURL(addr)

	log.Info().
		Str("func", "redisConnect").
		Str("addr", redisOpts.Addr).
		Str("user", redisOpts.Username).
		Str("pass", redisOpts.Password).
		Int("db", redisOpts.DB).
		Msg("Connecting to redis")

	rdb := redis.NewClient(&redisOpts)

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Err(err).
			Str("func", "redisConnect").
			Msg("Unable to connect to redis. Exiting!")
		os.Exit(1)
	}

	return rdb
}

// returns a redis.Options struct for use with redis.NewClient
func redisParseURL(url string) redis.Options {

	// Remove the redis:// prefix
	url = strings.TrimPrefix(url, "redis://")

	// Split the URL along the @ symbol
	urlParts := strings.Split(url, "@")

	auth := urlParts[0]
	addr := urlParts[1]

	// Split the auth along the : symbol to get the username and password
	authParts := strings.Split(auth, ":")
	username := authParts[0]
	password := authParts[1]

	// Split the addr along the / symbol to get the address and database
	// If there is no /, then the database is 0
	addrParts := strings.Split(addr, "/")
	addr = addrParts[0]
	db := 0

	// If there is a : in the address, then there is a port
	// If there is no : in the address, then there is no port
	if strings.Contains(addr, ":") {
		addrParts = strings.Split(addr, ":")
		addr = addrParts[0]
		port := addrParts[1]
		addr = addr + ":" + port
	} else {
		addr = addr + ":6379"
	}

	// Return the redis.Options struct
	return redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		Username: username,
	}
}
