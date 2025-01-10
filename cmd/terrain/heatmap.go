package terrain

import (
	"fmt"
	"pirate-wars/cmd/common"
)

const HeatmapUnprocessed = -1
const HeatmapQueued = -2

type DirectionCost struct {
	pos  common.Coordinates
	cost int
}

type HeatMap struct {
	grid [][]int
}

func (h *HeatMap) SetCost(c common.Coordinates, v int) {
	h.grid[c.X][c.Y] = v
}

func (h *HeatMap) GetCost(c common.Coordinates) int {
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
			//t.Logger.Debug(fmt.Sprintf("[town %v] Assigning cost %v, %v = %v [%v]", town, x, y, cost, t.Towns[town].heatMap[x][y]))
			if t.World.GetPositionType(c) == TypeShallowWater {
				// shallow water costs more (dangerous)
				cost = cost + 10
				town.heatMap.SetCost(c, cost)
			} else if t.World.GetPositionType(c) == TypeOpenWater {
				// open water faster than shallow, but not as fast as deep
				cost = cost + 5
				town.heatMap.SetCost(c, cost)
			} else {
				town.heatMap.SetCost(c, cost)
			}
			cost = cost + 1
		} else {
			if cost == 0 && t.World.GetPositionType(c) == TypeTown {
				// starting town is the cheapest
				town.heatMap.SetCost(c, cost)
			} else {
				// land currently impassible
				town.heatMap.SetCost(c, common.MaxMovementCost)
			}
		}

		// Explore neighbors
		for _, dir := range common.Directions {
			n := common.Coordinates{c.X + dir.X, c.Y + dir.Y}

			// Check if the new point is within bounds of the map and not visited
			if common.Inbounds(n) && town.heatMap.GetCost(n) == HeatmapUnprocessed {
				if t.World.IsLand(n) {
					town.heatMap.SetCost(n, common.MaxMovementCost)
				} else {
					//t.Logger.Debug(fmt.Sprintf("[town %v] (%v, %v) Adding direction %v, %v -- heatmap:%v", town, x, y, newX, newY, t.Towns[town].heatMap[newX][newY]))
					town.heatMap.SetCost(n, HeatmapQueued)
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
