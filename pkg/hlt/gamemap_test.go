package hlt_test

import (
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
	Describe("Testing ObstaclesBetween", func() {
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
				obstacles = append(obstacles, origin, target)
				collides, collider := ObstaclesBetween2(origin, target, obstacles, origin.ID(), target.ID())
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
				obstacles = append(obstacles, origin, target)
				collides, collider := ObstaclesBetween2(origin, target, obstacles, origin.ID(), target.ID())
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
				obstacles = append(obstacles, origin, target)
				collides, collider := ObstaclesBetween2(origin, target, obstacles, origin.ID(), target.ID())
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
				obstacles = append(obstacles, origin, target)
				collides, collider := ObstaclesBetween2(origin, target, obstacles, origin.ID(), target.ID())
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
				obstacles = append(obstacles, origin, target)
				collides, collider := ObstaclesBetween2(origin, target, obstacles, origin.ID(), target.ID())
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

				obstacles = append(obstacles, origin)
				collides, collider := ObstaclesBetween2(origin, target, obstacles, origin.ID(), target.ID())
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
				obstacles = append(obstacles, origin, target)
				collides, collider := ObstaclesBetween2(origin, target, obstacles, origin.ID(), target.ID())
				Expect(collides).To(BeTrue())
				Expect(collider).To(Equal(obstacles[1]))
			})
		})
	})
})
