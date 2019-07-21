package navigation

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/metalblueberry/halite-bot/pkg/astar"
)

type Grid struct {
	Width, Height int
	Tiles         [][]*Tile
}

func NewGrid(Width, Height int) *Grid {
	grid := &Grid{
		Width:  Width,
		Height: Height,
		Tiles:  make([][]*Tile, Height, Height),
	}
	for row := range grid.Tiles {
		grid.Tiles[row] = make([]*Tile, Width, Width)
		for col := range grid.Tiles[row] {
			grid.Tiles[row][col] = &Tile{
				Grid: grid,
				X:    col,
				Y:    row,
			}
		}
	}
	return grid
}

func (g *Grid) PaintShip(X float64, Y float64, shotRange float64) {
	g.Paint(X, Y, shotRange, ShotRange)
	g.Paint(X, Y, 0.5, Ship)
}

func (g *Grid) PaintPlanet(X float64, Y float64, radius float64) {
	g.Paint(X, Y, radius, Blocked)
}

type PaintMode int

const (
	EmptyOnly PaintMode = iota
	Replace
	Add
)

func (g *Grid) Paint(X float64, Y float64, radius float64, value TileType) {
	i := X - math.Ceil(radius)
	j := Y - math.Ceil(radius)

	for i := math.Max(i, 0); i < float64(g.Width) && i < X+math.Ceil(radius)*2; i++ {
		for j := math.Max(j, 0); j < float64(g.Height) && j < Y+math.Ceil(radius)*2; j++ {
			x := X - i
			y := Y - j
			if math.Sqrt(x*x+y*y) <= radius {
				tile := g.Tiles[int(j)][int(i)]
				switch value {
				case ShotRange:
					switch tile.Type {
					case Empty:
						tile.Type = value
					case ShotRange:
						tile.Type = ShotRange2
					case ShotRange2:
						tile.Type = ShotRange3
					}
				default:
					tile.Type = value
				}
			}
		}
	}
}

func (g *Grid) GetTile(x, y int) *Tile {
	if x < 0 || x >= g.Width || y < 0 || y >= g.Height {
		return nil
	}
	return g.Tiles[y][x]
}

func (g *Grid) String() string {
	mem := make([]byte, 0, g.Width*g.Height+g.Height+g.Width)
	buf := bytes.NewBuffer(mem)

	for _, row := range g.Tiles {
		for _, col := range row {
			buf.WriteString(col.Type.String())
		}
		buf.WriteRune('\n')
	}

	data, _ := ioutil.ReadAll(buf)
	return string(data)
}

type TileType int

const (
	Empty      TileType = 0
	Walked     TileType = 1
	ShotRange  TileType = 25
	ShotRange2 TileType = 50
	ShotRange3 TileType = 75
	Ship       TileType = 1000
	Blocked    TileType = 1000000
)

func (t TileType) String() string {
	repr := map[TileType]string{
		Empty:     "O",
		Walked:    "*",
		ShotRange: "#",
		ShotRange2: "%",
		ShotRange3: "@",
		Ship:      "V",
		Blocked:   "X",
	}
	return repr[t]
}

type Tile struct {
	Type TileType
	X    int
	Y    int
	Grid *Grid
}

func (c *Tile) String() string {
	return fmt.Sprintf("x:%d y:%d", c.X, c.Y)
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
	return math.Sqrt(math.Pow(float64(t.X-toT.X), 2)+math.Pow(float64(t.Y-toT.Y), 2)) * (float64(toT.Type) + 1)
}

// PathEstimatedCost is a heuristic method for estimating movement costs
// between non-adjacent nodes.
func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	return math.Sqrt(math.Pow(float64(t.X-toT.X), 2) + math.Pow(float64(t.Y-toT.Y), 2))
}
