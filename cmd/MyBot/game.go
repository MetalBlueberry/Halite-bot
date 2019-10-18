package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	halitedebug "github.com/metalblueberry/Halite-debug/pkg/client"
	"github.com/metalblueberry/halite-bot/pkg/control"
	"github.com/metalblueberry/halite-bot/pkg/hlt"
	log "github.com/sirupsen/logrus"
)

type GameConfig struct {
	Source   <-chan string
	Response chan<- string
	Debug    bool
}

func NewConf(source <-chan string, response chan<- string) GameConfig {
	return GameConfig{
		Source:   source,
		Response: response,
	}
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
	commander := control.NewCommander()

	gameturn := 1
	commander.SetMap(gameMap, gameturn)
	for {
		var start time.Time
		gameMap, start = conn.UpdateMap()
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1900)

		PrintDebugEntities(gameMap)

		commander.SetMap(gameMap, gameturn)
		commander.Command(ctx)
		cancel()

		//commandQueue := []string{}

		//myPlayer := commander.Players()[gameMap.MyID]
		//myShips := myPlayer.Ships

		//for i := 0; i < len(myShips); i++ {
		////shipStart := time.Now()
		//ship := myShips[i]
		//if ship.DockingStatus == hlt.UNDOCKED {
		//commandQueue = append(commandQueue, control.AstarStrategy(ship, *commander))
		//}
		////log.Printf("Time for ship %s, total %s", time.Since(shipStart), time.Since(start))
		//}

		commandQueue := commander.CommandQueue()
		log.Printf("Turn time %s, avg per ship %f", time.Since(start), time.Since(start).Seconds()/float64(len(commander.Pilots)))
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
}
