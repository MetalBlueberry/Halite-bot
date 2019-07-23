package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/metalblueberry/halite-bot/pkg/hlt"
)

// golang starter kit with logging and basic pathfinding
// Arjun Viswanathan 2017 / github arjunvis

func main() {

	
	logging := true
	botName := "GoBot"

	conn := hlt.NewConnection(botName)

	// set up logging
	if logging {
		fname := "logs/" + strconv.Itoa(conn.PlayerTag) + "_gamelog.log"
		f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("error opening file: %v\n", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	gameMap,_ := conn.UpdateMap()
	log.Println(gameMap.Grid.String())

	gameturn := 1
	for {
		var start time.Time
		gameMap,start = conn.UpdateMap()
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
