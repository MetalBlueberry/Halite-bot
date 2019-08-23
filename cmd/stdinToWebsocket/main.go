package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	// set up logging
	fname := "logs_" + strconv.Itoa(0) + "_fw.log"
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
	}
	defer f.Close()
	log.SetOutput(f)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Print("Sending response ", string(message))
			_, err = os.Stdout.Write(message)
			if err != nil {
				log.Panic(err)
			}
			_, err = os.Stdout.Write([]byte("\n"))
			if err != nil {
				log.Panic(err)
			}
		}
	}()

	stdin := make(chan []byte)
	scanner := bufio.NewScanner(os.Stdin)
	go func(scanner *bufio.Scanner) {
		defer close(done)

		for scanner.Scan() {
			stdin <- append(scanner.Bytes(), '\n')
		}

		if scanner.Err() != nil {
			panic(scanner.Err())
		}

	}(scanner)

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		case line := <-stdin:
			err := c.WriteMessage(websocket.TextMessage, line)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
