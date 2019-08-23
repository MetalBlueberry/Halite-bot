package main

import (
	"log"
	"runtime/debug"
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
	response := make(chan string)
	defer close(source)
	go func() {
		for {
			data, ok := <-response
			if !ok {
				log.Print("Game is finished, so closing websocket")
				c.Close()
				return
			}
			log.Println("ok")
			log.Printf("send: %s\n", data)
			c.WriteMessage(websocket.TextMessage, []byte(data))
		}
	}()

	go NewGame(source, response)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		source <- string(message)
		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func (ws *WebSocketHandler) CreateServer() {
	log.Print("Waiting for games")
	http.Handle("/echo", ws)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// golang starter kit with logging and basic pathfinding
// Arjun Viswanathan 2017 / github arjunvis

func main() {
	ws := WebSocketHandler{
		Upgrader: websocket.Upgrader{}, // use default options
	}
	ws.CreateServer()
}

func NewGame(source <-chan string, response chan<- string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
		}
	}()

	defer close(response)
	log.Print("Game Starts")

	// logging := true
	botName := "GoBot"

	conn := hlt.NewConnection(botName, source, response)

	// set up logging
	// if logging {
	// 	fname := "logs_" + strconv.Itoa(conn.PlayerTag) + "_gamelog.log"
	// 	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// 	if err != nil {
	// 		fmt.Printf("error opening file: %v\n", err)
	// 	}
	// 	defer f.Close()
	// 	log.SetOutput(f)
	// }

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
