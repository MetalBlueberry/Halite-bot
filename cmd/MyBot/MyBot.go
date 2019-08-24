package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"flag"
	"net/http"

	"github.com/metalblueberry/halite-bot/pkg/hlt"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type WebSocketHandler struct {
	Upgrader websocket.Upgrader
}

func (ws *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	source := make(chan string)
	defer close(source)

	response := make(chan string)

	go NewGame("WSBot", false, source, response)

	go func() {
		for {
			data, ok := <-response
			if !ok {
				log.Print("Game is finished, so closing websocket")
				c.Close()
				return
			}
			log.Printf("send: %s\n", data)
			err := c.WriteMessage(websocket.TextMessage, []byte(data))
			if err != nil {
				log.Panicf("message could not be sent over websocket %s", err)
			}
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		source <- string(message)
	}
}

func (ws *WebSocketHandler) CreateServer() {
	log.Print("Waiting for games")
	http.Handle("/echo", ws)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func main() {
	var server = flag.Bool("server", false, "if passed, the bot runs a websocket server, compatible with stdinToWebsocket")
	var botName = flag.String("name", "MyBot", "The name for the bot in local games")
	flag.Parse()

	if *server {
		log.Print("Running in server mode")
		ws := WebSocketHandler{
			Upgrader: websocket.Upgrader{}, // use default options
		}
		ws.CreateServer()
	} else {
		NewLocalGame(*botName)
	}
}

// NewLocalGame wraps stdin and stdout to be compatible with NewGame function.
func NewLocalGame(botName string) {
	done := make(chan struct{})

	stdin := make(chan string)
	stdout := make(chan string)

	go func() {
		defer close(done)
		for {
			message, ok := <-stdout
			if !ok {
				log.Println("stdout closed")
				return
			}
			_, err := fmt.Fprintf(os.Stdout, "%s\n", message)
			if err != nil {
				log.Panic(err)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	go func(scanner *bufio.Scanner) {
		defer close(done)

		for scanner.Scan() {
			stdin <- scanner.Text()
		}

		if scanner.Err() != nil {
			panic(scanner.Err())
		}

	}(scanner)

	NewGame(botName, true, stdin, stdout)

}

// NewGame creates a new game with a name and communication channels.
func NewGame(botName string, logToFile bool, source <-chan string, response chan<- string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Game finished due to: ", r)
		}
	}()

	defer close(response)

	// logging := true

	conn := hlt.NewConnection(botName, source, response)

	if logToFile {
		fname := "logs_" + strconv.Itoa(conn.PlayerTag) + "_gamelog.log"
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v\n", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	log.Print("Game Starts")

	defer func() {
		if r := recover(); r != nil {
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	gameMap, _ := conn.UpdateMap()
	// log.Println(gameMap.Grid.String())

	gameturn := 1
	for {
		var start time.Time
		gameMap, start = conn.UpdateMap()
		commandQueue := []string{}

		myPlayer := gameMap.Players[gameMap.MyID]
		myShips := myPlayer.Ships

		for i := 0; i < len(myShips); i++ {
			shipStart := time.Now()
			ship := myShips[i]
			if ship.DockingStatus == hlt.UNDOCKED {
				commandQueue = append(commandQueue, hlt.AstarStrategy(ship, gameMap))
			}
			log.Printf("Time for ship %s, total %s", time.Since(shipStart), time.Since(start))
		}

		log.Printf("Turn %v\n", gameturn)
		conn.SubmitCommands(commandQueue)
		gameturn++
	}

}
