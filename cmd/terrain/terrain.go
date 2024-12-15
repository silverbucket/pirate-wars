package terrain

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/ojrac/opensimplex-go"
	"math/rand"
	"pirate-wars/cmd/avatar"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/town"
)

// Icon ideas
// Towns: ⩎
// Boats: ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
// People: 옷

type Terrain struct {
	width       int
	height      int
	scale       float64
	lacunarity  float64
	persistence float64
	octaves     int
}

type Type int

type TypeQualities struct {
	symbol       rune
	style        lipgloss.Style
	Passable     bool
	RequiresBoat bool
}

var TypeLookup = map[Type]TypeQualities{
	common.TypeDeepWater:    {symbol: '⏖', style: createTerrainItem("18"), Passable: true, RequiresBoat: true},
	common.TypeOpenWater:    {symbol: '⏝', style: createTerrainItem("20"), Passable: true, RequiresBoat: true},
	common.TypeShallowWater: {symbol: '⏑', style: createTerrainItem("26"), Passable: true, RequiresBoat: true},
	common.TypeBeach:        {symbol: '~', style: createTerrainItem("#dad1ad"), Passable: true, RequiresBoat: false},
	common.TypeLowland:      {symbol: ':', style: createTerrainItem("113"), Passable: true, RequiresBoat: false},
	common.TypeHighland:     {symbol: ':', style: createTerrainItem("142"), Passable: true, RequiresBoat: false},
	common.TypeRock:         {symbol: '%', style: createTerrainItem("244"), Passable: true, RequiresBoat: false},
	common.TypePeak:         {symbol: '^', style: createTerrainItem("15"), Passable: false, RequiresBoat: false},
	common.TypeTown:         {symbol: '⩎', style: createTerrainItem("0"), Passable: true, RequiresBoat: false},
}

type World [][]Type

func createTerrainItem(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Background(color).Padding(0).Margin(0)
}

func Init() *Terrain {
	//default values for terrain map generation
	t := Terrain{
		width:       common.WorldWidth,
		height:      common.WorldHeight,
		scale:       60,
		lacunarity:  2.0,
		persistence: 0.5,
		octaves:     5,
	}
	return &t
}

func (tt Type) Render() string {
	return fmt.Sprintf(TypeLookup[tt].style.PaddingLeft(1).PaddingRight(1).Render("%c"), TypeLookup[tt].symbol)
}

func (tt Type) IsPassableByBoat() bool {
	if TypeLookup[tt].RequiresBoat {
		return true
	}
	return false
}

func GenerateTowns(world World, count int) {
	for i := 0; i <= count; i++ {
		for {
			coords := common.Coordinates{X: min(rand.Intn(common.WorldWidth-5), 5), Y: min(rand.Intn(common.WorldHeight-5), 5)}
			if coords.X > 1 && coords.Y > 1 &&
				coords.X < common.WorldWidth-1 && coords.Y < common.WorldHeight &&
				world[coords.X][coords.Y] == common.TypeBeach &&
				(world[coords.X+1][coords.Y] == common.TypeShallowWater ||
					world[coords.X][coords.Y+1] == common.TypeShallowWater ||
					world[coords.X-1][coords.Y] == common.TypeShallowWater ||
					world[coords.X][coords.Y-1] == common.TypeShallowWater) {
				town.Create(coords, '⩎')
				break
			}
		}
	}
}

func (t *Terrain) Generate() World {
	//var world [WorldWidth][WorldHeight]Type
	world := make([][]Type, common.WorldHeight)
	for i := range world {
		world[i] = make([]Type, common.WorldHeight)
	}

	noise := opensimplex.New(rand.Int63())

	for x := 0; x < t.width; x++ {
		for y := 0; y < t.height; y++ {
			//sample x and y and apply scale
			xFloat := float64(x) / t.scale
			yFloat := float64(y) / t.scale

			//init values for octave calculation
			frequency := 1.0
			amplitude := 1.0
			normalizeOctaves := 0.0
			total := 0.0

			//octave calculation
			for i := 0; i < t.octaves; i++ {
				total += noise.Eval2(xFloat*frequency, yFloat*frequency) * amplitude
				normalizeOctaves += amplitude
				amplitude *= t.persistence
				frequency *= t.lacunarity
			}

			//normalize to -1 to 1, and then from 0 to 1 (this is for the ability to use grayscale, if using colors could keep from -1 to 1)
			var s = (total/normalizeOctaves + 1) / 2
			if s > 0.60 {
				world[x][y] = common.TypeDeepWater
			} else if s > 0.46 {
				world[x][y] = common.TypeOpenWater
			} else if s > 0.42 {
				world[x][y] = common.TypeShallowWater
			} else if s > 0.40 {
				world[x][y] = common.TypeBeach
			} else if s > 0.31 {
				world[x][y] = common.TypeLowland
			} else if s > 0.26 {
				world[x][y] = common.TypeHighland
			} else if s > 0.21 {
				world[x][y] = common.TypeRock
			} else {
				world[x][y] = common.TypePeak
			}
		}
	}

	GenerateTowns(world, common.TotalTowns)
	for _, o := range town.List {
		world[o.GetX()][o.GetY()] = common.TypeTown
	}

	return world
}

func GetType(i int) (Type, error) {
	for k, _ := range TypeLookup {
		if int(k) == i {
			return k, nil
		}
	}
	return 0, errors.New("invalid type")
}

func (world World) RenderMiniMap() World {
	// Calculate new dimensions
	height := len(world) / common.MiniMapFactor
	width := len(world[0]) / common.MiniMapFactor

	// Create new 2D slice
	newArr := make([][]Type, height+1)
	for i := range newArr {
		newArr[i] = make([]Type, width+1)
	}

	// Down-sample
	for i, row := range world {
		for j, val := range row {
			// Calculate corresponding index in new slice
			newI := i / common.MiniMapFactor
			newJ := j / common.MiniMapFactor

			// Assign original value
			newArr[newI][newJ] = val
		}
	}

	return newArr
}

func (world World) Paint(avatar avatar.Type, isMiniMap bool) string {
	left := 0
	top := 0
	worldHeight := len(world)
	worldWidth := len(world[0])
	viewHeight := worldHeight
	viewWidth := worldWidth
	rowWidth := worldWidth
	avatarX := avatar.GetX()
	avatarY := avatar.GetY()

	if isMiniMap {
		avatarX = avatar.GetMiniMapX()
		avatarY = avatar.GetMiniMapY()
		for _, o := range town.List {
			world[o.GetMiniMapX()][o.GetMiniMapY()] = common.TypeTown
		}
	} else {
		// center viewport on avatar
		left = avatar.GetX() - (common.ViewWidth / 2)
		top = avatar.GetY() - (common.ViewHeight / 2)
		viewHeight = common.ViewHeight + top
		viewWidth = common.ViewWidth + left
		rowWidth = common.ViewWidth
	}

	viewport := table.New().BorderBottom(false).BorderTop(false).BorderLeft(false).BorderRight(false)

	for y := top; y < worldHeight && y < viewHeight; y++ {
		var row = make([]string, rowWidth)
		for x := left; x < worldWidth && x < viewWidth; x++ {
			if x == avatarX && y == avatarY {
				row[x-left] = avatar.Render()
			} else {
				//fmt.Printf("[%v , %v , %v ]\n", x, y, x-left)
				row[x-left] = world[x][y].Render()
			}
		}
		viewport.Row(row...).BorderColumn(false)
	}

	return fmt.Sprintln(viewport)
}
