package hlt

import (
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Connection performs all of the IO operations required to communicate
// game state and player movements with the Halite engine
type Connection struct {
	width, height int
	PlayerTag     int
	reader        <-chan string
	writer        chan<- string
}

func (c *Connection) sendString(input string) {
	c.writer <- input
}

func (c *Connection) getString() string {
	data, ok := <-c.reader
	if !ok {
		panic("game finished")
	}
	return strings.TrimSpace(data)
}

func (c *Connection) getInt() int {
	i, err := strconv.Atoi(c.getString())
	if err != nil {
		log.Printf("Errored on initial tag: %v", err)
	}
	return i
}

// NewConnection initializes a new connection for one of the bots
// participating in a match
func NewConnection(botName string, source <-chan string, response chan<- string) Connection {
	conn := Connection{
		reader: source,
		writer: response,
	}
	conn.PlayerTag = conn.getInt()
	sizeInfo := strings.Split(conn.getString(), " ")
	width, _ := strconv.Atoi(sizeInfo[0])
	height, _ := strconv.Atoi(sizeInfo[1])
	conn.width = width
	conn.height = height
	conn.sendString(botName)
	return conn
}

// UpdateMap decodes the current turn's game state from a string
func (c *Connection) UpdateMap() (Map, time.Time) {
	log.Printf("--- NEW TURN --- \n")
	gameString := c.getString()
	turnStart := time.Now()
	gameMap := ParseGameString(c, gameString)
	log.Printf("    Parsed map in %s", time.Since(turnStart))
	return gameMap, turnStart
}

// SubmitCommands encodes the player's commands into a string
func (c *Connection) SubmitCommands(commandQueue []string) {
	commandString := strings.Join(commandQueue, " ")
	log.Printf("Final string : %+v\n", commandString)
	c.sendString(commandString)
}
