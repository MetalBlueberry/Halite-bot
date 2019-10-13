#!/bin/bash
docker build -t unity .

cd cmd/MyBot
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build

