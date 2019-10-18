package twoD

import "math"

// LineMNFrom solves line equation for y = mx+n
// if m is infinity, n will be x = n
func LineMNFrom(A, B Positioner) (m, n float64) {
	Ax, Ay := A.Position()
	Bx, By := B.Position()
	m = (By - Ay) / (Bx - Ax)
	if math.IsInf(m, 0) {
		return m, Bx
	}
	n = Ay - Ax*m
	//n = (Ay*Bx - Ax*By) / (Bx - Ax)
	return m, n
}

// ClosestPointTo returns the closest point that is at least minDistance from the target
func ClosestPointTo(from, target Circler, minDitance float64) Positioner {
	tx, ty, radius := target.Circle()
	fromx, fromy := from.Position()
	dist := radius + minDitance
	norm := Distance(target, from)
	x := dist*(fromx-tx)/norm + tx
	y := dist*(fromy-ty)/norm + ty
	return NewPosition(x, y)
}

func DistancePointToLine(A, B, P Positioner) float64 {
	Px, Py := P.Position()
	m, n := LineMNFrom(A, B)
	if math.IsInf(m, 0) {
		return math.Abs(Px - n)
	}
	return math.Abs(Px*m+n-Py) / math.Sqrt(m*m+1)
}

func Project(A, B, P Positioner) float64 {
	Px, Py := P.Position()
	Ax, Ay := A.Position()
	Ux, Uy := UnitVector(A, B)
	return (Px-Ax)*Ux + (Py-Ay)*Uy
}
