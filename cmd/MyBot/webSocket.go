package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WebSocketHandler struct {
	Upgrader  websocket.Upgrader
	LogToFile bool
	Debug     bool
}

func (ws *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer socket.Close()

	source := make(chan string)
	response := make(chan string)

	defer close(source)

	conf := NewConf(source, response)
	conf.Debug = ws.Debug
	game := NewGame("WSBot", conf)

	go ws.ListenForGameUpdates(response, socket)
	go ws.ForwardMessages(source, socket)

	game.Loop()
}

func (ws *WebSocketHandler) CreateServer(addr string) {
	log.Print("Waiting for games")
	if ws.LogToFile {
		log.Print("The the game log will be saved to a file")
	}
	http.Handle("/echo", ws)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (ws *WebSocketHandler) ListenForGameUpdates(response <-chan string, socket *websocket.Conn) {
	for {
		data, ok := <-response
		if !ok {
			log.Print("Game is finished, so closing websocket")
			socket.Close()
			return
		}
		//log.Printf("send: %s\n", data)
		err := socket.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			log.Panicf("message could not be sent over websocket %s", err)
		}
	}
}

func (ws *WebSocketHandler) ForwardMessages(source chan<- string, socket *websocket.Conn) {
	for {
		_, message, err := socket.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		//log.Printf("recv: %s", message)
		source <- string(message)
	}
}
