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

	for _, planet := range planets {
		if (planet.Owned == 0 || planet.Owner == gameMap.MyID) && planet.NumDockedShips < planet.NumDockingSpots && planet.ID%2 == ship.ID%2 {
			if ship.CanDock(planet) {
				return ship.Dock(planet)
			}
			target := ship.ClosestPointTo(planet.Entity, 2)

			log.Printf("Planet %v, Point %v", planet.Entity, target)
			from := gameMap.Grid.GetTile(ship.X, ship.Y)
			to := gameMap.Grid.GetTile(target.X, target.Y)
			path, _, _, _ := gameMap.Grid.Path(from, to, -1)

			// log.Printf("Ship id %d", ship.ID)
			// if ship.ID == 0 {
				 log.Println(gameMap.Grid.PrintDebugPath(path, from, to))
			// }

			// log.Printf("Path %s", path)

			position := navigation.GetDirectionFromPath(path, 8)
			log.Printf("position %s", position)

			return ship.NavigateBasic2(Entity{
				X:      position.X,
				Y:      position.Y,
				Radius: 0,
			}, gameMap)
		}
	}

	for _, planet := range planets {
		if (planet.Owned == 1 && planet.Owner != gameMap.MyID) && planet.ID%2 == ship.ID%2 {
			target := ship.ClosestPointTo(planet.Entity, 2)
			for _, docked := range planet.DockedShips {
				if ship.CalculateDistanceTo(docked.Entity) < ship.CalculateDistanceTo(target) {
					target = ship.ClosestPointTo(docked.Entity, docked.Entity.Radius + 1)
				}
			}

			from := gameMap.Grid.GetTile(ship.X, ship.Y)
			to := gameMap.Grid.GetTile(target.X, target.Y)
			_, _, _, path := gameMap.Grid.Path(from, to, 10000)
			position := navigation.GetDirectionFromPath(path, 7)
			return ship.NavigateBasic2(Entity{
				X:      position.X,
				Y:      position.Y,
				Radius: 0,
			}, gameMap)
		}
	}

	return ""
}
