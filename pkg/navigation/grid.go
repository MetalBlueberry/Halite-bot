package navigation

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"

	"github.com/metalblueberry/halite-bot/pkg/astar"
)

type Grid struct {
	Width, Height int
	Tiles         []*Tile
}

func NewGrid(Width, Height int) *Grid {
	grid := &Grid{
		Width:  Width,
		Height: Height,
		Tiles:  make([]*Tile, Height*Width, Height*Width),
	}
	for index := range grid.Tiles {
		grid.Tiles[index] = &Tile{
			Grid: grid,
			X:    float64(index % Width),
			Y:    math.Floor(float64(index) / float64(Width)),
		}
	}
	return grid
}

func (g *Grid) PaintShip(X float64, Y float64, shotRange float64) {
	g.Paint(X, Y, shotRange, ShotRange)
	g.Paint(X, Y, 1.5, Ship)
}

func (g *Grid) PaintPlanet(X float64, Y float64, radius float64) {
	g.Paint(X, Y, radius+1, SafeMargin)
	g.Paint(X, Y, radius, Blocked)
}

func (g *Grid) Paint(X float64, Y float64, radius float64, value TileType) {
	i := X - math.Ceil(radius)
	j := Y - math.Ceil(radius)

	for i := math.Max(i, 0); i < float64(g.Width) && i < X+math.Ceil(radius)*2; i++ {
		for j := math.Max(j, 0); j < float64(g.Height) && j < Y+math.Ceil(radius)*2; j++ {
			x := X - i
			y := Y - j
			if math.Sqrt(x*x+y*y) <= radius {
				tile := g.GetTileSafe(i, j)
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

func (g *Grid) GetTile(x, y float64) *Tile {
	if x < 0 || x >= float64(g.Width) || y < 0 || y >= float64(g.Height) {
		return nil
	}
	return g.Tiles[int(int(y)*g.Width+int(x))]
}

func (g *Grid) GetTileSafe(x, y float64) *Tile {
	tile := g.GetTile(x, y)
	if tile == nil {
		log.Panicf("Index out of range \nx:%f y:%f\nw:%d h:%d\n", x, y, g.Width, g.Height)
	}
	return tile
}

func (g *Grid) SetTile(x, y float64, tile *Tile) {
	if x < 0 || x >= float64(g.Width) || y < 0 || y >= float64(g.Height) {
		log.Panicf("Index out of range \nx:%f y:%f\nw:%d h:%d\n", x, y, g.Width, g.Height)
	}
	g.Tiles[int(int(y)*g.Width+int(x))] = tile
}

func (g *Grid) Path(from, to *Tile, iterations int) (path []*Tile, distance float64, found bool, bestPath []*Tile) {
	result, distance, found, bestResult := astar.Path(from, to, iterations)
	path = make([]*Tile, len(result), len(result))
	for i, step := range result {
		path[i] = step.(*Tile)
	}
	bestPath = make([]*Tile, len(bestResult), len(bestResult))
	for i, step := range bestResult {
		bestPath[i] = step.(*Tile)
	}
	return path, distance, found, bestPath
}

func (g *Grid) String() string {
	mem := make([]byte, 0, int(g.Width*g.Height+g.Height+g.Width))
	buf := bytes.NewBuffer(mem)

	for index, tile := range g.Tiles {
		if index%g.Width == 0 {
			buf.WriteRune('\n')
		}
		buf.WriteString(tile.Type.String())
	}

	data, _ := ioutil.ReadAll(buf)
	return string(data)
}

func (g *Grid) PrintDebugPath(path []*Tile, from *Tile, to *Tile) string {
	mem := make([]byte, 0, int(g.Width*g.Height+g.Height+g.Width))
	buf := bytes.NewBuffer(mem)

	for index, tile := range g.Tiles {
		if from == tile {
			buf.WriteString("V")
			continue
		}
		if to == tile {
			buf.WriteString("8")
			continue
		}
		isPath := false
		for _, inPath := range path {
			if inPath == tile {
				isPath = true
				buf.WriteString("*")
			}
		}
		if !isPath {
			buf.WriteString(tile.Type.String())
		}
		if index%g.Width == 0 {
			buf.WriteRune('\n')
		}
	}

	data, _ := ioutil.ReadAll(buf)
	return string(data)
}
