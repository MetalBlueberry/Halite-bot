#!/bin/sh

rm logs/* 
#go build MyBot.go
halite -t -d "240 160" "go run ./cmd/stdinToWebsocket/main.go" "docker run --rm -i unity:v0.1.0" && \
find -type f -name "replay*" | grep -v "save" | sort | tail -n 1 | xargs -I{} chlorine -o {} && \
rm replay* && \
rm *.log
