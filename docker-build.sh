#! /bin/sh

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags "-s -w -extldflags '-static'" -o ./opt/bytebot
