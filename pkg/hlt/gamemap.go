package hlt

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/metalblueberry/halite-bot/pkg/navigation"
)

// Map describes the current state of the game
type Map struct {
	MyID, Width, Height int
	Planets             []Planet
	Players             []Player
	Entities            []Entity
	Grid                *navigation.Grid
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
		Entities: make([]Entity, 0),
		Grid:     navigation.NewGrid(c.width, c.height),
	}

	for i := 0; i < numPlayers; i++ {
		player, tokensnew := ParsePlayer(tokens)
		tokens = tokensnew
		gameMap.Players[player.ID] = player
		for j := 0; j < len(player.Ships); j++ {
			ship := player.Ships[j].Entity
			gameMap.Entities = append(gameMap.Entities, player.Ships[j].Entity)
			if player.ID == gameMap.MyID {
				gameMap.Grid.PaintShip(ship.X, ship.Y, 0)
			} else {
				gameMap.Grid.PaintShip(ship.X, ship.Y, 5)
			}
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
		gameMap.Grid.PaintPlanet(planet.Entity.X, planet.Entity.Y, planet.Entity.Radius)
	}

	return gameMap
}

type byX []Entity

func (a byX) Len() int           { return len(a) }
func (a byX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byX) Less(i, j int) bool { return a[i].X < a[j].X }

type byY []Entity

func (a byY) Len() int           { return len(a) }
func (a byY) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byY) Less(i, j int) bool { return a[i].Y < a[j].Y }

// ObstaclesBetween demonstrates how the player might determine if the path
// between two enitities is clear
func (gameMap Map) ObstaclesBetween(start Entity, end Entity) (bool, Entity) {
	x1 := start.X
	y1 := start.Y
	x2 := end.X
	y2 := end.Y
	dx := x2 - x1
	dy := y2 - y1
	a := dx*dx + dy*dy + 1e-8
	crossterms := x1*x1 - x1*x2 + y1*y1 - y1*y2

	for i := 0; i < len(gameMap.Entities); i++ {
		entity := gameMap.Entities[i]
		if entity.ID == start.ID || entity.ID == end.ID {
			continue
		}

		x0 := entity.X
		y0 := entity.Y

		closestDistance := end.CalculateDistanceTo(entity)
		if closestDistance < entity.Radius+1 {
			return true, entity
		}

		b := -2 * (crossterms + x0*dx + y0*dy)
		t := -b / (2 * a)

		if t <= 0 || t >= 1 {
			continue
		}

		closestX := start.X + dx*t
		closestY := start.Y + dy*t
		closestDistance = math.Sqrt(math.Pow(closestX-x0, 2) * +math.Pow(closestY-y0, 2))

		if closestDistance <= entity.Radius+start.Radius+1 {
			return true, entity
		}
	}
	return false, Entity{}
}

// NearestPlanetsByDistance orders all planets based on their proximity
// to a given ship from nearest for farthest
func (gameMap Map) NearestPlanetsByDistance(ship Ship) []Planet {
	planets := gameMap.Planets

	for i := 0; i < len(planets); i++ {

		planets[i].Distance = ship.CalculateDistanceTo(planets[i].Entity)
	}

	sort.Sort(byDist(planets))

	return planets
}

type byDist []Planet

func (a byDist) Len() int           { return len(a) }
func (a byDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDist) Less(i, j int) bool { return a[i].Distance < a[j].Distance }
