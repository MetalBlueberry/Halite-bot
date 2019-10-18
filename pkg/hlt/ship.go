package hlt

import (
	"fmt"
	"math"
	"strconv"

	"github.com/metalblueberry/halite-bot/pkg/twoD"
)

// Ship is a player controlled Entity made for the purpose of doing combat and mining Halite
type Ship struct {
	Entity
	VelX float64
	VelY float64

	PlanetID        int
	Planet          Planet
	DockingStatus   DockingStatus
	DockingProgress float64
	WeaponCooldown  float64
}

// ParseShip from a slice of game state tokens
func ParseShip(playerID int, tokens []string) (Ship, []string) {
	shipID, _ := strconv.Atoi(tokens[0])
	shipX, _ := strconv.ParseFloat(tokens[1], 64)
	shipY, _ := strconv.ParseFloat(tokens[2], 64)
	shipHealth, _ := strconv.ParseFloat(tokens[3], 64)
	shipVelX, _ := strconv.ParseFloat(tokens[4], 64)
	shipVelY, _ := strconv.ParseFloat(tokens[5], 64)
	shipDockingStatus, _ := strconv.Atoi(tokens[6])
	shipPlanetID, _ := strconv.Atoi(tokens[7])
	shipDockingProgress, _ := strconv.ParseFloat(tokens[8], 64)
	shipWeaponCooldown, _ := strconv.ParseFloat(tokens[9], 64)

	shipEntity := Entity{
		x:      shipX,
		y:      shipY,
		radius: .5,
		health: shipHealth,
		owner:  playerID,
		id:     shipID,
	}

	ship := Ship{
		PlanetID:        shipPlanetID,
		DockingStatus:   IntToDockingStatus(shipDockingStatus),
		DockingProgress: shipDockingProgress,
		WeaponCooldown:  shipWeaponCooldown,
		VelX:            shipVelX,
		VelY:            shipVelY,
		Entity:          shipEntity,
	}

	return ship, tokens[10:]
}

// Thrust generates a string describing the ship's intension to move during the current turn
func (ship Ship) Thrust(magnitude float64, angle float64) string {
	var boundedAngle int
	if angle > 0.0 {
		boundedAngle = int(math.Floor(angle + .5))
	} else {
		boundedAngle = int(math.Ceil(angle - .5))
	}
	boundedAngle = ((boundedAngle % 360) + 360) % 360
	return fmt.Sprintf("t %s %s %s", strconv.Itoa(ship.id), strconv.Itoa(int(magnitude)), strconv.Itoa(boundedAngle))
}

// Dock generates a string describing the ship's intension to dock during the current turn
func (ship Ship) Dock(planet Planet) string {
	return fmt.Sprintf("d %s %s", strconv.Itoa(ship.id), strconv.Itoa(planet.id))
}

// Undock generates a string describing the ship's intension to undock during the current turn
func (ship Ship) Undock() string {
	return fmt.Sprintf("u %s", strconv.Itoa(ship.id))
}

// NavigateBasic demonstrates how the player might move ships through space
func (ship Ship) NavigateBasic(target twoD.Positioner) string {
	distance := ship.CalculateDistanceTo(target)
	angle := ship.CalculateAngleTo(target)
	speed := 7.0 //TODO: get from environment

	speed = math.Min(speed, distance)
	return ship.Thrust(speed, angle)
}

// CanDock indicates that a ship is close enough to a given planet to dock
func (ship Ship) CanDock(planet Planet) bool {
	owner := planet.Owner()
	if owner != 0 && owner != ship.owner {
		return false
	}
	if planet.NumDockedShips == planet.NumDockingSpots {
		return false
	}
	dist := ship.CalculateDistanceTo(planet)
	_, _, radius := planet.Circle()
	return dist <= (ship.radius + radius + 4)
}

// IntToDockingStatus converts an int to a DockingStatus
func IntToDockingStatus(i int) DockingStatus {
	statuses := [4]DockingStatus{UNDOCKED, DOCKING, DOCKED, UNDOCKING}
	return statuses[i]
}
