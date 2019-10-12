package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/gorilla/websocket"
	halitedebug "github.com/metalblueberry/Halite-debug/pkg/client"

	"runtime/debug"

	"github.com/metalblueberry/halite-bot/pkg/hlt"
	log "github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	var server = flag.Bool("server", true, "if passed, the bot runs a websocket server, compatible with stdinToWebsocket")
	var botName = flag.String("name", "Unity", "The name for the bot in local games")
	var logToFile = flag.Bool("logToFile", false, "log to file, true if server is false")
	var debugf = flag.Bool("debug", true, "prints to stdout debug information to be used with halite-debug project")
	flag.Parse()

	// TODO: Configure logrus
	log.SetLevel(log.DebugLevel)
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
			Debug:     *debugf,
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
	Debug    bool
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
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))

		}
	}()

	defer g.End()

	conn := hlt.NewConnection(g.BotName, g.Conf.Source, g.Conf.Response)
	halitedebug.InitializeDefaultCanvas("http://localhost:8888", time.Now().Format("2006-01-02+15:04:05"), g.Conf.Debug)

	log.Print("Game Starts")

	gameMap, _ := conn.UpdateMap()
	commander := hlt.Commander{}
	commander.SetMap(&gameMap)

	gameturn := 1
	for {
		var start time.Time
		gameMap, start = conn.UpdateMap()
		commander.SetMap(&gameMap)
		commandQueue := []string{}

		myPlayer := commander.Players[gameMap.MyID]
		myShips := myPlayer.Ships

		PrintDebugEntities(gameMap)

		for i := 0; i < len(myShips); i++ {
			shipStart := time.Now()
			ship := myShips[i]
			if ship.DockingStatus == hlt.UNDOCKED {
				commandQueue = append(commandQueue, hlt.AstarStrategy(ship, commander))
			}
			log.Printf("Time for ship %s, total %s", time.Since(shipStart), time.Since(start))
		}

		log.Printf("Turn %v\n", gameturn)
		log.Printf("out %v\n", commandQueue)
		conn.SubmitCommands(commandQueue)
		halitedebug.Send(gameturn)
		gameturn++
	}
}

func PrintDebugEntities(gameMap hlt.Map) {
	for _, p := range gameMap.Planets {
		halitedebug.Circle(p.Entity, []string{"planet", fmt.Sprintf("player%d", int(p.Owned)*(1+p.Owner()))}...)
	}
	for _, ship := range gameMap.Ships {
		halitedebug.Circle(ship.Entity, []string{"ship", fmt.Sprintf("player%d", 1+ship.Owner())}...)
	}
	//for _, player := range gameMap.Players {
	//for _, ship := range player.Ships {
	//halitedebug.Circle(ship.Entity, []string{"ship", fmt.Sprintf("player%d", 1+ship.Owner())}...)
	//}
	//}

}
