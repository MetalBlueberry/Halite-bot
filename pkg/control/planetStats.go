package control

import (
	"github.com/metalblueberry/halite-bot/pkg/hlt"
	"github.com/metalblueberry/halite-bot/pkg/twoD"
)

type PlanetStats struct {
	hlt.Planet
	FlyingTo        []*Pilot
	lastTurnUpdated int

	StaticValue float64
	Value       float64
	//Distance    float64

	InOrbitShips []hlt.Ship

	PilotsInTheWay float64
}

type byValue []*PlanetStats

func (a byValue) Len() int           { return len(a) }
func (a byValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byValue) Less(i, j int) bool { return a[i].Value < a[j].Value }

func NewPlanetStats() *PlanetStats {
	return &PlanetStats{
		InOrbitShips: make([]hlt.Ship, 0),
	}
}
func (stats *PlanetStats) SetPlanet(planet hlt.Planet) {
	stats.Planet = planet
	stats.PilotsInTheWay = 0
}

func (stats *PlanetStats) setValueToImportance(pilot *Pilot) {
	_, _, radius := stats.Circle()

	stats.Value = stats.Health() / (radius * 255) *
		(radius - stats.PilotsInTheWay - twoD.Distance(stats, pilot)/hlt.Constants["MAX_SPEED"].(float64))
}

func (stats *PlanetStats) setValueToDistance(pos twoD.Positioner) {
	stats.Value = twoD.Distance(stats, pos)
}
