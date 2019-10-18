package control

import (
	//log "github.com/sirupsen/logrus"

	"strconv"

	halitedebug "github.com/metalblueberry/Halite-debug/pkg/client"

	"github.com/metalblueberry/halite-bot/pkg/hlt"
	"github.com/metalblueberry/halite-bot/pkg/navigation"
	"github.com/metalblueberry/halite-bot/pkg/twoD"
)

func (c *Commander) FindTarget(pilot *Pilot) twoD.Positioner {
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

func (gameMap *Commander) Navigate(pilot *Pilot, target twoD.Positioner) string {
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
		return ""
	}

	position := GetDirectionFromPath(gameMap.gameMap, pilot, path, 9)
	//log.Printf("position %s", position)

	halitedebug.Line(twoD.NewLine(pilot, position), "nextStep")

	return pilot.NavigateBasic(position)
}

// GetDirectionFromPath returns the tile at which you can move in straight line at the desired speed
func GetDirectionFromPath(gameMap hlt.Map, pilot *Pilot, path []*navigation.Tile, speed float64) *navigation.Tile {
	previousTarget := path[0]
	for _, tile := range path[1:] {
		totalDistance := tile.DistanceTo(pilot)
		if totalDistance > speed {
			return previousTarget
		}
		blocked, collider := gameMap.ObstaclesBetween(pilot, tile, pilot.ID())
		halitedebug.Line(twoD.NewLine(pilot, tile), "direction", "block"+strconv.FormatBool(blocked))
		if blocked {
			halitedebug.Circle(collider, "collider")
			return previousTarget
		}
		previousTarget = tile
	}
	return previousTarget
}
