package hlt

import (
	//log "github.com/sirupsen/logrus"
	debug "github.com/metalblueberry/sicase/pkg/client"
	"strconv"

	"github.com/metalblueberry/halite-bot/pkg/navigation"
)

// StrategyBasicBot demonstrates how the player might direct their ships
// in achieving victory
func StrategyBasicBot(ship Ship, gameMap Map) string {
	planets := gameMap.NearestPlanetsByDistance(ship)

	for i := 0; i < len(planets); i++ {
		planet := planets[i]
		if (planet.Owned == 0 || planet.Owner == gameMap.MyID) && planet.NumDockedShips < planet.NumDockingSpots && planet.ID%2 == ship.ID%2 {
			if ship.CanDock(planet) {
				return ship.Dock(planet)
			}
			return ship.Navigate(ship.ClosestPointTo(planet.Entity, 3), gameMap)
		}
	}

	return ""
}

type Line struct {
	X1, Y1 float64
	X2, Y2 float64
}

type Positioner interface {
	Position() (x, y float64)
}

func NewLine(a, b Positioner) Line {
	X1, Y1 := a.Position()
	X2, Y2 := b.Position()
	return Line{X1, Y1, X2, Y2}
}

func (l Line) Line() (float64, float64, float64, float64) {
	return l.X1, l.Y1, l.X2, l.Y2
}

func AstarStrategy(ship Ship, gameMap Map) string {
	planets := gameMap.NearestPlanetsByDistance(ship)

	for _, planet := range planets {
		if (planet.Owned == 0 || planet.Owner == gameMap.MyID) && planet.NumDockedShips < planet.NumDockingSpots {
			if ship.CanDock(planet) {
				return ship.Dock(planet)
			}
			target := ship.ClosestPointTo(planet.Entity, 2)

			//log.Printf("Planet %v, Point %v", planet.Entity, target)
			from := gameMap.Grid.GetTile(ship.X, ship.Y)
			to := gameMap.Grid.GetTile(target.X, target.Y)
			path, _, _, _ := gameMap.Grid.Path(from, to, -1)

			previous := from
			for _, t := range path {
				debug.Line(NewLine(previous, t), "path")
				previous = t
			}

			// log.Printf("Ship id %d", ship.ID)

			// log.Printf("Path %s", path)

			position := GetDirectionFromPath(&gameMap, ship, path, 8)
			//log.Printf("position %s", position)

			debug.Line(NewLine(&ship, position), "nextStep")

			return ship.NavigateBasic2(Entity{
				X:      position.X,
				Y:      position.Y,
				Radius: 0,
			}, gameMap)
		}
	}

	//for _, planet := range planets {
	//if planet.Owned == 1 && planet.Owner != gameMap.MyID {
	//target := ship.ClosestPointTo(planet.Entity, 2)
	//for _, docked := range planet.DockedShips {
	//if ship.CalculateDistanceTo(docked.Entity) < ship.CalculateDistanceTo(target) {
	//target = ship.ClosestPointTo(docked.Entity, docked.Entity.Radius+1)
	//}
	//}

	//from := gameMap.Grid.GetTile(ship.X, ship.Y)
	//to := gameMap.Grid.GetTile(target.X, target.Y)
	//_, _, _, path := gameMap.Grid.Path(from, to, 10000)
	//position := GetDirectionFromPath(&gameMap, ship, path, 7)
	//return ship.NavigateBasic2(Entity{
	//X:      position.X,
	//Y:      position.Y,
	//Radius: 0,
	//}, gameMap)
	//}
	//}

	return ""
}

// GetDirectionFromPath returns the tile at which you can move in straight line at the desired speed
func GetDirectionFromPath(gameMap *Map, ship Ship, path []*navigation.Tile, speed float64) *navigation.Tile {
	origin := &ship
	originEntity := PositionEntityFromTile(origin, ship.ID, 1)
	previousTarget := path[0]
	for _, tile := range path[1:] {
		totalDistance := tile.DistanceTo(origin)
		if totalDistance > speed {
			return previousTarget
		}
		targetEntity := PositionEntityFromTile(tile, ship.ID, 1)
		blocked, _ := gameMap.ObstaclesBetween(originEntity, targetEntity)
		debug.Line(NewLine(origin, previousTarget), "direction", "block"+strconv.FormatBool(blocked))
		if blocked {
			return previousTarget
		}
		previousTarget = tile
	}
	return previousTarget
}

func PositionEntityFromTile(p Positioner, ID int, radius float64) Entity {
	x, y := p.Position()
	return Entity{
		X:      x,
		Y:      y,
		Radius: radius,
		ID:     ID,
	}
}
