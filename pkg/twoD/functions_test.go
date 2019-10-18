package twoD_test

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/metalblueberry/halite-bot/pkg/twoD"
)

//type SimpleEntity struct {
//X, Y, R float64
//Id      int
//}

//func (o *SimpleEntity) Circle() (x, y, r float64) {
//return o.X, o.Y, o.R
//}

//func (o *SimpleEntity) Position() (x, y float64) {
//return o.X, o.Y
//}

//func (o *SimpleEntity) ID() int {
//return o.Id
//}

var _ = Describe("Gamemap", func() {
	Describe("Testing LineMNFrom", func() {
		It("Should return the right values", func() {
			m, n := LineMNFrom(
				NewPosition(1, 3),
				NewPosition(2, 6),
			)
			Expect(m).To(BeNumerically("==", 3))
			Expect(n).To(BeNumerically("==", 0))

		})
		It("Should handle vertical lines", func() {
			m, n := LineMNFrom(
				NewPosition(0, 0),
				NewPosition(0, 6),
			)
			Expect(math.IsInf(m, 0)).To(BeTrue())
			Expect(n).To(BeNumerically("==", 0))
		})
		It("Should handle vertical lines not passing from origin", func() {
			m, n := LineMNFrom(
				NewPosition(3, 0),
				NewPosition(3, 6),
			)
			Expect(math.IsInf(m, 0)).To(BeTrue())
			Expect(n).To(BeNumerically("==", 3))
		})
		It("Should handle horizontal lines", func() {
			m, n := LineMNFrom(
				NewPosition(0, 0),
				NewPosition(6, 0),
			)
			Expect(m).To(BeNumerically("==", 0))
			Expect(n).To(BeNumerically("==", 0))
		})
		It("Should handle horizontal lines not passing from origin", func() {
			m, n := LineMNFrom(
				NewPosition(0, 3),
				NewPosition(6, 3),
			)
			Expect(m).To(BeNumerically("==", 0))
			Expect(n).To(BeNumerically("==", 3))
		})
	})
	Describe("Testing DistancePointToLine", func() {
		It("Should handle horizontal lines", func() {
			dist := DistancePointToLine(
				NewPosition(0, 0),
				NewPosition(6, 0),
				NewPosition(0, 2),
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Should handle vertical lines", func() {
			dist := DistancePointToLine(
				NewPosition(0, 6),
				NewPosition(0, 0),
				NewPosition(2, 0),
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Should handle 45 line", func() {
			dist := DistancePointToLine(
				NewPosition(0, 0),
				NewPosition(3, 3),
				NewPosition(3, -3),
			)
			Expect(dist).To(BeNumerically("~", math.Sqrt(3*3+3*3), 0.001))
		})
		It("Should handle 45 line", func() {
			dist := DistancePointToLine(
				NewPosition(0, 0),
				NewPosition(3, 3),
				NewPosition(-3, 3),
			)
			Expect(dist).To(BeNumerically("~", math.Sqrt(3*3+3*3), 0.001))
		})
		It("Should handle inline points", func() {
			dist := DistancePointToLine(
				NewPosition(4, 4),
				NewPosition(3, 3),
				NewPosition(8, 8),
			)
			Expect(dist).To(BeNumerically("~", 0, 0.001))
		})
		It("Should handle inline points", func() {
			dist := DistancePointToLine(
				NewPosition(00, 10),
				NewPosition(20, 10),
				NewPosition(10, 10),
			)
			Expect(dist).To(BeNumerically("~", 0, 0.001))
		})
	})
	Describe("Testing Distance", func() {
		It("Shoud get the distance of positive X vector", func() {
			dist := Distance(
				NewPosition(0, 0),
				NewPosition(2, 0),
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance of negative X vector", func() {
			dist := Distance(
				NewPosition(0, 0),
				NewPosition(-2, 0),
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance of positive X vector", func() {
			dist := Distance(
				NewPosition(0, 0),
				NewPosition(0, -2),
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance of negative Y vector", func() {
			dist := Distance(
				NewPosition(0, 0),
				NewPosition(0, -2),
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance outside the origin", func() {
			dist := Distance(
				NewPosition(2, 0),
				NewPosition(0, 2),
			)
			Expect(dist).To(BeNumerically("==", math.Sqrt(2*2+2*2)))
		})
	})
	Describe("Testing UnitVector", func() {
		It("Should get x unit vectors", func() {
			x, y := UnitVector(
				NewPosition(0, 0),
				NewPosition(1, 0),
			)
			Expect(x).To(BeNumerically("==", 1))
			Expect(y).To(BeNumerically("==", 0))
		})
		It("Should get y unit vectors", func() {
			x, y := UnitVector(
				NewPosition(0, 0),
				NewPosition(0, 1),
			)
			Expect(x).To(BeNumerically("==", 0))
			Expect(y).To(BeNumerically("==", 1))
		})
		It("Should get negative x unit vectors", func() {
			x, y := UnitVector(
				NewPosition(0, 0),
				NewPosition(-1, 0),
			)
			Expect(x).To(BeNumerically("==", -1))
			Expect(y).To(BeNumerically("==", 0))
		})
		It("Should get negative y unit vectors", func() {
			x, y := UnitVector(
				NewPosition(0, 0),
				NewPosition(0, -1),
			)
			Expect(x).To(BeNumerically("==", 0))
			Expect(y).To(BeNumerically("==", -1))
		})
		It("Should get non 0 vectors", func() {
			x, y := UnitVector(
				NewPosition(1, 1),
				NewPosition(1, -1),
			)
			Expect(x).To(BeNumerically("==", 0))
			Expect(y).To(BeNumerically("==", -1))
		})
		It("Should get non 0 vectors", func() {
			x, y := UnitVector(
				NewPosition(2, 0),
				NewPosition(-1, 0),
			)
			Expect(x).To(BeNumerically("==", -1))
			Expect(y).To(BeNumerically("==", 0))
		})
	})
	Describe("Testing projection", func() {
		It("Should project over X axis", func() {
			projection := Project(
				NewPosition(0, 0),
				NewPosition(1, 0),
				NewPosition(2, 0),
			)
			Expect(projection).To(BeNumerically("==", 2))
		})
		It("Should project over Y axis", func() {
			projection := Project(
				NewPosition(0, 0),
				NewPosition(0, 1),
				NewPosition(2, 0),
			)
			Expect(projection).To(BeNumerically("==", 0))
		})
		It("Should project over 45degree axis", func() {
			projection := Project(
				NewPosition(0, 0),
				NewPosition(1, 1),
				NewPosition(2, 0),
			)
			Expect(projection).To(BeNumerically("~", math.Sqrt(2), 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				NewPosition(1, 0),
				NewPosition(2, 0),
				NewPosition(2, 6),
			)
			Expect(projection).To(BeNumerically("~", 1, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				NewPosition(1, 0),
				NewPosition(2, 0),
				NewPosition(0, 1),
			)
			Expect(projection).To(BeNumerically("~", -1, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				NewPosition(2, 0),
				NewPosition(1, 0),
				NewPosition(0, 1),
			)
			Expect(projection).To(BeNumerically("~", 2, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				NewPosition(5, 2),
				NewPosition(3, 2),
				NewPosition(2, 6),
			)
			Expect(projection).To(BeNumerically("~", 3, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				NewPosition(2, 5),
				NewPosition(2, 3),
				NewPosition(6, 2),
			)
			Expect(projection).To(BeNumerically("~", 3, 0.001))
		})
		It("Should project negative distances", func() {
			projection := Project(
				NewPosition(1, 1),
				NewPosition(2, 2),
				NewPosition(-0, -0),
			)
			Expect(projection).To(BeNumerically("~", -math.Sqrt(2), 0.001))
		})
	})
})
