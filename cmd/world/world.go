package world

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
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

var minimapPopup *widget.PopUp

const TileSizePx = 512
const TileSizeCells = TileSizePx / 12

type Props struct {
	scale       float64
	lacunarity  float64
	persistence float64
	octaves     int
}

type Tile struct {
	image *canvas.Image
	X, Y  int
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

type OverlayItems interface {
	GetPos() common.Coordinates
	GetTerrainType() terrain.Type
}

func (world *MapView) SetMapItem(m OverlayItems) {
	world.overlayItems = append(world.overlayItems, m)
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

func (world *MapView) GetPositionType(c common.Coordinates) terrain.Type {
	return world.terrain.Cells[c.X][c.Y]
}

func (world *MapView) SetPositionType(c common.Coordinates, tt terrain.Type) {
	world.terrain.Cells[c.X][c.Y] = tt
}

func (world *MapView) IsLand(c common.Coordinates) bool {
	tt := world.terrain.Cells[c.X][c.Y]
	if tt == terrain.TypeBeach || tt == terrain.TypeLowland || tt == terrain.TypeHighland || tt == terrain.TypePeak || tt == terrain.TypeRock {
		return true
	}
	return false
}

func (world *MapView) RandomPositionDeepWater() common.Coordinates {
	for {
		c := common.Coordinates{X: rand.Intn(common.WorldCols-2) + 1, Y: rand.Intn(common.WorldRows-2) + 1}
		//terrain.Logger.Info(fmt.Sprintf("Random position deep water at: %v, %v", c, terrain.World.GetPositionType(c)))
		if world.GetPositionType(c) == terrain.TypeDeepWater {
			return c
		}
	}
}

func (world *MapView) renderTile(tx, ty int) *Tile {
	img := image.NewRGBA(image.Rect(0, 0, TileSizePx, TileSizePx))
	for y := ty; y < ty+TileSizeCells && y < 800; y++ {
		for x := tx; x < tx+TileSizeCells && x < 800; x++ {
			for py := 0; py < 12; py++ {
				for px := 0; px < 12; px++ {
					img.Set((x-tx)*12+px, (y-ty)*12+py, world.terrain.Cells[y][x].GetColor())
				}
			}
		}
	}
	fyneImg := canvas.NewImageFromImage(img)
	fyneImg.FillMode = canvas.ImageFillStretch
	fyneImg.Resize(fyne.NewSize(TileSizePx, TileSizePx))
	return &Tile{image: fyneImg, X: tx, Y: ty}
}

func (world *MapView) generateBaseMinimapImage() {
	world.logger.Info(fmt.Sprintf("Generating minimap"))
	cols := common.WorldCols
	rows := common.WorldRows
	cellWidth := float32(window.MiniMapArea.Width) / float32(cols)
	cellHeight := float32(window.MiniMapArea.Height) / float32(rows)

	world.minimap = world.createRawMapImage(cellWidth, cellHeight, cols, rows, window.MiniMapArea.Width, window.MiniMapArea.Height)
}

func (world *MapView) createRawMapImage(cellWidth, cellHeight float32, cols, rows int, imageWidth, imageHeight int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	fmt.Print(".")
	count := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			for y := int(float32(r) * cellHeight); y < int(float32(r+1)*cellHeight); y++ {
				for x := int(float32(c) * cellWidth); x < int(float32(c+1)*cellWidth); x++ {
					img.Set(x, y, world.terrain.Cells[c][r].GetColor())
					count++
					if count%2000 == 0 {
						fmt.Print(".")
					}
				}
			}
		}
	}
	fmt.Println("done")
	return img
}

func (world *MapView) getMinimapWithDot(pos common.Coordinates) *image.RGBA {
	cols := common.WorldCols
	rows := common.WorldRows

	// Create a copy of the base image
	img := image.NewRGBA(world.minimap.Rect)
	draw.Draw(img, img.Bounds(), world.minimap, image.Point{}, draw.Src)

	// Calculate pixel position of player on minimap
	cellWidth := float32(window.MiniMapArea.Width) / float32(cols)
	cellHeight := float32(window.MiniMapArea.Height) / float32(rows)

	x := int(float32(pos.X) * cellWidth)
	y := int(float32(pos.Y) * cellHeight)

	// Draw a small white dot (e.g., 5x5 pixels)
	dotSize := 5
	for dy := -dotSize / 2; dy <= dotSize/2; dy++ {
		for dx := -dotSize / 2; dx <= dotSize/2; dx++ {
			px := x + dx
			py := y + dy
			if px >= 0 && px < window.MiniMapArea.Width && py >= 0 && py < window.MiniMapArea.Height {
				img.Set(px, py, color.White)
			}
		}
	}

	return img
}

func (world *MapView) ShowMinimapPopup(pos common.Coordinates, w fyne.Window) {
	overlay := canvas.NewRectangle(color.NRGBA{0, 0, 0, 128}) // 50% opacity
	overlay.Resize(fyne.NewSize(float32(window.Window.Width), float32(window.Window.Height)))

	minimapImage := canvas.NewImageFromImage(world.getMinimapWithDot(pos))
	minimapContainer := container.NewStack(overlay, minimapImage)
	minimapPopup = widget.NewPopUp(minimapContainer, w.Canvas())
	minimapPopup.Resize(fyne.NewSize(float32(window.MiniMapArea.Width), float32(window.MiniMapArea.Height)))
	minimapPopup.Move(fyne.NewPos(float32(window.Window.Width-window.MiniMapArea.Width)/2, float32(window.Window.Height-window.MiniMapArea.Height)/2))
	minimapPopup.Show()
}

func (world *MapView) HideMinimapPopup() {
	if minimapPopup != nil {
		minimapPopup.Hide()
	}
}

func (world *MapView) GetViewPort() *fyne.Container {
	return world.viewPort
}

func (world *MapView) Paint(avatar entities.AvatarReadOnly, npcs []entities.AvatarReadOnly, entity entities.ViewableEntity) {
	p := avatar.GetPos()
	h := entity.GetPos() // potential entity to highlight (selected)
	vpr := window.GetViewportRegion(p)

	// overlay map of all avatars, player and npcs
	// instead of terrain, in these overlay positions we generate the avatars
	overlay := make(map[string]entities.AvatarReadOnly)
	overlay[fmt.Sprintf("%03d%03d", p.X, p.Y)] = avatar

	for _, n := range npcs {
		c := n.GetPos()
		overlay[fmt.Sprintf("%03d%03d", c.X, c.Y)] = n
	}

	// if the entity to highlight has real coords, we add it to the overlay
	if h.X >= 0 {
		world.logger.Debug("[%v] highlighting", entity.GetID())
		// actual entity to examine, we should highlight it
		entity.Highlight()
		// Don't add entity to overlay as it doesn't implement AvatarReadOnly
		// overlay[fmt.Sprintf("%03d%03d", h.X, h.Y)] = entity
	}

	// world.logger.Info("--")
	// world.logger.Info("Player position %+v", p)
	// world.logger.Info("Painting world with %v viewable NPCs", len(npcs))
	// world.logger.Info("Viewport Region %+v", v)
	// world.logger.Info("Increment amount %+v", inc)
	// world.logger.Info("Grid length %v", len(g.Objects))

	vpIdx := 0
	for x := 0; x < vpr.Cols; x++ {
		for y := 0; y < vpr.Rows; y++ {
			// Calculate map coordinates
			mapX := vpr.X + x
			mapY := vpr.Y + y

			// Skip if outside map bounds
			if mapX < 0 || mapX >= common.WorldCols || mapY < 0 || mapY >= common.WorldRows {
				continue
			}

			var cell fyne.CanvasObject
			item, ok := overlay[fmt.Sprintf("%03d%03d", mapX, mapY)]
			if ok {
				cell = item.Render()
			} else {
				cell = world.terrain.Cells[mapX][mapY].Render()
			}

			cell.Resize(fyne.NewSize(float32(window.CellSize), float32(window.CellSize)))
			cell.Move(fyne.NewPos(float32(x*window.CellSize), float32(y*window.CellSize)))
			world.viewPort.Objects[vpIdx] = cell
			vpIdx++
		}
	}

	world.viewPort.Resize(fyne.NewSize(float32(window.ViewPort.Dimensions.Width), float32(window.ViewPort.Dimensions.Height)))
	fyne.Do(func() {
		world.viewPort.Refresh()
	})
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

	// initialize viewport cells
	for range window.ViewPort.Region.Cols {
		for range window.ViewPort.Region.Rows {
			world.viewPort.Add(canvas.NewRectangle(color.Black))
		}
	}

	world.generateBaseMinimapImage()
	return &world
}
