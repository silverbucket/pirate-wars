package terrain

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/ojrac/opensimplex-go"
	"go.uber.org/zap"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/town"
)

// Icon ideas
// Towns: ⩎
// Boats: ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
// People: 옷

type Props struct {
	width       int
	height      int
	scale       float64
	lacunarity  float64
	persistence float64
	octaves     int
}

type Terrain struct {
	logger  *zap.SugaredLogger
	props   Props
	World   MapView
	MiniMap MapView
}

type TerrainType int

type TypeQualities struct {
	symbol       rune
	style        lipgloss.Style
	Passable     bool
	RequiresBoat bool
}

var TypeLookup = map[TerrainType]TypeQualities{
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

type MapView struct {
	grid      [][]TerrainType
	logger    *zap.SugaredLogger
	isMiniMap bool
}

func createTerrainItem(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Background(color).Padding(0).Margin(0)
}

func Init(logger *zap.SugaredLogger) *Terrain {
	logger.Info("Initializing terrain")
	//default values for terrain map generation
	worldGrid := make([][]TerrainType, common.WorldHeight)
	for i := range worldGrid {
		worldGrid[i] = make([]TerrainType, common.WorldHeight)
	}

	// Calculate MiniMap dimensions
	height := len(worldGrid) / common.MiniMapFactor
	width := len(worldGrid[0]) / common.MiniMapFactor

	// Create new 2D slice
	miniMap := make([][]TerrainType, height+1)
	for i := range miniMap {
		miniMap[i] = make([]TerrainType, width+1)
	}

	return &Terrain{
		logger: logger,
		props: Props{
			width:       common.WorldWidth,
			height:      common.WorldHeight,
			scale:       60,
			lacunarity:  2.0,
			persistence: 0.5,
			octaves:     5,
		},
		World: MapView{
			isMiniMap: false,
			logger:    logger,
			grid:      worldGrid,
		},
		MiniMap: MapView{
			isMiniMap: true,
			logger:    logger,
			grid:      miniMap,
		},
	}
}

func (tt TerrainType) Render() string {
	return fmt.Sprintf(TypeLookup[tt].style.PaddingLeft(1).PaddingRight(1).Render("%c"), TypeLookup[tt].symbol)
}

func (t Terrain) genTownCoords() common.Coordinates {
	coords := common.Coordinates{X: rand.Intn(common.WorldWidth), Y: rand.Intn(common.WorldHeight)}
	t.logger.Debug(fmt.Sprintf("Generating random town coordinates: %v", coords))
	return coords
}

func (t *Terrain) generateTowns(fn func() common.Coordinates) {
	t.logger.Info("Initializing %v towns", common.TotalTowns)
	for i := 0; i <= common.TotalTowns; i++ {
		for {
			coords := fn()
			if coords.X > 1 && coords.Y > 1 &&
				coords.X < common.WorldWidth-1 && coords.Y < common.WorldHeight &&
				t.World.grid[coords.X][coords.Y] == common.TypeBeach {

				if t.World.isAdjacentToWater(coords) {
					town.Create(coords, '⩎')
					t.World.grid[coords.X][coords.Y] = common.TypeTown
					// grow towns
					for _, a := range t.World.GetAdjacentCoords(coords) {
						if (t.World.grid[a.X][a.Y] == common.TypeLowland || t.World.grid[a.X][a.Y] == common.TypeBeach) && t.World.isAdjacentToWater(a) {
							t.World.grid[a.X][a.Y] = common.TypeTown
						}
					}
					break
				}
			}
		}
	}
}

func (t *Terrain) GenerateTowns() {
	t.generateTowns(t.genTownCoords)
}

func (t *Terrain) GenerateWorld() {
	t.logger.Info("Initializing world")
	noise := opensimplex.New(rand.Int63())

	for x := 0; x < t.props.width; x++ {
		for y := 0; y < t.props.height; y++ {
			//sample x and y and apply scale
			xFloat := float64(x) / t.props.scale
			yFloat := float64(y) / t.props.scale

			//init values for octave calculation
			frequency := 1.0
			amplitude := 1.0
			normalizeOctaves := 0.0
			total := 0.0

			//octave calculation
			for i := 0; i < t.props.octaves; i++ {
				total += noise.Eval2(xFloat*frequency, yFloat*frequency) * amplitude
				normalizeOctaves += amplitude
				amplitude *= t.props.persistence
				frequency *= t.props.lacunarity
			}

			//normalize to -1 to 1, and then from 0 to 1 (this is for the ability to use grayscale, if using colors could keep from -1 to 1)
			var s = (total/normalizeOctaves + 1) / 2
			if s > 0.60 {
				t.World.grid[x][y] = common.TypeDeepWater
			} else if s > 0.46 {
				t.World.grid[x][y] = common.TypeOpenWater
			} else if s > 0.42 {
				t.World.grid[x][y] = common.TypeShallowWater
			} else if s > 0.40 {
				t.World.grid[x][y] = common.TypeBeach
			} else if s > 0.31 {
				t.World.grid[x][y] = common.TypeLowland
			} else if s > 0.26 {
				t.World.grid[x][y] = common.TypeHighland
			} else if s > 0.21 {
				t.World.grid[x][y] = common.TypeRock
			} else {
				t.World.grid[x][y] = common.TypePeak
			}
		}
	}
	t.GenerateMiniMap()
}

func GetType(i int) (TerrainType, error) {
	for k := range TypeLookup {
		if int(k) == i {
			return k, nil
		}
	}
	return 0, errors.New("invalid type")
}

func (t *Terrain) GenerateMiniMap() {
	// Down-sample
	for i, row := range t.World.grid {
		for j, val := range row {
			// Calculate corresponding index in new slice
			newI := i / common.MiniMapFactor
			newJ := j / common.MiniMapFactor

			// Assign original value
			t.MiniMap.grid[newI][newJ] = val
		}
	}
}

func (t *Terrain) RandomPositionDeepWater() common.Coordinates {
	for {
		coords := common.Coordinates{X: rand.Intn(common.WorldWidth), Y: rand.Intn(common.WorldHeight)}
		if t.World.grid[coords.X][coords.Y] == common.TypeDeepWater {
			return coords
		}
	}
}

type AvatarReadOnly interface {
	GetX() int
	GetY() int
	GetMiniMapX() int
	GetMiniMapY() int
	Render() string
}

func (world MapView) GetTerrainType(x int, y int) TerrainType {
	return world.grid[x][y]
}

func (world MapView) isAdjacentToWater(coords common.Coordinates) bool {
	adjacentCoords := world.GetAdjacentCoords(coords)
	isAdjacentWater := false
	for _, a := range adjacentCoords {
		if world.grid[a.X][a.Y] == common.TypeShallowWater {
			isAdjacentWater = true
			break
		}
	}
	return isAdjacentWater
}

func (world MapView) GetAdjacentCoords(coords common.Coordinates) []common.Coordinates {
	var adjacentCoords []common.Coordinates
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			adjX := coords.X + i
			adjY := coords.Y + j
			if adjX < 0 || adjX >= world.GetWidth() || adjY < 0 || adjY >= world.GetHeight() {
				continue
			}
			adjacentCoords = append(adjacentCoords, common.Coordinates{X: adjX, Y: adjY})
		}
	}
	return adjacentCoords
}

func (world MapView) GetWidth() int {
	return len(world.grid[0])
}

func (world MapView) GetHeight() int {
	return len(world.grid)
}

func (world MapView) Paint(avatar AvatarReadOnly) string {
	left := 0
	top := 0
	worldHeight := len(world.grid)
	worldWidth := len(world.grid[0])
	viewHeight := worldHeight
	viewWidth := worldWidth
	rowWidth := worldWidth
	avatarX := avatar.GetX()
	avatarY := avatar.GetY()

	if world.isMiniMap {
		avatarX = avatar.GetMiniMapX()
		avatarY = avatar.GetMiniMapY()
		for _, o := range town.List {
			world.grid[o.GetMiniMapX()][o.GetMiniMapY()] = common.TypeTown
		}
	} else {
		// center viewport on avatar
		left = avatarX - (common.ViewWidth / 2)
		top = avatarY - (common.ViewHeight / 2)
		if left < 0 {
			left = 0
		}
		if top < 0 {
			top = 0
		}
		viewHeight = common.ViewHeight + top
		viewWidth = common.ViewWidth + left
		rowWidth = common.ViewWidth
	}

	viewport := table.New().BorderBottom(false).BorderTop(false).BorderLeft(false).BorderRight(false)

	world.logger.Debug(fmt.Sprintf("avatar position:  X:%v Y:%v", avatarX, avatarY))
	world.logger.Debug(fmt.Sprintf("viewport:  top:%v left:%v", top, left))
	world.logger.Debug(fmt.Sprintf("world:  height:%v width:%v", worldHeight, worldWidth))
	for y := top; y < worldHeight && y < viewHeight; y++ {
		var row = make([]string, rowWidth)
		for x := left; x < worldWidth && x < viewWidth; x++ {
			if x == avatarX && y == avatarY {
				row[x-left] = avatar.Render()
			} else {
				row[x-left] = world.grid[x][y].Render()
			}
		}
		viewport.Row(row...).BorderColumn(false)
	}

	return fmt.Sprintln(viewport)
}

func (world MapView) IsPassableByBoat(coordinates common.Coordinates) bool {
	tt := world.grid[coordinates.X][coordinates.Y]
	if TypeLookup[tt].RequiresBoat {
		return true
	}
	return false
}
