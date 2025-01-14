package terrain

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/table"
	"go.uber.org/zap"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/screen"
)

type MapView struct {
	grid      [][]TerrainType
	logger    *zap.SugaredLogger
	isMiniMap bool
}

func (world MapView) isAdjacentToWater(c common.Coordinates) bool {
	adjacentCoords := world.GetAdjacentCoords(c)
	isAdjacentWater := false
	for _, a := range adjacentCoords {
		if world.GetPositionType(a) == TypeShallowWater {
			isAdjacentWater = true
			break
		}
	}
	return isAdjacentWater
}

func (world MapView) GetAdjacentCoords(c common.Coordinates) []common.Coordinates {
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

func (world MapView) GetWidth() int {
	return len(world.grid[0])
}

func (world MapView) GetHeight() int {
	return len(world.grid)
}

func (world MapView) Paint(avatar common.AvatarReadOnly, npcs []common.AvatarReadOnly, entity common.ViewableEntity) string {
	v := common.ViewableArea{}
	rowWidth := screen.ViewWidth

	viewport := table.New().BorderBottom(false).BorderTop(false).BorderLeft(false).BorderRight(false)

	// overlay map of all avatars
	overlay := make(map[string]common.AvatarReadOnly)

	world.logger.Info(fmt.Sprintf("ViewPort set to %v, %v", screen.ViewWidth, screen.ViewHeight))

	if world.isMiniMap {
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

func (world MapView) IsPassableByBoat(c common.Coordinates) bool {
	tt := world.GetPositionType(c)
	return TypeLookup[tt].RequiresBoat
}

func (world MapView) IsPassable(c common.Coordinates) bool {
	tt := world.GetPositionType(c)
	return TypeLookup[tt].Passable
}

func (world MapView) GetPositionType(c common.Coordinates) TerrainType {
	return world.grid[c.X][c.Y]
}

func (world MapView) SetPositionType(c common.Coordinates, t TerrainType) {
	world.grid[c.X][c.Y] = t
}

func (world MapView) IsLand(c common.Coordinates) bool {
	tt := world.grid[c.X][c.Y]
	if tt == TypeBeach || tt == TypeLowland || tt == TypeHighland || tt == TypePeak || tt == TypeRock {
		return true
	}
	return false
}
