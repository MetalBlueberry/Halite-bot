package hlt

import (
	"log"

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

func AstarStrategy(ship Ship, gameMap Map) string {
	planets := gameMap.NearestPlanetsByDistance(ship)

	for i := 0; i < len(planets); i++ {
		planet := planets[i]
		if (planet.Owned == 0 || planet.Owner == gameMap.MyID) && planet.NumDockedShips < planet.NumDockingSpots && planet.ID%2 == ship.ID%2 {
			if ship.CanDock(planet) {
				return ship.Dock(planet)
			}
			toEntity := ship.ClosestPointTo(planet.Entity, planet.Radius+1)

			from := gameMap.Grid.GetTile(ship.X, ship.Y)
			to := gameMap.Grid.GetTile(toEntity.X, toEntity.Y)
			path, _, _ := gameMap.Grid.Path(from, to, 100)

			log.Printf("Path %s", path)

			position := navigation.GetDirectionFromPath(path, 7)

			return ship.Navigate(Entity{
				X: position.X,
				Y: position.Y,
			}, gameMap)
		}
	}

	return ""
}
