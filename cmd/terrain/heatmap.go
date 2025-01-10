package terrain

import (
	"fmt"
	"pirate-wars/cmd/common"
)

func (t *Terrain) GenerateTownHeatMap(town *Town) bool {
	// Define the starting point
	startX, startY := town.pos.X, town.pos.Y

	// Queue contains points to visit, cost
	queue := [][]int{{startX, startY, 0}}
	count := 0
	// Perform Breadth-First Search
	for len(queue) > 0 {
		count++
		//t.Logger.Debugf("Queue length: %d", len(queue))
		j := queue[0]
		queue = queue[1:]

		x, y, cost := j[0], j[1], j[2]

		//t.Logger.Infof("[towm %v] Processing %v, %v", t`own, x, y)
		if TypeLookup[t.World.grid[x][y]].RequiresBoat {
			//t.Logger.Debug(fmt.Sprintf("[town %v] Assigning cost %v, %v = %v [%v]", town, x, y, cost, t.Towns[town].heatMap[x][y]))
			if t.World.grid[x][y] == TypeShallowWater {
				// shallow water costs more (dangerous)
				town.heatMap[x][y] = cost + 1
			} else {
				town.heatMap[x][y] = cost
			}
			cost++
		} else {
			// land or town
			if cost == 0 && t.World.grid[x][y] == TypeTown {
				town.heatMap[x][y] = cost
			} else {
				town.heatMap[x][y] = common.MaxMovementCost
			}
		}

		// Explore neighbors
		for _, dir := range common.Directions {
			newX, newY := x+dir.X, y+dir.Y
			//t.Logger.Debug(fmt.Sprintf("[town %v] Prepping direction %v, %v cost:%v", town, newX, newY, cost))

			// Check if the new point is within bounds and not visited
			if newX >= 0 && newX < common.WorldHeight && newY >= 0 && newY < common.WorldWidth && town.heatMap[newX][newY] == -1 {
				if !TypeLookup[t.World.grid[newX][newY]].RequiresBoat && t.World.grid[newX][newY] != TypeTown {
					town.heatMap[newX][newY] = common.MaxMovementCost
				} else {
					//t.Logger.Debug(fmt.Sprintf("[town %v] (%v, %v) Adding direction %v, %v -- heatmap:%v", town, x, y, newX, newY, t.Towns[town].heatMap[newX][newY]))
					town.heatMap[newX][newY] = -2
					queue = append(queue, []int{newX, newY, cost})
				}
			}
		}
	}
	if count < 200 {
		t.Logger.Debug(fmt.Sprintf("Town at %v, %v heatmap aborted at just %v iterations", town.pos.X, town.pos.Y, count))
		return false
	} else {
		return true
	}
}

func (t *Terrain) GenerateHeatMaps() {
	for _, town := range t.Towns {
		t.Logger.Info(fmt.Sprintf("Generating heatmap for town at %v, %v", town.pos.X, town.pos.Y))
		t.GenerateTownHeatMap(&town)
	}
}
