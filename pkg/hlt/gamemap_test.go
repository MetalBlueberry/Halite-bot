package hlt_test

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/metalblueberry/halite-bot/pkg/hlt"
)

type SimpleEntity struct {
	Entity
	X, Y, R float64
	Id      int
}

func (o *SimpleEntity) Circle() (x, y, r float64) {
	return o.X, o.Y, o.R
}

func (o *SimpleEntity) Position() (x, y float64) {
	return o.X, o.Y
}

func (o *SimpleEntity) ID() int {
	return o.Id
}

var _ = Describe("Gamemap", func() {
	Describe("Testing LineMNFrom", func() {
		It("Should return the right values", func() {
			m, n := LineMNFrom(
				&SimpleEntity{
					X: 1,
					Y: 3,
				},
				&SimpleEntity{
					X: 2,
					Y: 6,
				},
			)
			Expect(m).To(BeNumerically("==", 3))
			Expect(n).To(BeNumerically("==", 0))

		})
		It("Should handle vertical lines", func() {
			m, n := LineMNFrom(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: 6,
				},
			)
			Expect(math.IsInf(m, 0)).To(BeTrue())
			Expect(n).To(BeNumerically("==", 0))
		})
		It("Should handle vertical lines not passing from origin", func() {
			m, n := LineMNFrom(
				&SimpleEntity{
					X: 3,
					Y: 0,
				},
				&SimpleEntity{
					X: 3,
					Y: 6,
				},
			)
			Expect(math.IsInf(m, 0)).To(BeTrue())
			Expect(n).To(BeNumerically("==", 3))
		})
		It("Should handle horizontal lines", func() {
			m, n := LineMNFrom(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 6,
					Y: 0,
				},
			)
			Expect(m).To(BeNumerically("==", 0))
			Expect(n).To(BeNumerically("==", 0))
		})
		It("Should handle horizontal lines not passing from origin", func() {
			m, n := LineMNFrom(
				&SimpleEntity{
					X: 0,
					Y: 3,
				},
				&SimpleEntity{
					X: 6,
					Y: 3,
				},
			)
			Expect(m).To(BeNumerically("==", 0))
			Expect(n).To(BeNumerically("==", 3))
		})
	})
	Describe("Testing DistancePointToLine", func() {
		It("Should handle horizontal lines", func() {
			dist := DistancePointToLine(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 6,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: 2,
				},
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Should handle vertical lines", func() {
			dist := DistancePointToLine(
				&SimpleEntity{
					X: 0,
					Y: 6,
				},
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Should handle 45 line", func() {
			dist := DistancePointToLine(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 3,
					Y: 3,
				},
				&SimpleEntity{
					X: 3,
					Y: -3,
				},
			)
			Expect(dist).To(BeNumerically("~", math.Sqrt(3*3+3*3), 0.001))
		})
		It("Should handle 45 line", func() {
			dist := DistancePointToLine(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 3,
					Y: 3,
				},
				&SimpleEntity{
					X: -3,
					Y: 3,
				},
			)
			Expect(dist).To(BeNumerically("~", math.Sqrt(3*3+3*3), 0.001))
		})
		It("Should handle inline points", func() {
			dist := DistancePointToLine(
				&SimpleEntity{
					X: 4,
					Y: 4,
				},
				&SimpleEntity{
					X: 3,
					Y: 3,
				},
				&SimpleEntity{
					X: 8,
					Y: 8,
				},
			)
			Expect(dist).To(BeNumerically("~", 0, 0.001))
		})
		It("Should handle inline points", func() {
			dist := DistancePointToLine(
				&SimpleEntity{
					X: 00,
					Y: 10,
				},
				&SimpleEntity{
					X: 20,
					Y: 10,
				},
				&SimpleEntity{
					X: 10,
					Y: 10,
				},
			)
			Expect(dist).To(BeNumerically("~", 0, 0.001))
		})
	})
	Describe("Testing Distance", func() {
		It("Shoud get the distance of positive X vector", func() {
			dist := Distance(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance of negative X vector", func() {
			dist := Distance(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: -2,
					Y: 0,
				},
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance of positive X vector", func() {
			dist := Distance(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: -2,
				},
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance of negative Y vector", func() {
			dist := Distance(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: -2,
				},
			)
			Expect(dist).To(BeNumerically("==", 2))
		})
		It("Shoud get the distance outside the origin", func() {
			dist := Distance(
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: 2,
				},
			)
			Expect(dist).To(BeNumerically("==", math.Sqrt(2*2+2*2)))
		})
	})
	Describe("Testing UnitVector", func() {
		It("Should get x unit vectors", func() {
			x, y := UnitVector(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 1,
					Y: 0,
				},
			)
			Expect(x).To(BeNumerically("==", 1))
			Expect(y).To(BeNumerically("==", 0))
		})
		It("Should get y unit vectors", func() {
			x, y := UnitVector(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: 1,
				},
			)
			Expect(x).To(BeNumerically("==", 0))
			Expect(y).To(BeNumerically("==", 1))
		})
		It("Should get negative x unit vectors", func() {
			x, y := UnitVector(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: -1,
					Y: 0,
				},
			)
			Expect(x).To(BeNumerically("==", -1))
			Expect(y).To(BeNumerically("==", 0))
		})
		It("Should get negative y unit vectors", func() {
			x, y := UnitVector(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: -1,
				},
			)
			Expect(x).To(BeNumerically("==", 0))
			Expect(y).To(BeNumerically("==", -1))
		})
		It("Should get non 0 vectors", func() {
			x, y := UnitVector(
				&SimpleEntity{
					X: 1,
					Y: 1,
				},
				&SimpleEntity{
					X: 1,
					Y: -1,
				},
			)
			Expect(x).To(BeNumerically("==", 0))
			Expect(y).To(BeNumerically("==", -1))
		})
		It("Should get non 0 vectors", func() {
			x, y := UnitVector(
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
				&SimpleEntity{
					X: -1,
					Y: 0,
				},
			)
			Expect(x).To(BeNumerically("==", -1))
			Expect(y).To(BeNumerically("==", 0))
		})
	})
	Describe("Testing projection", func() {
		It("Should project over X axis", func() {
			projection := Project(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 1,
					Y: 0,
				},
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
			)
			Expect(projection).To(BeNumerically("==", 2))
		})
		It("Should project over Y axis", func() {
			projection := Project(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: 1,
				},
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
			)
			Expect(projection).To(BeNumerically("==", 0))
		})
		It("Should project over 45degree axis", func() {
			projection := Project(
				&SimpleEntity{
					X: 0,
					Y: 0,
				},
				&SimpleEntity{
					X: 1,
					Y: 1,
				},
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
			)
			Expect(projection).To(BeNumerically("~", math.Sqrt(2), 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				&SimpleEntity{
					X: 1,
					Y: 0,
				},
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
				&SimpleEntity{
					X: 2,
					Y: 6,
				},
			)
			Expect(projection).To(BeNumerically("~", 1, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				&SimpleEntity{
					X: 1,
					Y: 0,
				},
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: 1,
				},
			)
			Expect(projection).To(BeNumerically("~", -1, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				&SimpleEntity{
					X: 2,
					Y: 0,
				},
				&SimpleEntity{
					X: 1,
					Y: 0,
				},
				&SimpleEntity{
					X: 0,
					Y: 1,
				},
			)
			Expect(projection).To(BeNumerically("~", 2, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				&SimpleEntity{
					X: 5,
					Y: 2,
				},
				&SimpleEntity{
					X: 3,
					Y: 2,
				},
				&SimpleEntity{
					X: 2,
					Y: 6,
				},
			)
			Expect(projection).To(BeNumerically("~", 3, 0.001))
		})
		It("Should project to a non origin vector", func() {
			projection := Project(
				&SimpleEntity{
					X: 2,
					Y: 5,
				},
				&SimpleEntity{
					X: 2,
					Y: 3,
				},
				&SimpleEntity{
					X: 6,
					Y: 2,
				},
			)
			Expect(projection).To(BeNumerically("~", 3, 0.001))
		})
		It("Should project negative distances", func() {
			projection := Project(
				&SimpleEntity{
					X: 1,
					Y: 1,
				},
				&SimpleEntity{
					X: 2,
					Y: 2,
				},
				&SimpleEntity{
					X: -0,
					Y: -0,
				},
			)
			Expect(projection).To(BeNumerically("~", -math.Sqrt(2), 0.001))
		})
	})
	FDescribe("Testing ObstaclesBetween", func() {
		var (
			obstacles []Entitier
		)

		Describe("Avoid a single planet", func() {
			BeforeEach(func() {
				obstacles = []Entitier{
					&SimpleEntity{
						X:  10,
						Y:  10,
						R:  2,
						Id: 1,
					},
				}
			})
			It("Should detect collision passing through the origin", func() {
				origin := &SimpleEntity{
					X:  00,
					Y:  10,
					R:  1,
					Id: 0,
				}
				target := &SimpleEntity{
					X:  20,
					Y:  10,
					R:  1,
					Id: -1,
				}
				collides, collider := ObstaclesBetween2(origin, target, obstacles)
				Expect(collides).To(BeTrue())
				Expect(collider).To(Equal(obstacles[0]))
			})
			It("Should detect collision by side", func() {
				origin := &SimpleEntity{
					X:  00,
					Y:  11,
					R:  1,
					Id: 0,
				}
				target := &SimpleEntity{
					X:  20,
					Y:  11,
					R:  1,
					Id: -1,
				}
				collides, collider := ObstaclesBetween2(origin, target, obstacles)
				Expect(collides).To(BeTrue())
				Expect(collider).To(Equal(obstacles[0]))
			})
			It("Should detect collision 45degees", func() {
				origin := &SimpleEntity{
					X:  0,
					Y:  0,
					R:  1,
					Id: 0,
				}
				target := &SimpleEntity{
					X:  20,
					Y:  20,
					R:  1,
					Id: -1,
				}
				collides, collider := ObstaclesBetween2(origin, target, obstacles)
				Expect(collides).To(BeTrue())
				Expect(collider).To(Equal(obstacles[0]))
			})
			It("Should allow with 1 unit margin", func() {
				origin := &SimpleEntity{
					X:  00,
					Y:  14,
					R:  1,
					Id: 0,
				}
				target := &SimpleEntity{
					X:  20,
					Y:  14,
					R:  1,
					Id: -1,
				}
				collides, collider := ObstaclesBetween2(origin, target, obstacles)
				Expect(collides).To(BeFalse())
				Expect(collider).To(BeNil())
			})
		})
		Describe("Avoid a multiple obstacles", func() {
			BeforeEach(func() {
				obstacles = []Entitier{
					&SimpleEntity{
						X:  2,
						Y:  10,
						R:  2,
						Id: 1,
					},
					&SimpleEntity{
						X:  8,
						Y:  10,
						R:  2,
						Id: 1,
					},
				}
			})
			It("Should detect collision from left to rigth", func() {
				origin := &SimpleEntity{
					X:  0,
					Y:  12,
					R:  1,
					Id: 0,
				}
				target := &SimpleEntity{
					X:  10,
					Y:  8,
					R:  1,
					Id: -1,
				}
				collides, collider := ObstaclesBetween2(origin, target, obstacles)
				Expect(collides).To(BeTrue())
				Expect(collider).To(Equal(obstacles[0]))
			})
			It("Should detect collision from right to left", func() {
				origin := &SimpleEntity{
					X:  10,
					Y:  8,
					R:  1,
					Id: -1,
				}
				target := &SimpleEntity{
					X:  0,
					Y:  12,
					R:  1,
					Id: 0,
				}
				collides, collider := ObstaclesBetween2(origin, target, obstacles)
				Expect(collides).To(BeTrue())
				Expect(collider).To(Equal(obstacles[0]))
			})
			It("Should detect collision only with objects in path", func() {
				origin := &SimpleEntity{
					X:  10,
					Y:  10,
					R:  1,
					Id: -1,
				}
				target := &SimpleEntity{
					X:  5,
					Y:  10,
					R:  1,
					Id: 0,
				}
				collides, collider := ObstaclesBetween2(origin, target, obstacles)
				Expect(collides).To(BeTrue())
				Expect(collider).To(Equal(obstacles[1]))
			})
		})
	})
})
