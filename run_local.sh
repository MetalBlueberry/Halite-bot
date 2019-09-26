#!/bin/sh

rm logs/* 
#go build MyBot.go
halite -t -d "240 160" "go run ./cmd/MyBot/*.go --server=FALSE --logToFile=TRUE" "/home/victor/Documents/Projects/Halite-II/airesources/Python3/MyBot.py" && \
find -type f -name "replay*" | grep -v "save" | sort | tail -n 1 | xargs -I{} chlorine -o {} && \
rm replay* && \
rm *.log
