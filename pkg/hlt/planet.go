package hlt

import "strconv"

// Planet object from which Halite is mined
type Planet struct {
	Entity
	NumDockingSpots    float64
	NumDockedShips     float64
	CurrentProduction  float64
	RemainingResources float64
	DockedShipIDs      []int
	Owned              float64
	//DockedShips        []Ship
	//Distance           float64
}

// ParsePlanet from a slice of game state tokens
func ParsePlanet(tokens []string) (Planet, []string) {
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

	planetEntity := Entity{
		x:      planetX,
		y:      planetY,
		radius: planetRadius,
		health: planetHealth,
		owner:  planetOwner,
		id:     planetID,
	}

	planet := Planet{
		NumDockingSpots:    planetNumDockingSpots,
		NumDockedShips:     planetNumDockedShips,
		CurrentProduction:  planetCurrentProduction,
		RemainingResources: planetRemainingResources,
		DockedShipIDs:      nil,
		Owned:              planetOwned,
		Entity:             planetEntity,
	}

	for i := 0; i < int(planetNumDockedShips); i++ {
		dockedShipID, _ := strconv.Atoi(tokens[11+i])
		planet.DockedShipIDs = append(planet.DockedShipIDs, dockedShipID)
	}
	return planet, tokens[11+int(planetNumDockedShips):]
}
