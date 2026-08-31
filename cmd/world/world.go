package world

import (
	"image"
	"image/color"

	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/terrain"

	"github.com/ojrac/opensimplex-go"
	"go.uber.org/zap"
)

type ViewType int

const ViewTypeMainMap = 0
const ViewTypeHeatMap = 1
const ViewTypeMiniMap = 2
const ViewTypeExamine = 3
const ViewTypeDock = 4
const ViewTypeHail = 5

// ViewTypeHelp is the controls and sailing-model screen behind "?".
const ViewTypeHelp = 6

// ViewTypeQuitConfirm guards Ctrl+Q, which otherwise ends the voyage instantly.
const ViewTypeQuitConfirm = 7

type Props struct {
	scale       float64
	lacunarity  float64
	persistence float64
	octaves     int
}

var WorldProps = Props{
	scale:       60,
	lacunarity:  2.0,
	persistence: 0.5,
	octaves:     5,
}

type MapView struct {
	logger       *zap.SugaredLogger
	terrain      *terrain.Terrain
	minimap      *image.RGBA
	overlayItems []OverlayItems
	// glides is renderer-only smoothing state: each ship's on-screen position
	// chasing its logical cell. See visualPos in render.go.
	glides map[string]*glideState
}

type MinimapOverlay struct {
	pos   common.Coordinates
	color color.Color
}

type OverlayItems interface {
	GetPos() common.Coordinates
	GetTerrainType() common.TerrainType
	GetTileImage() image.Image
}

// viewablePaintWrapper adapts ViewableEntity for paint overlays that need sailing fields.
type viewablePaintWrapper struct {
	entities.ViewableEntity
}

func (v viewablePaintWrapper) GetFacing() common.Facing { return common.FacingN }
func (v viewablePaintWrapper) MovedThisTick() bool      { return false }

func (world *MapView) SetMapItem(m OverlayItems) {
	world.overlayItems = append(world.overlayItems, m)
}

func (world *MapView) IsAdjacentToWater(c common.Coordinates) bool {
	adjacentCoords := world.GetAdjacentCoords(c)
	isAdjacentWater := false
	for _, a := range adjacentCoords {
		if world.GetPositionType(a) == common.TerrainTypeShallowWater {
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
	return len(world.terrain.Cells[0])
}

func (world *MapView) GetHeight() int {
	return len(world.terrain.Cells)
}

func (world *MapView) IsPassableByBoat(c common.Coordinates) bool {
	tt := world.GetPositionType(c)
	return terrain.TypeLookup[tt].RequiresBoat
}

func (world *MapView) IsPassable(c common.Coordinates) bool {
	tt := world.GetPositionType(c)
	return terrain.TypeLookup[tt].Passable
}

func (world *MapView) GetPositionType(c common.Coordinates) common.TerrainType {
	return world.terrain.Cells[c.X][c.Y]
}

func (world *MapView) SetPositionType(c common.Coordinates, tt common.TerrainType) {
	world.terrain.Cells[c.X][c.Y] = tt
}

func (world *MapView) IsLand(c common.Coordinates) bool {
	tt := world.terrain.Cells[c.X][c.Y]
	if tt == common.TerrainTypeBeach || tt == common.TerrainTypeLowland || tt == common.TerrainTypeHighland || tt == common.TerrainTypePeak || tt == common.TerrainTypeRock {
		return true
	}
	return false
}

func (world *MapView) RandomPositionDeepWater() common.Coordinates {
	for {
		c := common.Coordinates{X: rand.Intn(common.WorldCols-2) + 1, Y: rand.Intn(common.WorldRows-2) + 1}
		//terrain.Logger.Info(fmt.Sprintf("Random position deep water at: %v, %v", c, terrain.World.GetPositionType(c)))
		if world.GetPositionType(c) == common.TerrainTypeDeepWater {
			return c
		}
	}
}
func Init(logger *zap.SugaredLogger) *MapView {
	t := &terrain.Terrain{}
	world := MapView{
		logger:       logger,
		terrain:      t,
		overlayItems: []OverlayItems{},
	}

	world.logger.Info("Initializing world...")
	noise := opensimplex.New(rand.Int63())

	for x := 0; x < common.WorldCols; x++ {
		for y := 0; y < common.WorldRows; y++ {
			// sample x and y and apply scale
			xFloat := float64(x) / WorldProps.scale
			yFloat := float64(y) / WorldProps.scale

			// init values for octave calculation
			frequency := 1.0
			amplitude := 1.0
			normalizeOctaves := 0.0
			total := 0.0

			// octave calculation
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
			// normalize to -1 to 1, and then from 0 to 1 (this is for the ability to use grayscale, if using colors could keep from -1 to 1)
			var s = (total/normalizeOctaves + 1) / 2
			var terrain common.TerrainType
			if s > 0.59 {
				terrain = common.TerrainTypeDeepWater
			} else if s > 0.44 {
				terrain = common.TerrainTypeOpenWater
			} else if s > 0.42 {
				terrain = common.TerrainTypeShallowWater
			} else if s > 0.40 {
				terrain = common.TerrainTypeBeach
			} else if s > 0.31 {
				terrain = common.TerrainTypeLowland
			} else if s > 0.26 {
				terrain = common.TerrainTypeHighland
			} else if s > 0.21 {
				terrain = common.TerrainTypeRock
			} else {
				terrain = common.TerrainTypePeak
			}
			world.SetPositionType(c, terrain)
		}
	}

	world.generateMinimapImage()
	return &world
}
