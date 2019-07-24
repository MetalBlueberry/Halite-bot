package navigation_test

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	navigation "github.com/metalblueberry/halite-bot/pkg/navigation"
)

var re = regexp.MustCompile(`\s+`)

func generate(data string) string {
	return strings.TrimSuffix(re.ReplaceAllString(data, "\n"), "\n")
}

var _ = Describe("Grid", func() {
	var ()
	BeforeEach(func() {
	})
	Describe("When Initialized", func() {
		It("Should be allocated with desired Tiles", func() {
			grid := navigation.NewGrid(40, 50)
			memory := len(grid.Tiles)
			Expect(memory).To(BeNumerically("==", grid.Width*grid.Height))
		})
		Specify("Grid should be referenced by tiles", func() {
			grid := navigation.NewGrid(5, 10)
			for _, tile := range grid.Tiles {
				Expect(tile.Grid).To(Equal(grid))
			}
		})
		Specify("Tiles must store their position", func() {
			grid := navigation.NewGrid(5, 10)
			for index, tile := range grid.Tiles {
				fmt.Fprintf(GinkgoWriter, "Step %d\n", index)
				By(fmt.Sprintf("Step %d\n", index))
				Expect(tile.X).To(BeNumerically("==", index%grid.Width))
				Expect(tile.Y).To(BeNumerically("==", math.Floor(float64(index)/float64(grid.Width))))
			}
		})
		Specify("GetTile should return values for all range", func() {
			X := 5
			Y := 5
			grid := navigation.NewGrid(X, Y)
			for x := 0; x < X; x++ {
				for y := 0; y < Y; y++ {
					tile := grid.GetTile(float64(x), float64(y))
					Expect(tile).ToNot(BeNil())
				}
			}
		})
		Specify("GetTile should return the value set by SetTile", func() {
			X := 5
			Y := 5
			grid := navigation.NewGrid(X, Y)
			testTile := &navigation.Tile{}
			for x := 0; x < X; x++ {
				for y := 0; y < Y; y++ {
					grid.SetTile(float64(x), float64(y), testTile)
					tile := grid.GetTile(float64(x), float64(y))
					Expect(tile).To(Equal(testTile))
				}
			}
		})
	})
	Describe("When printed as string", func() {
		It("Should return an empty ASCII map", func() {
			grid := navigation.NewGrid(4, 2)
			expected := generate(`
			OOOO
			OOOO
			`)
			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
		It("Should print closed as X", func() {
			grid := navigation.NewGrid(9, 5)
			grid.PaintPlanet(4, 2, 1)
			expected := generate(`
			OOOO+OOOO
			OOO+X+OOO
			OO+XXX+OO
			OOO+X+OOO
			OOOO+OOOO
			`)
			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
		It("Should handle borders", func() {
			grid := navigation.NewGrid(9, 5)
			grid.PaintPlanet(0, 0, 2)
			expected := generate(`
			XXX+OOOOO
			XX+OOOOOO
			X++OOOOOO
			+OOOOOOOO
			OOOOOOOOO
			`)
			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
		It("Should handle borders", func() {
			grid := navigation.NewGrid(9, 5)
			grid.PaintPlanet(8, 4, 2)
			expected := generate(`
			OOOOOOOOO
			OOOOOOOO+
			OOOOOO++X
			OOOOOO+XX
			OOOOO+XXX
			`)
			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
	})
	Describe("When finding paths", func() {
		It("Should find string line in clear path", func() {
			grid := navigation.NewGrid(10, 3)
			start := grid.GetTile(0, 1)
			end := grid.GetTile(9, 1)
			path, distance, found, _ := grid.Path(start, end, 10)

			for _, step := range path {
				step.Type = navigation.Walked
			}

			expected := generate(`
			OOOOOOOOOO
			**********
			OOOOOOOOOO
			`)

			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))

		})
		It("Should return parcial path if cannot reach destiny", func() {
			grid := navigation.NewGrid(20, 3)
			start := grid.GetTile(0, 1)
			end := grid.GetTile(19, 1)
			path, distance, found, _ := grid.Path(start, end, 10)

			for _, step := range path {
				step.Type = navigation.Walked
			}

			expected := generate(`
			OOOOOOOOOOOOOOOOOOOO
			***********OOOOOOOOO
			OOOOOOOOOOOOOOOOOOOO
			`)

			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeFalse())
			Expect(distance).To(BeNumerically(">", 0))

		})
		It("Should avoid obstacles", func() {
			grid := navigation.NewGrid(11, 7)
			grid.PaintPlanet(5, 4, 3)
			start := grid.GetTile(0, 3)
			end := grid.GetTile(10, 3)
			path, distance, found, _ := grid.Path(start, end, 200)

			for _, step := range path {
				step.Type = navigation.Walked
			}

			expected := generate(`
			OOO*****OOO
			OO*++X++*OO
			O*+XXXXX+*O
			*O+XXXXX+O*
			O+XXXXXXX+O
			OO+XXXXX+OO
			OO+XXXXX+OO
			`)

			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))

		})
		It("Should avoid Ships", func() {
			grid := navigation.NewGrid(19, 15)
			grid.PaintShip(4, 7, 5)
			grid.PaintShip(14, 7, 5)

			start := grid.GetTile(9, 0)
			end := grid.GetTile(9, 14)
			path, distance, found, _ := grid.Path(start, end, 200)

			for _, step := range path {
				step.Type = navigation.Walked
			}

			expected := generate(`
			OOOOOOOOO*OOOOOOOOO
			OOOOOOOOO*OOOOOOOOO
			OOOO#OOOO*OOOO#OOOO
			O#######O*O#######O
			#########*#########
			#########*#########
			###VVV###*###VVV###
			###VVV##*%###VVV###
			###VVV###*###VVV###
			#########*#########
			#########*#########
			O#######O*O#######O
			OOOO#OOOO*OOOO#OOOO
			OOOOOOOOO*OOOOOOOOO
			OOOOOOOOO*OOOOOOOOO
			`)

			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))

		})
		It("Should return the best posible path", func() {
			grid := navigation.NewGrid(19, 15)
			grid.PaintShip(4, 7, 5)
			grid.PaintShip(14, 7, 5)

			start := grid.GetTile(9, 0)
			end := grid.GetTile(9, 14)
			path, distance, found, bestPath := grid.Path(start, end, 40)

			for _, step := range path {
				step.Type = navigation.Walked
			}
			for _, step := range bestPath {
				step.Type = navigation.ShotRange3
			}

			expected := generate(`
			OOOOOOOOO@****OOOOO
			OOOOOOOOO@OOOOOOOOO
			OOOO#OOOO@OOOO#OOOO
			O#######O@O#######O
			#########@#########
			#########@#########
			###VVV###@###VVV###
			###VVV###@###VVV###
			###VVV###O###VVV###
			#########O#########
			#########O#########
			O#######OOO#######O
			OOOO#OOOOOOOOO#OOOO
			OOOOOOOOOOOOOOOOOOO
			OOOOOOOOOOOOOOOOOOO
			`)

			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeFalse())
			Expect(distance).To(BeNumerically(">", 0))

		})
	})
	Describe("When finding the dirction", func() {
		It("Should return a point at desired distance in horizontal", func() {
			grid := navigation.NewGrid(10, 5)
			start := grid.GetTile(1, 2)
			end := grid.GetTile(8, 2)
			path, distance, found, _ := grid.Path(start, end, 200)
			destiny := navigation.GetDirectionFromPath(path, 5)

			for _, step := range path {
				step.Type = navigation.Walked
			}
			destiny.Type = navigation.Ship

			expected := generate(`
			OOOOOOOOOO
			OOOOOOOOOO
			O*****V**O
			OOOOOOOOOO
			OOOOOOOOOO
			`)

			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))
		})
		It("Should return a point at desired distance in diagonal", func() {
			grid := navigation.NewGrid(10, 5)
			start := grid.GetTile(9, 4)
			end := grid.GetTile(0, 0)
			path, distance, found, _ := grid.Path(start, end, 200)
			destiny := navigation.GetDirectionFromPath(path, 5)

			for _, step := range path {
				step.Type = navigation.Walked
			}
			destiny.Type = navigation.Ship

			expected := generate(`
			*OOOOOOOOO
			O***OOOOOO
			OOOO*VOOOO
			OOOOOO**OO
			OOOOOOOO**
			`)

			result := grid.String()
			fmt.Fprintln(GinkgoWriter, "")
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))
		})
	})
})
