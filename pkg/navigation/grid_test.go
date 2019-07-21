package navigation_test

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	navigation "github.com/metalblueberry/halite-bot/pkg/navigation"
	"github.com/metalblueberry/halite-bot/pkg/astar"
)

var re = regexp.MustCompile(`\s+`)

func generate(data string) string {
	return strings.TrimPrefix(re.ReplaceAllString(data, "\n"), "\n")
}

var _ = Describe("Grid", func() {
	var ()
	BeforeEach(func() {
	})
	Describe("When Initialized", func() {
		It("Should be allocated with desired Tiles", func() {
			grid := navigation.NewGrid(50, 50)
			width := len(grid.Tiles)
			Expect(width).To(Equal(grid.Width))
			for _, col := range grid.Tiles {
				height := len(col)
				Expect(height).To(Equal(grid.Height))
			}
		})
		Specify("Grid should be referenced by tiles", func() {
			grid := navigation.NewGrid(5, 5)
			for _, row := range grid.Tiles {
				for _, tile := range row {
					Expect(tile.Grid).To(Equal(grid))
				}
			}
		})
		Specify("Tiles must store theri position", func() {
			grid := navigation.NewGrid(5, 5)
			for y, row := range grid.Tiles {
				for x, tile := range row {
					Expect(tile.X).To(Equal(x))
					Expect(tile.Y).To(Equal(y))
				}
			}
		})
	})
	Describe("When printed as string", func() {
		It("Should return an empty ASCII map", func() {
			grid := navigation.NewGrid(10, 5)
			expected := generate(`
			OOOOOOOOOO
			OOOOOOOOOO
			OOOOOOOOOO
			OOOOOOOOOO
			OOOOOOOOOO
			`)
			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
		It("Should print closed as X", func() {
			grid := navigation.NewGrid(9, 5)
			grid.PaintPlanet(4, 2, 1)
			expected := generate(`
			OOOOOOOOO
			OOOOXOOOO
			OOOXXXOOO
			OOOOXOOOO
			OOOOOOOOO
			`)
			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
		It("Should handle borders", func() {
			grid := navigation.NewGrid(9, 5)
			grid.PaintPlanet(0, 0, 2)
			expected := generate(`
			XXXOOOOOO
			XXOOOOOOO
			XOOOOOOOO
			OOOOOOOOO
			OOOOOOOOO
			`)
			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
		It("Should handle borders", func() {
			grid := navigation.NewGrid(9, 5)
			grid.PaintPlanet(8, 4, 2)
			expected := generate(`
			OOOOOOOOO
			OOOOOOOOO
			OOOOOOOOX
			OOOOOOOXX
			OOOOOOXXX
			`)
			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))
		})
	})
	Describe("When finding paths", func() {
		It("Should find string line in clear path", func() {
			grid := navigation.NewGrid(10, 3)
			start := grid.GetTile(0, 1)
			end := grid.GetTile(9, 1)
			path, distance, found := astar.Path(start, end, 10)

			for _, step := range path {
				t := step.(*navigation.Tile)
				t.Type = navigation.Walked
			}

			expected := generate(`
			OOOOOOOOOO
			**********
			OOOOOOOOOO
			`)

			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))

		})
		It("Should return parcial path if cannot reach destiny", func() {
			grid := navigation.NewGrid(20, 3)
			start := grid.GetTile(0, 1)
			end := grid.GetTile(19, 1)
			path, distance, found := astar.Path(start, end, 10)

			for _, step := range path {
				t := step.(*navigation.Tile)
				t.Type = navigation.Walked
			}

			expected := generate(`
			OOOOOOOOOOOOOOOOOOOO
			***********OOOOOOOOO
			OOOOOOOOOOOOOOOOOOOO
			`)

			result := grid.String()
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
			path, distance, found := astar.Path(start, end, 200)

			for _, step := range path {
				t := step.(*navigation.Tile)
				t.Type = navigation.Walked
			}

			expected := generate(`
			OOOOO*OOOOO
			OOO**X**OOO
			OO*XXXXX**O
			**OXXXXXOO*
			OOXXXXXXXOO
			OOOXXXXXOOO
			OOOXXXXXOOO
			`)

			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))

		})
		It("Should avoid Ships", func() {
			grid := navigation.NewGrid(19, 15)
			grid.PaintShip(4,7,5)
			grid.PaintShip(14,7,5)

			start := grid.GetTile(9, 0)
			end := grid.GetTile(9, 14)
			path, distance, found := astar.Path(start, end, 200)

			for _, step := range path {
				t := step.(*navigation.Tile)
				t.Type = navigation.Walked
			}

			expected := generate(`
			OOOOOOOOO*OOOOOOOOO
			OOOOOOOOO*OOOOOOOOO
			OOOO#OOOO*OOOO#OOOO
			O#######O*O#######O
			#########*#########
			#########*#########
			#########*#########
			####V####%*###V####
			#########*#########
			#########*#########
			#########*#########
			O#######O*O#######O
			OOOO#OOOO*OOOO#OOOO
			OOOOOOOOO*OOOOOOOOO
			OOOOOOOOO*OOOOOOOOO
			`)

			result := grid.String()
			fmt.Fprint(GinkgoWriter, result)
			Expect(result).To(Equal(expected))

			Expect(found).To(BeTrue())
			Expect(distance).To(BeNumerically(">", 0))

		})
	})
})
