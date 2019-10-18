package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")
	var server = flag.Bool("server", true, "if passed, the bot runs a websocket server, compatible with stdinToWebsocket")
	var botName = flag.String("name", "Unity "+UnityVersion, "The name for the bot in local games")
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
		ws.CreateServer(*addr)
	} else {
		log.Print("Running in local mode")
		conf := NewLocalConf()
		game := NewGame(*botName, conf)
		game.Loop()
	}
}
