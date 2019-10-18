package hlt

import (
	"github.com/metalblueberry/halite-bot/pkg/twoD"
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
	twoD.Circler
	Health() float64
	Owner() int
	ID() int
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

func (e Entity) Position() (x, y float64) {
	return e.x, e.y
}

func (e Entity) Circle() (x, y, radius float64) {
	return e.x, e.y, e.radius
}

func (e Entity) Health() float64 {
	return e.health
}
func (e Entity) Owner() int {
	return e.owner
}
func (e Entity) ID() int {
	return e.id
}

// CalculateDistanceTo returns a euclidean distance to the target
func (entity Entity) CalculateDistanceTo(target twoD.Positioner) float64 {
	return twoD.Distance(entity, target)
}

// CalculateAngleTo returns an angle in degrees to the target
func (entity Entity) CalculateAngleTo(target twoD.Positioner) float64 {
	return twoD.RadToDeg(twoD.CalculateRadAngleTo(entity, target))
}
