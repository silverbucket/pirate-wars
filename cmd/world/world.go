package world

import (
	"image"
	"image/color"

	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/resources"
	"pirate-wars/cmd/terrain"
	"pirate-wars/cmd/window"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ojrac/opensimplex-go"
	"go.uber.org/zap"
)

type ViewType int

const ViewTypeMainMap = 0
const ViewTypeHeatMap = 1
const ViewTypeMiniMap = 2
const ViewTypeExamine = 3

var minimapPopup *widget.PopUp

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
	viewPort     *fyne.Container
	minimap      *image.RGBA
	overlayItems []OverlayItems
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
func (world *MapView) GetViewPort() *fyne.Container {
	return world.viewPort
}

func (world *MapView) generateViewPort() {
	// initialize viewport cells
	for x := 0; x < window.ViewPort.Region.Cols; x++ {
		for y := 0; y < window.ViewPort.Region.Rows; y++ {
			cell := container.NewStack(
				canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, window.CellSize, window.CellSize))),
				canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, window.CellSize, window.CellSize))),
			)
			cell.Resize(fyne.NewSize(float32(window.CellSize), float32(window.CellSize)))
			cell.Move(fyne.NewPos(float32(x*window.CellSize), float32(y*window.CellSize)))
			world.viewPort.Add(cell)
		}
	}
}

func (world *MapView) Paint(avatar entities.AvatarReadOnly, npcs []entities.AvatarReadOnly, highlight entities.ViewableEntity) {
	p := avatar.GetPos()
	h := highlight.GetPos()
	vpr := window.GetViewportRegion(p)

	// Create overlay map
	overlay := make(map[int]entities.AvatarReadOnly, len(npcs)+2)
	overlay[common.CoordToKey(p)] = avatar
	for _, n := range npcs {
		overlay[common.CoordToKey(n.GetPos())] = n
	}

	// if the entity to highlight has real coords, we add it to the overlay
	if h.X >= 0 {
		world.logger.Debug("[%v] highlighting", highlight.GetID())
		highlight.Highlight(true)
		overlay[common.CoordToKey(h)] = highlight
	}

	vpIdx := 0
	needsRefresh := false

	// Pre-calculate cell positions to avoid repeated calculations
	cellPositions := make([]common.Coordinates, vpr.Cols*vpr.Rows)
	for x := 0; x < vpr.Cols; x++ {
		for y := 0; y < vpr.Rows; y++ {
			mapX := vpr.X + x
			mapY := vpr.Y + y
			if mapX >= 0 && mapX < common.WorldCols && mapY >= 0 && mapY < common.WorldRows {
				cellPositions[vpIdx] = common.Coordinates{X: mapX, Y: mapY}
			}
			vpIdx++
		}
	}

	// Process all cells in the viewport
	vpIdx = 0
	for _, pos := range cellPositions {
		if pos.X < 0 || pos.X >= common.WorldCols || pos.Y < 0 || pos.Y >= common.WorldRows {
			vpIdx++
			continue
		}

		cell := world.viewPort.Objects[vpIdx].(*fyne.Container)
		terrainImg := cell.Objects[0].(*canvas.Image)
		entityImg := cell.Objects[1].(*canvas.Image)

		var newTerrainImage image.Image
		var newEntityImage image.Image

		if item, ok := overlay[common.CoordToKey(pos)]; ok {
			newEntityImage = item.GetTileImage()
		} else {
			newEntityImage = image.NewRGBA(image.Rect(0, 0, window.CellSize, window.CellSize))
		}

		newTerrainImage = resources.GetTerrainTile(world.terrain.Cells[pos.X][pos.Y])

		if terrainImg.Image != newTerrainImage {
			terrainImg.Image = newTerrainImage
			needsRefresh = true
		}

		if entityImg.Image != newEntityImage {
			entityImg.Image = newEntityImage
			needsRefresh = true
		}
		vpIdx++
	}

	if needsRefresh {
		world.viewPort.Refresh()
	}
}

func Init(logger *zap.SugaredLogger) *MapView {
	t := &terrain.Terrain{}
	world := MapView{
		logger:       logger,
		terrain:      t,
		viewPort:     container.NewWithoutLayout(),
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

	world.generateViewPort()
	world.generateMinimapImage()
	return &world
}
