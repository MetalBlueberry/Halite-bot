#!/bin/sh

#go build MyBot.go
docker exec -u 1000 -it --workdir "/home/victor/Halite-bot/"  halite-bot_devcontainer_development_1 go build
halite -d "240 160" "./cmd/MyBot/MyBot" "./cmd/MyBot/MyBot"
find -type f -name "replay*" | sort | tail -n 1 | xargs -I{} chlorine -o {}