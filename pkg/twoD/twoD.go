package twoD

type Positioner interface {
	Position() (x, y float64)
}

type Circler interface {
	Positioner
	Circle() (x, y, r float64)
}
type Liner interface {
	Line() (float64, float64, float64, float64)
}

type point struct{ x, y float64 }

func NewPosition(x, y float64) Positioner {
	return point{x: x, y: y}
}
func (p point) Position() (x, y float64) {
	return p.x, p.y
}

type line struct {
	X1, Y1 float64
	X2, Y2 float64
}

func NewLine(a, b Positioner) line {
	X1, Y1 := a.Position()
	X2, Y2 := b.Position()
	return line{X1, Y1, X2, Y2}
}

func (l line) Line() (float64, float64, float64, float64) {
	return l.X1, l.Y1, l.X2, l.Y2
}
