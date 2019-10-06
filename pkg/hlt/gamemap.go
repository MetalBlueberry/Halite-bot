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
	Planets             []*Planet
	Players             []*Player
	Entities            []Entitier
	Grid                *navigation.Grid
}

// Player has an ID for establishing ownership, and a number of ships
type Player struct {
	ID    int
	Ships []*Ship
}

// ParsePlayer from a slice of game state tokens
func ParsePlayer(tokens []string) (*Player, []string) {
	playerID, _ := strconv.Atoi(tokens[0])
	playerNumShips, _ := strconv.ParseFloat(tokens[1], 64)

	player := &Player{
		ID:    playerID,
		Ships: []*Ship{},
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
		Players:  make([]*Player, numPlayers),
		Entities: make([]Entitier, 0),
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
				gameMap.Grid.PaintShip(ship.x, ship.y, 0)
			} else {
				gameMap.Grid.PaintShip(ship.x, ship.y, 5)
			}
		}
	}

	numPlanets, _ := strconv.Atoi(tokens[0])
	gameMap.Planets = make([]*Planet, 0, numPlanets)
	tokens = tokens[1:]

	for i := 0; i < numPlanets; i++ {
		planet, tokensnew := ParsePlanet(tokens)
		tokens = tokensnew
		gameMap.Planets = append(gameMap.Planets, planet)
		gameMap.Entities = append(gameMap.Entities, planet.Entity)
		gameMap.Grid.PaintPlanet(planet.Entity.x, planet.Entity.y, planet.Entity.radius)
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

// ObstaclesBetween demonstrates how the player might determine if the path
// between two enitities is clear
func (gameMap Map) ObstaclesBetween(start Entitier, end Entitier) (bool, Entitier) {
	return ObstaclesBetween(start, end, gameMap.Entities)
}

func ObstaclesBetween(start Entitier, end Entitier, Obstacles []Entitier) (bool, Entitier) {
	x1, y1 := start.Position()
	x2, y2 := end.Position()
	dx := x2 - x1
	dy := y2 - y1
	a := dx*dx + dy*dy + 1e-8
	crossterms := x1*x1 - x1*x2 + y1*y1 - y1*y2

	for i := 0; i < len(Obstacles); i++ {
		entity := Obstacles[i]
		if entity.ID() == start.ID() || entity.ID() == end.ID() {
			continue
		}

		x0, y0, radius := entity.Circle()

		closestDistance := end.CalculateDistanceTo(entity)
		if closestDistance < radius {
			return true, entity
		}

		b := -2 * (crossterms + x0*dx + y0*dy)
		t := -b / (2 * a)

		if t <= 0 || t >= 1 {
			continue
		}

		sx, sy, sradius := start.Circle()
		closestX := sx + dx*t
		closestY := sy + dy*t
		closestDistance = math.Sqrt(math.Pow(closestX-x0, 2) * +math.Pow(closestY-y0, 2))

		if closestDistance <= radius+sradius {
			return true, entity
		}
	}
	return false, nil
}

func (gameMap Map) ObstaclesBetween2(start Entitier, end Entitier) (bool, Entitier) {
	return ObstaclesBetween2(start, end, gameMap.Entities)
}

func ObstaclesBetween2(start Entitier, end Entitier, obstacles []Entitier) (bool, Entitier) {
	_, _, r1 := start.Circle()
	//_, _ := end.Position()
	StartToEnd := Distance(start, end)
	for i := 0; i < len(obstacles); i++ {
		entity := obstacles[i]
		if entity.ID() == start.ID() || entity.ID() == end.ID() {
			continue
		}

		_, _, r := entity.Circle()
		margin := r1 + r
		dist := DistancePointToLine(start, end, entity)
		if dist > margin {
			continue
		}

		projection := Project(start, end, entity)
		relative := projection / StartToEnd
		if relative > 0 && relative < 1 {
			return true, entity
		}

		endDist := Distance(end, entity)
		if endDist < margin {
			return true, entity
		}
	}
	return false, nil
}

func DistancePointToLine(A, B, P Positioner) float64 {
	Px, Py := P.Position()
	m, n := LineMNFrom(A, B)
	if math.IsInf(m, 0) {
		return math.Abs(Px - n)
	}
	return math.Abs(Px*m+n-Py) / math.Sqrt(m*m+1)
}

func Project(A, B, P Positioner) float64 {
	Px, Py := P.Position()
	Ax, Ay := A.Position()
	Ux, Uy := UnitVector(A, B)
	return (Px-Ax)*Ux + (Py-Ay)*Uy
}

func UnitVector(A, B Positioner) (x, y float64) {
	Ax, Ay := A.Position()
	Bx, By := B.Position()
	mod := Distance(A, B)
	return (Bx - Ax) / mod, (By - Ay) / mod
}

func Distance(A, B Positioner) float64 {
	Ax, Ay := A.Position()
	Bx, By := B.Position()
	return math.Sqrt(math.Pow(Ax-Bx, 2) + math.Pow(Ay-By, 2))
}

// LineMNFrom solves line equation for y = mx+n
// if m is infinity, n will be x = n
func LineMNFrom(A, B Positioner) (m, n float64) {
	Ax, Ay := A.Position()
	Bx, By := B.Position()
	m = (By - Ay) / (Bx - Ax)
	if math.IsInf(m, 0) {
		return m, Bx
	}
	n = Ay - Ax*m
	//n = (Ay*Bx - Ax*By) / (Bx - Ax)
	return m, n
}

// NearestPlanetsByDistance orders all planets based on their proximity
// to a given ship from nearest for farthest
func (gameMap Map) NearestPlanetsByDistance(ship *Ship) []*Planet {
	planets := gameMap.Planets

	for i := 0; i < len(planets); i++ {

		planets[i].Distance = ship.CalculateDistanceTo(planets[i].Entity)
	}

	sort.Sort(byDist(planets))

	return planets
}

type byDist []*Planet

func (a byDist) Len() int           { return len(a) }
func (a byDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDist) Less(i, j int) bool { return a[i].Distance < a[j].Distance }
