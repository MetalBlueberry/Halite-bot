#!/bin/bash
git tag $UNITY_VERSION
docker build -t unity:$UNITY_VERSION .

cd cmd/MyBot
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build

