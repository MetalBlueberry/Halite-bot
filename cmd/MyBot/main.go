package main

import (
	"fmt"
	"os"
	"time"

	"flag"

	"github.com/metalblueberry/halite-bot/pkg/hlt"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	var server = flag.Bool("server", false, "if passed, the bot runs a websocket server, compatible with stdinToWebsocket")
	var botName = flag.String("name", "MyBot", "The name for the bot in local games")
	var logToFile = flag.Bool("logToFile", false, "log to file, true if server is false")
	var debug = flag.Bool("debug", false, "prints to stdout debug information to be used with halite-debug project")
	flag.Parse()

	// TODO: Configure logrus
	if *logToFile {
		fname := fmt.Sprintf("logs_%s_gamelog.log", *botName)
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	if *server {
		log.Print("Running in server mode")
		ws := WebSocketHandler{
			Upgrader:  websocket.Upgrader{}, // use default options
			LogToFile: *logToFile,
			Debug:     *debug,
		}
		ws.CreateServer()
	} else {
		log.Print("Running in local mode")
		conf := NewLocalConf()
		game := NewGame(*botName, conf)
		game.Loop()
	}
}

type GameConfig struct {
	Source   <-chan string
	Response chan<- string
}

type Game struct {
	BotName string
	Conf    GameConfig
}

func NewGame(botName string, conf GameConfig) *Game {
	return &Game{
		BotName: botName,
		Conf:    conf,
	}
}

func NewConf(source <-chan string, response chan<- string) GameConfig {
	return GameConfig{
		Source:   source,
		Response: response,
	}
}

func (g Game) End() {
	close(g.Conf.Response)
}

// NewGame creates a new game with a name and communication channels.
func (g Game) Loop() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Game finished due to: ", r)
		}
	}()

	defer g.End()

	conn := hlt.NewConnection(g.BotName, g.Conf.Source, g.Conf.Response)

	log.Print("Game Starts")

	gameMap, _ := conn.UpdateMap()
	log.Println(gameMap.Grid.String())

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