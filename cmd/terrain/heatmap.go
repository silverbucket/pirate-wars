package terrain

import (
	"fmt"
	"pirate-wars/cmd/common"
)

func (t *Terrain) GenerateTownHeatMap(town *Town) bool {
	// Define the starting point
	startX, startY := town.GetX(), town.GetY()

	// Queue contains points to visit, cost
	queue := [][]int{{startX, startY, 0}}
	count := 0
	// Perform Breadth-First Search
	for len(queue) > 0 {
		count++
		//t.Logger.Debugf("Queue length: %d", len(queue))
		j := queue[0]
		queue = queue[1:]

		c := common.Coordinates{j[0], j[1]}
		cost := j[2]

		//t.Logger.Infof("[towm %v] Processing %v, %v", t`own, x, y)
		if t.World.IsPassableByBoat(c) {
			//t.Logger.Debug(fmt.Sprintf("[town %v] Assigning cost %v, %v = %v [%v]", town, x, y, cost, t.Towns[town].heatMap[x][y]))
			if t.World.GetPositionType(c) == TypeShallowWater {
				// shallow water costs more (dangerous)
				town.SetHeatmapCost(c, cost+3)
				cost = cost + 4
			} else if t.World.GetPositionType(c) == TypeOpenWater {
				// open water faster than shallow, but not as fast as deep
				town.SetHeatmapCost(c, cost+1)
				cost = cost + 2
			} else {
				town.SetHeatmapCost(c, cost)
				cost = cost + 1
			}
		} else {
			if cost == 0 && t.World.GetPositionType(c) == TypeTown {
				// starting town is the cheapest
				town.SetHeatmapCost(c, cost)
			} else {
				// land is impassible
				town.SetHeatmapCost(c, common.MaxMovementCost)
			}
		}

		// Explore neighbors
		for _, dir := range common.Directions {
			n := common.Coordinates{c.X + dir.X, c.Y + dir.Y}

			//t.Logger.Debug(fmt.Sprintf("[town %v] Prepping direction %v, %v cost:%v", town, newX, newY, cost))

			// Check if the new point is within bounds and not visited
			if common.Inbounds(n) && town.GetHeatmapCost(n) == -1 {
				if !t.World.IsPassableByBoat(n) && t.World.GetPositionType(n) != TypeTown {
					town.SetHeatmapCost(n, common.MaxMovementCost)
				} else {
					//t.Logger.Debug(fmt.Sprintf("[town %v] (%v, %v) Adding direction %v, %v -- heatmap:%v", town, x, y, newX, newY, t.Towns[town].heatMap[newX][newY]))
					town.SetHeatmapCost(n, -2)
					queue = append(queue, []int{n.X, n.Y, cost})
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
