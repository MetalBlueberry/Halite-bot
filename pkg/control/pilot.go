package control

import (
	"github.com/metalblueberry/halite-bot/pkg/hlt"
	"github.com/metalblueberry/halite-bot/pkg/twoD"
)

type Pilot struct {
	hlt.Ship
	Command         string
	lastTurnUpdated int
	target          twoD.Positioner
}

func NewPilot() *Pilot {
	return &Pilot{}
}

func (pilot *Pilot) SetShip(ship hlt.Ship) {
	pilot.Ship = ship
}
