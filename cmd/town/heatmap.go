package town

import (
	"fmt"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/world"
)

const HeatmapUnprocessed = -1
const HeatmapQueued = -2
const MaxMovementCost = HeatMapCost(999999)
const LandMovementBase = HeatMapCost(5000)

type DirectionCost struct {
	Pos  common.Coordinates
	Cost HeatMapCost
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

func (town *Town) generateHeatMap(world *world.MapView) bool {
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

		c := j.Pos
		cost := j.Cost

		//t.Logger.Infof("[towm %v] Processing %v, %v", t`own, x, y)
		if world.IsPassableByBoat(c) {
			//t.Logger.Debug(fmt.Sprintf("[town %v] Assigning cost %v, %v = %v [%v]", town, x, y, cost, t.Towns[town].HeatMap[x][y]))
			if world.GetPositionType(c) == common.TerrainTypeShallowWater {
				// shallow water costs more (dangerous)
				cost = cost + 10
				town.HeatMap.SetCost(c, cost)
			} else if world.GetPositionType(c) == common.TerrainTypeOpenWater {
				// open water faster than shallow, but not as fast as deep
				cost = cost + 5
				town.HeatMap.SetCost(c, cost)
			} else {
				town.HeatMap.SetCost(c, cost)
			}
			cost = cost + 1
		} else {
			if cost == 0 && world.GetPositionType(c) == common.TerrainTypeTown {
				// starting town is the cheapest
				town.HeatMap.SetCost(c, cost)
			} else {
				// land currently impassible
				town.HeatMap.SetCost(c, MaxMovementCost)
			}
		}

		// Explore neighbors
		for _, dir := range common.Directions {
			n := common.Coordinates{X: c.X + dir.X, Y: c.Y + dir.Y}

			// Check if the new point is within bounds of the map and not visited
			if common.Inbounds(n) && town.HeatMap.GetCost(n) == HeatmapUnprocessed {
				if world.IsLand(n) {
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
		town.logger.Debug(fmt.Sprintf("[%v] Town at %v heatmap aborted with %v iterations", town.GetID(), town.GetPos(), count))
		return false
	} else {
		town.logger.Debug(fmt.Sprintf("[%v] Town at %v heatmap completed with %v iterations", town.GetID(), town.GetPos(), count))
		return true
	}
}

//func (h *HeatMap) Paint(avatar npc.AvatarReadOnly, npcs []npc.AvatarReadOnly, highlight common.ViewableEntity) *fyne.Container {
//	// center viewport on avatar
//	v := window.GetViewport(avatar.GetPos(), window.ViewableArea)
//
//	viewport := container.NewGridWithColumns(64)
//
//	// overlay map of all avatars
//	overlay := make(map[string]npc.AvatarReadOnly)
//	c := avatar.GetPos()
//	overlay[fmt.Sprintf("%03d%03d", c.X, c.Y)] = avatar
//
//	// on the world map we draw the NPCs
//	for _, n := range npcs {
//		p := n.GetPos()
//		overlay[fmt.Sprintf("%03d%03d", p.X, p.Y)] = n
//	}
//
//	for y := v.Top; y < v.Bottom; y++ {
//		for x := v.Left; x < v.Right; x++ {
//			item, ok := overlay[fmt.Sprintf("%03d%03d", x, y)]
//			if ok {
//				viewport.Add(item.Render())
//			} else {
//				viewport.Add(h.grid[x][y].Render())
//			}
//		}
//	}
//
//	return viewport
//}

func DecideDirection(o []DirectionCost, dest common.Coordinates) DirectionCost {
	lowestCost := MaxMovementCost
	choice := DirectionCost{}
	for _, e := range o {
		if e.Cost <= lowestCost && e.Cost >= 0 {
			lowestCost = e.Cost
			//possibilities = append(possibilities, e.pos)
			choice = e
		}
	}
	//return common.ClosestTo(dest, possibilities)
	return choice
}
