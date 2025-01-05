package terrain

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/table"
	"go.uber.org/zap"
	"pirate-wars/cmd/common"
)

type MapView struct {
	grid      [][]TerrainType
	logger    *zap.SugaredLogger
	isMiniMap bool
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
		if world.grid[a.X][a.Y] == TypeShallowWater {
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

func (world MapView) Paint(a AvatarReadOnly) string {
	left := 0
	top := 0
	worldHeight := len(world.grid)
	worldWidth := len(world.grid[0])
	viewHeight := worldHeight
	viewWidth := worldWidth
	rowWidth := worldWidth
	avatarX := a.GetX()
	avatarY := a.GetY()

	if world.isMiniMap {
		avatarX = a.GetMiniMapX()
		avatarY = a.GetMiniMapY()
		for _, o := range Towns {
			world.grid[o.GetMiniMapX()][o.GetMiniMapY()] = TypeTown
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

	overlay := make(map[string]AvatarReadOnly)
	overlay[fmt.Sprintf("%v%v", avatarX, avatarY)] = a
	if !world.isMiniMap {
		// on the world map we draw the npcs
		for _, n := range NPCs {
			overlay[fmt.Sprintf("%v%v", n.GetX(), n.GetY())] = &n
		}
	}

	//world.logger.Debug(fmt.Sprintf("avatar position:  X:%v Y:%v", avatarX, avatarY))
	//world.logger.Debug(fmt.Sprintf("viewport:  top:%v left:%v", top, left))
	//world.logger.Debug(fmt.Sprintf("world:  height:%v width:%v", worldHeight, worldWidth))
	for y := top; y < worldHeight && y < viewHeight; y++ {
		var row = make([]string, rowWidth)
		for x := left; x < worldWidth && x < viewWidth; x++ {
			item, ok := overlay[fmt.Sprintf("%v%v", x, y)]
			if ok {
				row[x-left] = item.Render()
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
	return TypeLookup[tt].RequiresBoat
}

func (world MapView) IsPassable(coordinates common.Coordinates) bool {
	tt := world.grid[coordinates.X][coordinates.Y]
	return TypeLookup[tt].Passable
}
