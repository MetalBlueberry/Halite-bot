FROM golang:1.13 as dev
WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY cmd/ cmd/
COPY pkg/ pkg/

RUN go test ./...

WORKDIR cmd/MyBot

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go install

FROM alpine:3.9 
COPY --from=dev /go/bin/MyBot /bin/bot
RUN chmod +x /bin/bot
ENTRYPOINT [ "bot" ]
CMD ["--logToFile=true", "--server=false"] 
