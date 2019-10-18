package twoD

import "math"

func CalculateRadAngleTo(from, to Positioner) float64 {
	tox, toy := to.Position()
	fromx, fromy := from.Position()
	dx := tox - fromx
	dy := toy - fromy

	return math.Atan2(dy, dx)
}

func CalculateAngleTo(from, to Positioner) float64 {
	return RadToDeg(CalculateRadAngleTo(from, to))
}
func Distance(A, B Positioner) float64 {
	Ax, Ay := A.Position()
	Bx, By := B.Position()
	return math.Sqrt(math.Pow(Ax-Bx, 2) + math.Pow(Ay-By, 2))
}

func UnitVector(A, B Positioner) (x, y float64) {
	Ax, Ay := A.Position()
	Bx, By := B.Position()
	mod := Distance(A, B)
	return (Bx - Ax) / mod, (By - Ay) / mod
}

// VectorTo returns the distance and direction between two entities.
func VectorTo(from, to Positioner) Positioner {
	x, y := to.Position()
	fromx, fromy := to.Position()
	return NewPosition(x-fromx, y-fromy)
}

// Rotate returns the position rotated an angle from the point origin
func Rotate(point, origin Positioner, angle float64) Positioner {
	// x2=r−u=cosβx1−sinβy1y2=t+s=sinβx1+cosβy1
	ox, oy := origin.Position()
	px, py := point.Position()
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return NewPosition(
		cos*px-sin*py+ox,
		sin*px+cos*py+oy,
	)
}

// DegToRad converts degrees to radians
func DegToRad(d float64) float64 {
	return d / 180 * math.Pi
}

// RadToDeg converts radians to degrees
func RadToDeg(r float64) float64 {
	return r / math.Pi * 180
}
