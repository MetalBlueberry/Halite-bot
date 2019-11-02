package control

import (
	//log "github.com/sirupsen/logrus"

	"errors"
	"strconv"

	halitedebug "github.com/metalblueberry/Halite-debug/pkg/client"

	"github.com/metalblueberry/halite-bot/pkg/hlt"
	"github.com/metalblueberry/halite-bot/pkg/navigation"
	"github.com/metalblueberry/halite-bot/pkg/twoD"
)

func (c *Commander) FindTarget(pilot *Pilot) twoD.Positioner {

	// if Pilot can fight, look for trouble
	if pilot.Health() > 255/4 && pilot.ClosestPlanet.Owner() == c.gameMap.MyID {
		for _, inOrbitShip := range pilot.ClosestPlanet.InOrbitShips {
			if inOrbitShip.Owner() != c.gameMap.MyID {
				return inOrbitShip
			}
		}
	}

	planets := c.GetPlanetsByImportance(pilot)

	for _, planet := range planets {
		if (planet.Owned == 0 || planet.Owner() == c.gameMap.MyID) && planet.NumDockedShips < planet.NumDockingSpots {

			// Select enemy ship if is close to the planet
			for _, ship := range planet.InOrbitShips {
				if ship.Owner() == c.gameMap.MyID {
					continue
				}
				_, _, r := planet.Circle()
				if twoD.Distance(ship, planet)-r < 3*hlt.Constants["DOCK_RADIUS"].(float64) {
					return ship
				}
			}

			return planet
		}
		if planet.Owner() != c.gameMap.MyID {
			//TODO: find nearest ship
			for _, enemy := range planet.DockedShipIDs {
				return c.gameMap.Ships[enemy]
			}
		}
	}

	return nil
}

func (gameMap *Commander) CalculatePath(pilot *Pilot, target twoD.Positioner) ([]*navigation.Tile, error) {
	// If target has a radius, calculate a near position with a margin
	switch targetType := target.(type) {
	case twoD.Circler:
		target = twoD.ClosestPointTo(pilot, targetType, 2)
	}

	//log.Printf("Planet %v, Point %v", planet.Entity, target)
	from := gameMap.Grid.GetTile(pilot.Position())
	to := gameMap.Grid.GetTile(target.Position())
	_, _, found, path := gameMap.Grid.Path(from, to, 300)

	if !found {
		halitedebug.Line(twoD.NewLine(from, to), "notFound")
	}

	//Print debug information
	{
		previous := from
		for _, t := range path {
			halitedebug.Line(twoD.NewLine(previous, t), "path")
			previous = t
		}
	}

	if len(path) == 0 {
		return nil, errors.New("Path not found")
	}
	return path, nil
}

// GetPathForTurn returns the tile at which you can move in straight line at the desired speed
func GetPathForTurn(gameMap hlt.Map, pilot *Pilot, path []*navigation.Tile, speed float64) []*navigation.Tile {
	for i, tile := range path[0:] {
		totalDistance := tile.DistanceTo(pilot)
		if totalDistance > speed {
			return path[0:i]
		}
		blocked, collider := gameMap.ObstaclesBetween(pilot, tile, pilot.ID())
		halitedebug.Line(twoD.NewLine(pilot, tile), "direction", "block"+strconv.FormatBool(blocked))
		if blocked {
			halitedebug.Circle(collider, "collider")
			return path[0:i]
		}
	}
	return path
}
