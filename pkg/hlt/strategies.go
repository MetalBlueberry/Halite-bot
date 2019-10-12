package hlt

import (
	//log "github.com/sirupsen/logrus"
	halitedebug "github.com/metalblueberry/Halite-debug/pkg/client"
	"strconv"

	"github.com/metalblueberry/halite-bot/pkg/navigation"
)

// StrategyBasicBot demonstrates how the player might direct their ships
// in achieving victory
func StrategyBasicBot(ship *Ship, gameMap Map) string {
	planets := gameMap.NearestPlanetsByDistance(ship)

	for i := 0; i < len(planets); i++ {
		planet := planets[i]
		if (planet.Owned == 0 || planet.owner == gameMap.MyID) && planet.NumDockedShips < planet.NumDockingSpots && planet.id%2 == ship.id%2 {
			if ship.CanDock(planet) {
				return ship.Dock(planet)
			}
			return ship.Navigate(ship.ClosestPointTo(planet, 3), gameMap)
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

func AstarStrategy(ship *Ship, gameMap Commander) string {
	planets := gameMap.NearestPlanetsByDistance(ship)

	if ship.CanDock(planets[0]) {
		return ship.Dock(planets[0])
	}

	target := gameMap.FindTarget(ship, planets)

	if target == nil {
		return ""
	}
	targetEntity, ok := target.(Entitier)
	if ok {
		target = ship.ClosestPointTo(targetEntity, 2)
	}

	return gameMap.Navigate(ship, target)

}

func (gameMap Commander) FindTarget(ship *Ship, planets []*Planet) Positioner {
	for _, planet := range planets {
		if (planet.Owned == 0 || planet.owner == gameMap.MyID) && planet.NumDockedShips < planet.NumDockingSpots {
			return planet
		}
		if planet.owner != gameMap.MyID {
			for _, enemy := range planet.DockedShipIDs {
				return gameMap.Ships[enemy]
			}
		}
	}
	//for _, planet := range planets {
	//if planet.owner != gameMap.MyID {
	//for _, enemy := range planet.DockedShipIDs {
	//return gameMap.Ships[enemy]
	//}
	//}
	//}
	return nil
}

func (gameMap Map) Navigate(ship *Ship, target Positioner) string {
	x, y := target.Position()
	//log.Printf("Planet %v, Point %v", planet.Entity, target)
	from := gameMap.Grid.GetTile(ship.x, ship.y)
	to := gameMap.Grid.GetTile(x, y)
	path, _, _, _ := gameMap.Grid.Path(from, to, -1)

	previous := from
	for _, t := range path {
		halitedebug.Line(NewLine(previous, t), "path")
		previous = t
	}

	position := GetDirectionFromPath(&gameMap, ship, path, 9)
	//log.Printf("position %s", position)

	halitedebug.Line(NewLine(ship, position), "nextStep")

	return ship.NavigateBasic2(&Entity{
		x:      position.X,
		y:      position.Y,
		radius: 0,
	}, gameMap)
}

// GetDirectionFromPath returns the tile at which you can move in straight line at the desired speed
func GetDirectionFromPath(gameMap *Map, ship *Ship, path []*navigation.Tile, speed float64) *navigation.Tile {
	previousTarget := path[0]
	for _, tile := range path[1:] {
		totalDistance := tile.DistanceTo(ship)
		if totalDistance > speed {
			return previousTarget
		}
		targetEntity := PositionEntityFromTile(tile, ship.id, 0)
		blocked, collider := gameMap.ObstaclesBetween2(ship, targetEntity)
		halitedebug.Line(NewLine(ship, targetEntity), "direction", "block"+strconv.FormatBool(blocked))
		if blocked {
			halitedebug.Circle(collider, "collider")
			return previousTarget
		}
		previousTarget = tile
	}
	return previousTarget
}

func PositionEntityFromTile(p Positioner, ID int, radius float64) *Entity {
	x, y := p.Position()
	return &Entity{
		x:      x,
		y:      y,
		radius: radius,
		id:     ID,
	}
}
