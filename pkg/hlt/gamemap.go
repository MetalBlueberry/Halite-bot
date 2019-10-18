package hlt

import (
	"strconv"
	"strings"

	"github.com/metalblueberry/halite-bot/pkg/twoD"
)

// Map describes the current state of the game
type Map struct {
	MyID, Width, Height int
	Planets             []Planet
	Players             []Player
	Ships               map[int]Ship
	Entities            []Entitier
}

// Player has an ID for establishing ownership, and a number of ships
type Player struct {
	ID    int
	Ships []Ship
}

// ParsePlayer from a slice of game state tokens
func ParsePlayer(tokens []string) (Player, []string) {
	playerID, _ := strconv.Atoi(tokens[0])
	playerNumShips, _ := strconv.ParseFloat(tokens[1], 64)

	player := Player{
		ID:    playerID,
		Ships: []Ship{},
	}

	tokens = tokens[2:]
	for i := 0; float64(i) < playerNumShips; i++ {
		ship, tokensnew := ParseShip(playerID, tokens)
		tokens = tokensnew
		player.Ships = append(player.Ships, ship)
	}

	return player, tokens
}

// ParseGameString from a slice of game state tokens
func ParseGameString(c *Connection, gameString string) Map {
	tokens := strings.Split(gameString, " ")
	numPlayers, _ := strconv.Atoi(tokens[0])
	tokens = tokens[1:]

	gameMap := Map{
		MyID:     c.PlayerTag,
		Width:    c.width,
		Height:   c.height,
		Planets:  nil,
		Players:  make([]Player, numPlayers),
		Ships:    make(map[int]Ship),
		Entities: make([]Entitier, 0),
	}

	for i := 0; i < numPlayers; i++ {
		player, tokensnew := ParsePlayer(tokens)
		tokens = tokensnew
		gameMap.Players[player.ID] = player
		for j := 0; j < len(player.Ships); j++ {
			ship := player.Ships[j]
			gameMap.Entities = append(gameMap.Entities, ship.Entity)
			gameMap.Ships[ship.id] = ship
		}
	}

	numPlanets, _ := strconv.Atoi(tokens[0])
	gameMap.Planets = make([]Planet, 0, numPlanets)
	tokens = tokens[1:]

	for i := 0; i < numPlanets; i++ {
		planet, tokensnew := ParsePlanet(tokens)
		tokens = tokensnew
		gameMap.Planets = append(gameMap.Planets, planet)
		gameMap.Entities = append(gameMap.Entities, planet.Entity)
	}

	return gameMap
}

type byX []Entity

func (a byX) Len() int           { return len(a) }
func (a byX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byX) Less(i, j int) bool { return a[i].x < a[j].x }

type byY []Entity

func (a byY) Len() int           { return len(a) }
func (a byY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byY) Less(i, j int) bool { return a[i].y < a[j].y }

func (gameMap Map) ObstaclesBetween(start twoD.Circler, end twoD.Positioner, ignoreIDs ...int) (bool, Entitier) {
	return ObstaclesBetween(start, end, gameMap.Entities, ignoreIDs...)
}

func ObstaclesBetween(start twoD.Circler, end twoD.Positioner, obstacles []Entitier, ignoreIDs ...int) (bool, Entitier) {
	_, _, r1 := start.Circle()
	//_, _ := end.Position()
	StartToEnd := twoD.Distance(start, end)
	for i := 0; i < len(obstacles); i++ {
		entity := obstacles[i]

		if contains(entity.ID(), ignoreIDs) {
			continue
		}

		_, _, r := entity.Circle()
		margin := r1 + r
		dist := twoD.DistancePointToLine(start, end, entity)
		if dist > margin {
			continue
		}

		projection := twoD.Project(start, end, entity)
		relative := projection / StartToEnd
		if relative > 0 && relative < 1 {
			return true, entity
		}

		endDist := twoD.Distance(end, entity)
		if endDist < margin {
			return true, entity
		}
	}
	return false, nil
}
func contains(entity int, ignoreIDs []int) bool {
	for _, id := range ignoreIDs {
		if id == entity {
			return true
		}
	}
	return false
}

// NearestPlanetsByDistance orders all planets based on their proximity
// to a given ship from nearest for farthest
//func NearestPlanetsByDistance(ship Ship, planets []Planet) []Planet {
//for i := 0; i < len(planets); i++ {

//planets[i].Distance = ship.CalculateDistanceTo(planets[i].Entity)
//}

//sort.Sort(byDist(planets))

//return planets
//}

//type byDist []Planet

//func (a byDist) Len() int           { return len(a) }
//func (a byDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a byDist) Less(i, j int) bool { return a[i].Distance < a[j].Distance }
