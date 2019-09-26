package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/metalblueberry/halite-bot/pkg/hlt"
	log "github.com/sirupsen/logrus"
)

type CanvasServer struct {
	Url string
}

type Action map[string]interface{}

//Method string
//Params []interface{}
//}

func NewCanvasServer() *CanvasServer {
	return &CanvasServer{
		Url: "http://127.0.0.1:8888/test/%d",
	}
}
func (c CanvasServer) Entity(turn int, p hlt.Entity, style []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	// initialize http client
	client := &http.Client{}

	//colors := map[int]string{
	//0: "",
	//1: "player1",
	//2: "player2",
	//3: "player3",
	//4: "player4",
	//}
	log.
		WithField("Owner", p.Owner).
		//WithField("Owner", p.Owned).
		Debug("Draw Debug Planet")

	actions := []Action{
		Action{
			"Method": "Entity",
			"X":      p.X,
			"Y":      p.Y,
			"R":      p.Radius,
			"Class":  style,
		},
	}

	// marshal User to json
	json, err := json.Marshal(actions)
	if err != nil {
		panic(err)
	}

	// set the HTTP method, url, and request body
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf(c.Url, turn), bytes.NewBuffer(json))
	if err != nil {
		log.Printf("error building request for debug server %v", err)
		return
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Unable to reach debug draw server %v", err)
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	//log.Println(resp.StatusCode)
	//log.Println(body)
}
