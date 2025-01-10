package terrain

import (
	"fmt"
	"github.com/charmbracelet/lipgloss/table"
	"pirate-wars/cmd/common"
)

const HeatmapUnprocessed = -1
const HeatmapQueued = -2
const MaxMovementCost = HeatMapCost(999999)
const LandMovementBase = HeatMapCost(5000)

type DirectionCost struct {
	pos  common.Coordinates
	cost HeatMapCost
}

type HeatMapCost int

type HeatMap struct {
	grid [][]HeatMapCost
}

func (h *HeatMap) SetCost(c common.Coordinates, v HeatMapCost) {
	h.grid[c.X][c.Y] = v
}

func (h *HeatMap) GetCost(c common.Coordinates) HeatMapCost {
	return h.grid[c.X][c.Y]
}

func (t *Terrain) GenerateTownHeatMap(town *Town) bool {
	// Define the starting point
	// Queue contains points to visit, cost
	queue := []DirectionCost{{town.GetPos(), 0}}
	count := 0
	// Perform Breadth-First Search
	for len(queue) > 0 {
		count++
		//t.Logger.Debugf("Queue length: %d", len(queue))
		j := queue[0]
		queue = queue[1:]

		c := j.pos
		cost := j.cost

		//t.Logger.Infof("[towm %v] Processing %v, %v", t`own, x, y)
		if t.World.IsPassableByBoat(c) {
			//t.Logger.Debug(fmt.Sprintf("[town %v] Assigning cost %v, %v = %v [%v]", town, x, y, cost, t.Towns[town].HeatMap[x][y]))
			if t.World.GetPositionType(c) == TypeShallowWater {
				// shallow water costs more (dangerous)
				cost = cost + 10
				town.HeatMap.SetCost(c, cost)
			} else if t.World.GetPositionType(c) == TypeOpenWater {
				// open water faster than shallow, but not as fast as deep
				cost = cost + 5
				town.HeatMap.SetCost(c, cost)
			} else {
				town.HeatMap.SetCost(c, cost)
			}
			cost = cost + 1
		} else {
			if cost == 0 && t.World.GetPositionType(c) == TypeTown {
				// starting town is the cheapest
				town.HeatMap.SetCost(c, cost)
			} else {
				// land currently impassible
				town.HeatMap.SetCost(c, MaxMovementCost)
			}
		}

		// Explore neighbors
		for _, dir := range common.Directions {
			n := common.Coordinates{c.X + dir.X, c.Y + dir.Y}

			// Check if the new point is within bounds of the map and not visited
			if common.Inbounds(n) && town.HeatMap.GetCost(n) == HeatmapUnprocessed {
				if t.World.IsLand(n) {
					town.HeatMap.SetCost(n, MaxMovementCost)
				} else {
					//t.Logger.Debug(fmt.Sprintf("[town %v] (%v, %v) Adding direction %v, %v -- heatmap:%v", town, x, y, newX, newY, t.Towns[town].HeatMap[newX][newY]))
					town.HeatMap.SetCost(n, HeatmapQueued)
					queue = append(queue, DirectionCost{n, cost})
				}
			}
		}
	}
	if count < 200 {
		t.Logger.Debug(fmt.Sprintf("[%v] Town at %v heatmap aborted with %v iterations", town.GetId(), town.GetPos(), count))
		return false
	} else {
		t.Logger.Debug(fmt.Sprintf("[%v] Town at %v heatmap completed with %v iterations", town.GetId(), town.GetPos(), count))
		return true
	}
}

//func (t *Terrain) GenerateHeatMaps() {
//	for _, town := range t.Towns {
//		t.Logger.Info(fmt.Sprintf("Generating heatmap for town at %v, %v", town.pos.X, town.pos.Y))
//		t.GenerateTownHeatMap(&town)
//	}
//}

func (h *HeatMap) Paint(avatar AvatarReadOnly, npcs []AvatarReadOnly) string {
	left := 0
	top := 0
	worldHeight := len(h.grid)
	worldWidth := len(h.grid[0])
	viewHeight := worldHeight
	viewWidth := worldWidth
	rowWidth := worldWidth

	viewport := table.New().BorderBottom(false).BorderTop(false).BorderLeft(false).BorderRight(false)

	// overlay map of all avatars
	overlay := make(map[string]AvatarReadOnly)

	// center viewport on avatar
	left = avatar.GetX() - (common.ViewWidth / 3)
	top = avatar.GetY() - (common.ViewHeight / 3)
	if left < 0 {
		left = 0
	}
	if top < 0 {
		top = 0
	}
	viewHeight = common.ViewHeight + top
	viewWidth = common.ViewWidth + left
	rowWidth = common.ViewWidth

	overlay[fmt.Sprintf("%03d%03d", avatar.GetX(), avatar.GetY())] = avatar
	// on the world map we draw the NPCs
	for _, n := range npcs {
		overlay[fmt.Sprintf("%03d%03d", n.GetX(), n.GetY())] = n
	}

	//world.logger.Debug(fmt.Sprintf("avatar position:  X:%v Y:%v", avs[0].GetX, avs[0].GetY()))
	for y := top; y < worldHeight && y < viewHeight; y++ {
		var row = make([]string, rowWidth)
		for x := left; x < worldWidth && x < viewWidth; x++ {
			item, ok := overlay[fmt.Sprintf("%03d%03d", x, y)]
			if ok {
				row[x-left] = item.Render()
			} else {
				row[x-left] = h.grid[x][y].Render()
			}
		}
		viewport.Row(row...).BorderColumn(false)
	}

	return fmt.Sprintln(viewport)
}

func (hc *HeatMapCost) Render() string {
	return fmt.Sprintf(createTerrainItem("0").PaddingLeft(1).PaddingRight(1).Render("%v"), *hc)
}
