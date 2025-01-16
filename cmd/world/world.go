package world

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/ojrac/opensimplex-go"
	"go.uber.org/zap"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/screen"
	"pirate-wars/cmd/terrain"
)

type ViewType int

const ViewTypeMainMap = 0
const ViewTypeHeatMap = 1
const ViewTypeMiniMap = 2

type Props struct {
	width       int
	height      int
	scale       float64
	lacunarity  float64
	persistence float64
	octaves     int
}

var WorldProps = Props{
	width:       common.WorldWidth,
	height:      common.WorldHeight,
	scale:       60,
	lacunarity:  2.0,
	persistence: 0.5,
	octaves:     5,
}

type MiniMapView struct {
	grid [][]terrain.TerrainType
}
type MapView struct {
	grid     [][]terrain.TerrainType
	logger   *zap.SugaredLogger
	miniMap  MiniMapView
	mapItems []MapItem
}

type MapItem interface {
	GetPos() common.Coordinates
	GetTerrainType() terrain.TerrainType
}

func (world *MapView) SetMapItem(m MapItem) {
	world.mapItems = append(world.mapItems, m)
}

func (world *MapView) IsAdjacentToWater(c common.Coordinates) bool {
	adjacentCoords := world.GetAdjacentCoords(c)
	isAdjacentWater := false
	for _, a := range adjacentCoords {
		if world.GetPositionType(a) == terrain.TypeShallowWater {
			isAdjacentWater = true
			break
		}
	}
	return isAdjacentWater
}

func (world *MapView) GetAdjacentCoords(c common.Coordinates) []common.Coordinates {
	var adjacentCoords []common.Coordinates
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			adjX := c.X + i
			adjY := c.Y + j
			if adjX < 0 || adjX >= world.GetWidth() || adjY < 0 || adjY >= world.GetHeight() {
				continue
			}
			adjacentCoords = append(adjacentCoords, common.Coordinates{X: adjX, Y: adjY})
		}
	}
	return adjacentCoords
}

func (world *MapView) GetWidth() int {
	return len(world.grid[0])
}

func (world *MapView) GetHeight() int {
	return len(world.grid)
}

func (world *MapView) IsPassableByBoat(c common.Coordinates) bool {
	tt := world.GetPositionType(c)
	return terrain.TypeLookup[tt].RequiresBoat
}

func (world *MapView) IsPassable(c common.Coordinates) bool {
	tt := world.GetPositionType(c)
	return terrain.TypeLookup[tt].Passable
}

func (world *MapView) GetPositionType(c common.Coordinates) terrain.TerrainType {
	return world.grid[c.X][c.Y]
}

func (world *MapView) SetPositionType(c common.Coordinates, tt terrain.TerrainType) {
	world.grid[c.X][c.Y] = tt
}

func (world *MapView) IsLand(c common.Coordinates) bool {
	tt := world.grid[c.X][c.Y]
	if tt == terrain.TypeBeach || tt == terrain.TypeLowland || tt == terrain.TypeHighland || tt == terrain.TypePeak || tt == terrain.TypeRock {
		return true
	}
	return false
}

func (world *MapView) RandomPositionDeepWater() common.Coordinates {
	for {
		c := common.Coordinates{X: rand.Intn(common.WorldWidth-2) + 1, Y: rand.Intn(common.WorldHeight-2) + 1}
		//t.Logger.Info(fmt.Sprintf("Random position deep water at: %v, %v", c, t.World.GetPositionType(c)))
		if world.GetPositionType(c) == terrain.TypeDeepWater {
			return c
		}
	}
}

func (world *MapView) Paint(avatar common.AvatarReadOnly, npcs []common.AvatarReadOnly, entity common.ViewableEntity, viewType ViewType) string {
	v := common.ViewableArea{}
	rowWidth := screen.ViewWidth

	viewport := table.New().BorderBottom(false).BorderTop(false).BorderLeft(false).BorderRight(false)

	// overlay map of all avatars
	overlay := make(map[string]common.AvatarReadOnly)

	world.logger.Info(fmt.Sprintf("ViewPort set to %v, %v", screen.ViewWidth, screen.ViewHeight))

	if viewType == ViewTypeMiniMap {
		v = common.ViewableArea{0, 0, len(world.grid[0]), len(world.grid)}
		// mini map views the whole map
		rowWidth = common.WorldWidth
		// always display main character avatar on the minimap
		mm := common.GetMiniMapScale(avatar.GetPos())
		overlay[fmt.Sprintf("%03d%03d", mm.X, mm.Y)] = avatar
	} else {
		v = common.GetViewableArea(avatar.GetPos())
		p := avatar.GetPos()
		overlay[fmt.Sprintf("%03d%03d", p.X, p.Y)] = avatar
		// on the world map we draw the NPCs
		for _, n := range npcs {
			c := n.GetPos()
			overlay[fmt.Sprintf("%03d%03d", c.X, c.Y)] = n
		}
	}

	h := entity.GetPos()
	if h.X >= 0 {
		world.logger.Debug(fmt.Sprintf("[%v] highlighting", entity.GetID()))
		// actual entity to examine, we should highlight it
		entity.Highlight()
		overlay[fmt.Sprintf("%03d%03d", h.X, h.Y)] = entity
	}

	world.logger.Info(fmt.Sprintf("Viewable Area %v", v))
	world.logger.Info(fmt.Sprintf("Player position %v", avatar.GetPos()))
	world.logger.Info(fmt.Sprintf("Painting world with %v viewable NPCs", len(npcs)))

	for y := v.Top; y < v.Bottom; y++ {
		var row = make([]string, rowWidth)
		for x := v.Left; x < v.Right; x++ {

			item, ok := overlay[fmt.Sprintf("%03d%03d", x, y)]
			if ok {
				row[x-v.Left] = item.Render()
			} else {
				//world.logger.Debug(
				//	fmt.Sprintf("row[%v] = world.grid[%v][%v] [row len(%v), gridX len(%v), gridY len(%v)]",
				//		x-v.Left, x, y, len(row), len(world.grid), len(world.grid[0])))
				row[x-v.Left] = world.grid[x][y].Render()
			}
		}
		viewport.Row(row...).BorderColumn(false)
	}

	return fmt.Sprintln(viewport)
}

func Init(logger *zap.SugaredLogger) *MapView {
	worldGrid := make([][]terrain.TerrainType, common.WorldHeight)
	for i := range worldGrid {
		worldGrid[i] = make([]terrain.TerrainType, common.WorldHeight)
	}
	world := MapView{
		logger:   logger,
		grid:     worldGrid,
		mapItems: []MapItem{},
	}
	world.logger.Info("Initializing world...")
	noise := opensimplex.New(rand.Int63())

	for x := 0; x < WorldProps.width; x++ {
		for y := 0; y < WorldProps.height; y++ {
			//sample x and y and apply scale
			xFloat := float64(x) / WorldProps.scale
			yFloat := float64(y) / WorldProps.scale

			//init values for octave calculation
			frequency := 1.0
			amplitude := 1.0
			normalizeOctaves := 0.0
			total := 0.0

			//octave calculation
			for i := 0; i < WorldProps.octaves; i++ {
				total += noise.Eval2(xFloat*frequency, yFloat*frequency) * amplitude
				normalizeOctaves += amplitude
				amplitude *= WorldProps.persistence
				frequency *= WorldProps.lacunarity
			}

			c := common.Coordinates{
				X: x,
				Y: y,
			}
			//normalize to -1 to 1, and then from 0 to 1 (this is for the ability to use grayscale, if using colors could keep from -1 to 1)
			var s = (total/normalizeOctaves + 1) / 2
			if s > 0.59 {
				world.SetPositionType(c, terrain.TypeDeepWater)
			} else if s > 0.44 {
				world.SetPositionType(c, terrain.TypeOpenWater)
			} else if s > 0.42 {
				world.SetPositionType(c, terrain.TypeShallowWater)
			} else if s > 0.40 {
				world.SetPositionType(c, terrain.TypeBeach)
			} else if s > 0.31 {
				world.SetPositionType(c, terrain.TypeLowland)
			} else if s > 0.26 {
				world.SetPositionType(c, terrain.TypeHighland)
			} else if s > 0.21 {
				world.SetPositionType(c, terrain.TypeRock)
			} else {
				world.SetPositionType(c, terrain.TypePeak)
			}
		}
	}
	return &world
}
