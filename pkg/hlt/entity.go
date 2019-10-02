package hlt

import (
	"fmt"
	"math"
	"strconv"
)

// DockingStatus represents possible ship.DockingStatus values
type DockingStatus int

const (
	// UNDOCKED ship.DockingStatus value
	UNDOCKED DockingStatus = iota
	// DOCKING ship.DockingStatus value
	DOCKING
	// DOCKED ship.DockingStatus value
	DOCKED
	// UNDOCKING ship.DockingStatus value
	UNDOCKING
)

type Entitier interface {
	Positioner
	Circle() (x, y, radius float64)
	Health() float64
	Owner() int
	ID() int
	CalculateDistanceTo(other Positioner) float64
	CalculateAngleTo(other Positioner) float64
	CalculateRadAngleTo(target Positioner) float64
}

// Entity captures spacial and ownership state for Planets and Ships
type Entity struct {
	x      float64
	y      float64
	radius float64
	health float64
	owner  int
	id     int
}

func (e *Entity) Position() (x, y float64) {
	return e.x, e.y
}

func (e *Entity) Circle() (x, y, radius float64) {
	return e.x, e.y, e.radius
}

func (e *Entity) Health() float64 {
	return e.health
}
func (e *Entity) Owner() int {
	return e.owner
}
func (e *Entity) ID() int {
	return e.id
}

// Planet object from which Halite is mined
type Planet struct {
	*Entity
	NumDockingSpots    float64
	NumDockedShips     float64
	CurrentProduction  float64
	RemainingResources float64
	DockedShipIDs      []int
	DockedShips        []Ship
	InOrbitShips       []Ship
	Owned              float64
	Distance           float64
}

// Ship is a player controlled Entity made for the purpose of doing combat and mining Halite
type Ship struct {
	*Entity
	VelX float64
	VelY float64

	PlanetID        int
	Planet          Planet
	DockingStatus   DockingStatus
	DockingProgress float64
	WeaponCooldown  float64
}

func (entity Entity) UniqueID() string {
	return fmt.Sprintf("%f%f", entity.x, entity.y)
}

// CalculateDistanceTo returns a euclidean distance to the target
func (entity Entity) CalculateDistanceTo(target Positioner) float64 {
	x, y := target.Position()
	dx := x - entity.x
	dy := y - entity.y

	return math.Sqrt(dx*dx + dy*dy)
}

// CalculateAngleTo returns an angle in degrees to the target
func (entity *Entity) CalculateAngleTo(target Positioner) float64 {
	return RadToDeg(entity.CalculateRadAngleTo(target))
}

// CalculateRadAngleTo returns an angle in radians to the target
func (entity *Entity) CalculateRadAngleTo(target Positioner) float64 {
	tx, ty := target.Position()
	dx := tx - entity.x
	dy := ty - entity.y

	return math.Atan2(dy, dx)
}

// ClosestPointTo returns the closest point that is at least minDistance from the target
func (entity *Entity) ClosestPointTo(target Entitier, minDitance float64) *Entity {
	tx, ty, radius := target.Circle()
	dist := radius + minDitance
	norm := target.CalculateDistanceTo(entity)
	x := dist*(entity.x-tx)/norm + tx
	y := dist*(entity.y-ty)/norm + ty
	return &Entity{
		x:      x,
		y:      y,
		radius: 0,
		health: 0,
		owner:  -1,
		id:     -1,
	}
}

// VectorTo returns the distance and direction between two entities.
func (entity *Entity) VectorTo(other Positioner) Entity {
	x, y := other.Position()
	return Entity{
		x: x - entity.x,
		y: y - entity.y,
	}
}

// Rotate returns the position rotated an angle from the point origin
func (entity *Entity) Rotate(origin Entity, angle float64) Entity {
	// x2=r−u=cosβx1−sinβy1y2=t+s=sinβx1+cosβy1
	vector := origin.VectorTo(entity)
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	vector.x = cos*entity.x - sin*entity.y + origin.x
	vector.y = sin*entity.x + cos*entity.y + origin.y
	return vector
}

// ParseShip from a slice of game state tokens
func ParseShip(playerID int, tokens []string) (*Ship, []string) {
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

	shipEntity := &Entity{
		x:      shipX,
		y:      shipY,
		radius: .5,
		health: shipHealth,
		owner:  playerID,
		id:     shipID,
	}

	ship := &Ship{
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

// ParsePlanet from a slice of game state tokens
func ParsePlanet(tokens []string) (*Planet, []string) {
	planetID, _ := strconv.Atoi(tokens[0])
	planetX, _ := strconv.ParseFloat(tokens[1], 64)
	planetY, _ := strconv.ParseFloat(tokens[2], 64)
	planetHealth, _ := strconv.ParseFloat(tokens[3], 64)
	planetRadius, _ := strconv.ParseFloat(tokens[4], 64)
	planetNumDockingSpots, _ := strconv.ParseFloat(tokens[5], 64)
	planetCurrentProduction, _ := strconv.ParseFloat(tokens[6], 64)
	planetRemainingResources, _ := strconv.ParseFloat(tokens[7], 64)
	planetOwned, _ := strconv.ParseFloat(tokens[8], 64)
	planetOwner, _ := strconv.Atoi(tokens[9])
	planetNumDockedShips, _ := strconv.ParseFloat(tokens[10], 64)

	planetEntity := &Entity{
		x:      planetX,
		y:      planetY,
		radius: planetRadius,
		health: planetHealth,
		owner:  planetOwner,
		id:     planetID,
	}

	planet := &Planet{
		NumDockingSpots:    planetNumDockingSpots,
		NumDockedShips:     planetNumDockedShips,
		CurrentProduction:  planetCurrentProduction,
		RemainingResources: planetRemainingResources,
		DockedShipIDs:      nil,
		DockedShips:        nil,
		Owned:              planetOwned,
		Entity:             planetEntity,
	}

	for i := 0; i < int(planetNumDockedShips); i++ {
		dockedShipID, _ := strconv.Atoi(tokens[11+i])
		planet.DockedShipIDs = append(planet.DockedShipIDs, dockedShipID)
	}
	return planet, tokens[11+int(planetNumDockedShips):]
}

// IntToDockingStatus converts an int to a DockingStatus
func IntToDockingStatus(i int) DockingStatus {
	statuses := [4]DockingStatus{UNDOCKED, DOCKING, DOCKED, UNDOCKING}
	return statuses[i]
}

// Thrust generates a string describing the ship's intension to move during the current turn
func (ship *Ship) Thrust(magnitude float64, angle float64) string {
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
func (ship *Ship) Dock(planet *Planet) string {
	return fmt.Sprintf("d %s %s", strconv.Itoa(ship.id), strconv.Itoa(planet.id))
}

// Undock generates a string describing the ship's intension to undock during the current turn
func (ship *Ship) Undock() string {
	return fmt.Sprintf("u %s", strconv.Itoa(ship.id))
}

// NavigateBasic demonstrates how the player might move ships through space
func (ship *Ship) NavigateBasic(target Positioner, gameMap Map) string {
	distance := ship.CalculateDistanceTo(target)
	safeDistance := distance - ship.Entity.radius - .1

	angle := ship.CalculateAngleTo(target)
	speed := 7.0
	if distance < 10 {
		speed = 3.0
	}

	speed = math.Min(speed, safeDistance)
	return ship.Thrust(speed, angle)
}

// NavigateBasic demonstrates how the player might move ships through space
func (ship *Ship) NavigateBasic2(target Positioner, gameMap Map) string {
	distance := ship.CalculateDistanceTo(target)
	safeDistance := distance // - ship.Entity.Radius - target.Radius - .1 //disbled  as pathfinding is handling this

	angle := ship.CalculateAngleTo(target)
	speed := 7.0

	speed = math.Min(speed, safeDistance)
	return ship.Thrust(speed, angle)
}

// CanDock indicates that a ship is close enough to a given planet to dock
func (ship *Ship) CanDock(planet Entitier) bool {
	dist := ship.CalculateDistanceTo(planet)

	_, _, radius := planet.Circle()
	return dist <= (ship.radius + radius + 4)
}

// Navigate demonstrates how the player might negotiate obsticles between
// a ship and its target
func (ship *Ship) Navigate(target Entitier, gameMap Map) string {
	ob, _ := gameMap.ObstaclesBetween(ship, target)

	if !ob {
		return ship.NavigateBasic(target, gameMap)
	}

	tx, ty := target.Position()
	x0 := math.Min(ship.x, tx)
	x2 := math.Max(ship.x, tx)
	y0 := math.Min(ship.y, ty)
	y2 := math.Max(ship.y, ty)

	dx := (x2 - x0) / 5
	dy := (y2 - y0) / 5
	bestdist := 1000.0
	bestTarget := target

	for x1 := x0; x1 <= x2; x1 += dx {
		for y1 := y0; y1 <= y2; y1 += dy {
			intermediateTarget := &Entity{
				x:      x1,
				y:      y1,
				radius: 0,
				health: 0,
				owner:  0,
				id:     -1,
			}
			ob1, _ := gameMap.ObstaclesBetween(ship, intermediateTarget)
			if !ob1 {
				ob2, _ := gameMap.ObstaclesBetween(intermediateTarget, target)
				if !ob2 {
					totdist := math.Sqrt(math.Pow(x1-x0, 2)+math.Pow(y1-y0, 2)) + math.Sqrt(math.Pow(x1-x2, 2)+math.Pow(y1-y2, 2))
					if totdist < bestdist {
						bestdist = totdist
						bestTarget = intermediateTarget

					}
				}
			}
		}
	}

	return ship.NavigateBasic(bestTarget, gameMap)
}
