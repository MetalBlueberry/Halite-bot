package navigation

import (
	"fmt"
	"math"

	"github.com/metalblueberry/halite-bot/pkg/astar"
)

type TileType int

const (
	Empty      TileType = 0
	SafeMargin TileType = 2
	Walked     TileType = 4
	ShotRange  TileType = 6
	ShotRange2 TileType = 12
	ShotRange3 TileType = 18
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

type Weighter interface {
	GetWeight(float64, float64) float64
}

type Tile struct {
	Type TileType
	X    float64
	Y    float64
	Grid *Grid
}

func (t *Tile) Position() (x, y float64) {
	return t.X, t.Y
}

func (t *Tile) String() string {
	return fmt.Sprintf("x:%f y:%f", t.X, t.Y)
}

type Positioner interface {
	Position() (x, y float64)
}

func (t *Tile) DistanceTo(other Positioner) float64 {
	x2, y2 := other.Position()
	return math.Sqrt(math.Pow(float64(t.X-x2), 2) + math.Pow(float64(t.Y-y2), 2))
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
