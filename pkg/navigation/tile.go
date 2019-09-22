package navigation

import (
	"fmt"
	"math"

	"github.com/metalblueberry/halite-bot/pkg/astar"
)

type TileType int

const (
	Empty      TileType = 0
	SafeMargin TileType = 1
	Walked     TileType = 2
	ShotRange  TileType = 3
	ShotRange2 TileType = 6
	ShotRange3 TileType = 9
	Ship       TileType = 1000
	Blocked    TileType = 1000000
)

func (t TileType) String() string {
	repr := map[TileType]string{
		Empty:      "O",
		Walked:     "*",
		SafeMargin: "+",
		ShotRange:  "#",
		ShotRange2: "%",
		ShotRange3: "@",
		Ship:       "V",
		Blocked:    "X",
	}
	return repr[t]
}

type Weighter interface{
	GetWeight(float64, float64) float64
}

type Tile struct {
	Type TileType
	X    float64
	Y    float64
	Grid *Grid
}

func (t *Tile) String() string {
	return fmt.Sprintf("x:%f y:%f", t.X, t.Y)
}

func (t *Tile) DistanceTo(other *Tile) float64 {
	return math.Sqrt(math.Pow(float64(t.X-other.X), 2) + math.Pow(float64(t.Y-other.Y), 2))
}

// PathNeighbors returns the direct neighboring nodes of this node which
// can be pathed to.
func (t *Tile) PathNeighbors() []astar.Pather {
	neighbors := make([]astar.Pather, 0, 4)

	appendIfNotNull := func(item astar.Pather) {
		switch v := item.(type) {
		case *Tile:
			if v != nil {
				neighbors = append(neighbors, v)
			}
		default:
			return
		}
	}

	up := t.Grid.GetTile(t.X, t.Y-1)
	appendIfNotNull(up)
	down := t.Grid.GetTile(t.X, t.Y+1)
	appendIfNotNull(down)
	left := t.Grid.GetTile(t.X-1, t.Y)
	appendIfNotNull(left)
	right := t.Grid.GetTile(t.X+1, t.Y)
	appendIfNotNull(right)

	upleft := t.Grid.GetTile(t.X-1, t.Y-1)
	appendIfNotNull(upleft)
	upright := t.Grid.GetTile(t.X+1, t.Y-1)
	appendIfNotNull(upright)
	downleft := t.Grid.GetTile(t.X-1, t.Y+1)
	appendIfNotNull(downleft)
	downright := t.Grid.GetTile(t.X+1, t.Y+1)
	appendIfNotNull(downright)
	return neighbors
}

// PathNeighborCost calculates the exact movement cost to neighbor nodes.
func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	return t.DistanceTo(toT) * (float64(toT.Type) + 1)
}

// PathEstimatedCost is a heuristic method for estimating movement costs
// between non-adjacent nodes.
func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	return t.DistanceTo(toT)
}
