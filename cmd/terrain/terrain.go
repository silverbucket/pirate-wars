package terrain

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/ojrac/opensimplex-go"
	"math"
	"math/rand"
	"pirate-wars/cmd/avatar"
)

// Icon ideas
// Towns: ⩎
// Boats: ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
// People: 옷

const (
	WorldWidth       = 600
	WorldHeight      = 600
	ViewWidth        = 75
	ViewHeight       = 50
	BufferWidth      = 21
	BufferHeight     = 26
	MiniMapFactor    = 11
	TypeDeepWater    = 0
	TypeOpenWater    = 1
	TypeShallowWater = 2
	TypeBeach        = 3
	TypeLowland      = 4
	TypeHighland     = 5
	TypeRock         = 6
	TypePeak         = 7
)

type ViewPort struct {
	width   int
	height  int
	topLeft int
}

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
	TypeDeepWater:    {symbol: '⏖', style: createTerrainItem("18"), Passable: true, RequiresBoat: true},
	TypeOpenWater:    {symbol: '⏝', style: createTerrainItem("20"), Passable: true, RequiresBoat: true},
	TypeShallowWater: {symbol: '⏑', style: createTerrainItem("26"), Passable: true, RequiresBoat: true},
	TypeBeach:        {symbol: '~', style: createTerrainItem("#dad1ad"), Passable: true, RequiresBoat: false},
	TypeLowland:      {symbol: ':', style: createTerrainItem("113"), Passable: true, RequiresBoat: false},
	TypeHighland:     {symbol: ':', style: createTerrainItem("142"), Passable: true, RequiresBoat: false},
	TypeRock:         {symbol: '%', style: createTerrainItem("244"), Passable: true, RequiresBoat: false},
	TypePeak:         {symbol: '^', style: createTerrainItem("15"), Passable: false, RequiresBoat: false},
}

type World [][]Type

func createTerrainItem(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Background(color).Padding(0).Margin(0)
}

func Init() *Terrain {
	//default values for terrain map generation
	t := Terrain{
		width:       WorldWidth,
		height:      WorldHeight,
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

func (t *Terrain) GenerateWorld() World {
	//var world [WorldWidth][WorldHeight]Type
	world := make([][]Type, WorldHeight)
	for i := range world {
		world[i] = make([]Type, WorldHeight)
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
				world[x][y] = TypeDeepWater
			} else if s > 0.46 {
				world[x][y] = TypeOpenWater
			} else if s > 0.42 {
				world[x][y] = TypeShallowWater
			} else if s > 0.40 {
				world[x][y] = TypeBeach
			} else if s > 0.31 {
				world[x][y] = TypeLowland
			} else if s > 0.26 {
				world[x][y] = TypeHighland
			} else if s > 0.21 {
				world[x][y] = TypeRock
			} else {
				world[x][y] = TypePeak
			}
		}
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
	height := len(world) / MiniMapFactor
	width := len(world[0]) / MiniMapFactor

	// Create new 2D slice
	newArr := make([][]Type, height+1)
	for i := range newArr {
		newArr[i] = make([]Type, width+1)
	}

	// Down-sample
	for i, row := range world {
		for j, val := range row {
			// Calculate corresponding index in new slice
			newI := i / MiniMapFactor
			newJ := j / MiniMapFactor

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
	avatarCoordsX := avatar.GetX()
	avatarCoordsY := avatar.GetY()

	if !isMiniMap {
		left = int(math.Max(float64(avatar.GetX()-ViewWidth+BufferWidth), 0))
		top = int(math.Max(float64(avatar.GetY()-ViewHeight+BufferHeight), 0))
		viewHeight = ViewHeight + top
		viewWidth = ViewWidth + left
		rowWidth = ViewWidth
	} else {
		avatarCoordsX = avatarCoordsX / MiniMapFactor
		avatarCoordsY = avatarCoordsY / MiniMapFactor
	}

	viewport := table.New().BorderBottom(false).BorderTop(false).BorderLeft(false).BorderRight(false)

	for y := top; y < worldHeight && y < viewHeight; y++ {
		var row = make([]string, rowWidth)
		for x := left; x < worldWidth && x < viewWidth; x++ {
			//fmt.Println("[%v %v == %v %v]", x, y, avatarCoordsX, avatarCoordsY)
			if x == avatarCoordsX && y == avatarCoordsY {
				row[x-left] = avatar.Render()
			} else {
				row[x-left] = world[x][y].Render()
			}
		}
		viewport.Row(row...).BorderColumn(false)
	}

	return fmt.Sprintln(viewport)
}
