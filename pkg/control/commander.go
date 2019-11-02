package control

import (
	"context"
	"sort"

	halitedebug "github.com/metalblueberry/Halite-debug/pkg/client"
	"github.com/metalblueberry/halite-bot/pkg/hlt"
	"github.com/metalblueberry/halite-bot/pkg/navigation"
	"github.com/metalblueberry/halite-bot/pkg/twoD"
)

// Commander persist between turns and manage the different ships
type Commander struct {
	gameMap hlt.Map

	currentTurn int

	Grid    *navigation.Grid
	Planets map[int]*PlanetStats
	Pilots  map[int]*Pilot
}

func (c *Commander) PreCalculations() {

	for _, player := range c.gameMap.Players {
		for _, ship := range player.Ships {
			closestPlanet := c.GetPlanetsByDistance(ship)[0]
			closestPlanet.InOrbitShips = append(closestPlanet.InOrbitShips, ship)

			if ship.Owner() == c.gameMap.MyID {
				c.Pilots[ship.ID()].ClosestPlanet = closestPlanet
			}
		}
	}
}

func (c *Commander) Command(ctx context.Context) {

	c.PreCalculations()

	for _, pilot := range c.GetPilotsByHealth() {
		if ctx.Err() != nil {
			return
		}
		if pilot.DockingStatus != hlt.UNDOCKED {
			continue
		}

		target := c.FindTarget(pilot)

		if target == nil {
			continue
		}

		switch targetType := target.(type) {
		case *PlanetStats:
			targetType.PilotsInTheWay += 1.0
			if pilot.CanDock(targetType.Planet) {
				pilot.Command = pilot.Dock(targetType.Planet)
				continue
			}
		default:
			if pilot.DockingStatus != hlt.UNDOCKED {
				pilot.Command = pilot.Undock()
				continue
			}
		}

		//log.Printf("position %s", position)

		path, err := c.CalculatePath(pilot, target)
		if err != nil {
			continue
		}

		turnpath := GetPathForTurn(c.gameMap, pilot, path, 8)
		for _, step := range turnpath {
			step.Type = navigation.Blocked
		}

		if len(turnpath) == 0 {
			continue
		}

		halitedebug.Line(twoD.NewLine(pilot, turnpath[len(turnpath)-1]), "nextStep")

		pilot.Command = pilot.NavigateBasic(turnpath[len(turnpath)-1])
	}
}

type byHealth []*Pilot

func (a byHealth) Len() int      { return len(a) }
func (a byHealth) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byHealth) Less(i, j int) bool {
	return a[i].Health() < a[j].Health() || (a[i].Health() == a[j].Health() && a[i].ID() < a[j].ID())
}

func NewCommander() *Commander {
	return &Commander{
		Planets: make(map[int]*PlanetStats),
		Pilots:  make(map[int]*Pilot),
	}
}

func (c *Commander) CommandQueue() []string {
	commandQueue := make([]string, 0, len(c.Pilots))
	for _, pilot := range c.Pilots {
		commandQueue = append(commandQueue, pilot.Command)
	}
	return commandQueue
}

func (c *Commander) Players() []hlt.Player {
	return c.gameMap.Players
}

func (c *Commander) Me() hlt.Player {
	return c.gameMap.Players[c.gameMap.MyID]
}

func (c *Commander) SetMap(Map hlt.Map, turn int) {
	c.currentTurn = turn
	c.gameMap = Map

	c.findPilotShips()
	c.findPlanetsStats()
	c.generateGrid()
}

func (c *Commander) GetPilots() []*Pilot {
	pilots := make([]*Pilot, 0, len(c.Pilots))
	for _, pilot := range c.Pilots {
		pilots = append(pilots, pilot)
	}
	return pilots
}

func (c *Commander) GetPlanets() []*PlanetStats {
	planets := make([]*PlanetStats, 0, len(c.Planets))
	for _, stats := range c.Planets {
		planets = append(planets, stats)
	}
	return planets
}

func (c *Commander) GetPilotsByHealth() []*Pilot {
	pilots := c.GetPilots()
	sort.Sort(byHealth(pilots))
	return pilots
}

func (c *Commander) GetPlanetsByImportance(pilot *Pilot) []*PlanetStats {
	planets := c.GetPlanets()

	for _, planet := range planets {
		planet.setValueToImportance(pilot)
	}

	sort.Sort(sort.Reverse(byValue(planets)))
	return planets
}

func (c *Commander) GetPlanetsByDistance(pos twoD.Positioner) []*PlanetStats {
	planets := c.GetPlanets()

	for _, planet := range planets {
		planet.setValueToDistance(pos)
	}

	sort.Sort(byValue(planets))
	return planets
}

func (c *Commander) findPlanetsStats() {
	for _, planet := range c.gameMap.Planets {
		stats, exist := c.Planets[planet.ID()]
		if !exist {
			stats = NewPlanetStats()
			c.Planets[planet.ID()] = stats
		}

		stats.SetPlanet(planet)
		stats.lastTurnUpdated = c.currentTurn
		stats.InOrbitShips = make([]hlt.Ship, 0)
	}
	c.removeDeadPlanets()
}

func (c *Commander) removeDeadPlanets() {
	for planetID, stats := range c.Planets {
		if stats.lastTurnUpdated != c.currentTurn {
			delete(c.Planets, planetID)
		}
	}
}

func (c *Commander) findPilotShips() {
	for _, ship := range c.Me().Ships {
		pilot, exist := c.Pilots[ship.ID()]
		if !exist {
			pilot = NewPilot()
			c.Pilots[ship.ID()] = pilot
		}

		pilot.SetShip(ship)
		pilot.lastTurnUpdated = c.currentTurn
	}
	c.removeDeadPilots()
}

func (c *Commander) removeDeadPilots() {
	for shipID, pilot := range c.Pilots {
		if pilot.lastTurnUpdated != c.currentTurn {
			delete(c.Pilots, shipID)
		}
	}
}

func (c *Commander) generateGrid() {
	c.Grid = navigation.NewGrid(c.gameMap.Width, c.gameMap.Height)
	for _, player := range c.gameMap.Players {
		for _, ship := range player.Ships {
			x, y := ship.Position()
			if player.ID == c.gameMap.MyID {
				c.Grid.PaintShip(x, y, 0)
			} else {
				c.Grid.PaintShip(x, y, 5)
			}
		}
	}
	for _, planet := range c.gameMap.Planets {
		c.Grid.PaintPlanet(planet.Circle())
	}
}
