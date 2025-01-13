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
		t.Logger.Debug(fmt.Sprintf("[%v] Town at %v heatmap aborted with %v iterations", town.GetID(), town.GetPos(), count))
		return false
	} else {
		t.Logger.Debug(fmt.Sprintf("[%v] Town at %v heatmap completed with %v iterations", town.GetID(), town.GetPos(), count))
		return true
	}
}

func (h *HeatMap) Paint(avatar common.AvatarReadOnly, npcs []common.AvatarReadOnly, highlight common.ViewableEntity) string {
	// center viewport on avatar
	v := common.GetViewableArea(avatar.GetPos())
	rowWidth := common.ViewWidth

	viewport := table.New().BorderBottom(false).BorderTop(false).BorderLeft(false).BorderRight(false)

	// overlay map of all avatars
	overlay := make(map[string]common.AvatarReadOnly)
	c := avatar.GetPos()
	overlay[fmt.Sprintf("%03d%03d", c.X, c.Y)] = avatar

	// on the world map we draw the NPCs
	for _, n := range npcs {
		p := n.GetPos()
		overlay[fmt.Sprintf("%03d%03d", p.X, p.Y)] = n
	}

	for y := v.Top; y < v.Bottom; y++ {
		var row = make([]string, rowWidth)
		for x := v.Left; x < v.Right; x++ {
			item, ok := overlay[fmt.Sprintf("%03d%03d", x, y)]
			if ok {
				row[x-v.Left] = item.Render()
			} else {
				row[x-v.Left] = h.grid[x][y].Render()
			}
		}
		viewport.Row(row...).BorderColumn(false)
	}

	return fmt.Sprintln(viewport)
}

func (hc *HeatMapCost) Render() string {
	return fmt.Sprintf(createTerrainItem("0").PaddingLeft(1).PaddingRight(1).Render("%v"), *hc)
}
