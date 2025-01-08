package terrain

import (
	"pirate-wars/cmd/common"
	"sync"
)

func (t *Terrain) assignCost(town int, c common.Coordinates, v int) bool {
	canContinue := false
	if t.Towns[town].heatMap[c.X][c.Y] > 0 {
		// grid point has already been processed
		return false
	}
	if TypeLookup[t.World.grid[c.X][c.Y]].RequiresBoat {
		t.Towns[town].heatMap[c.X][c.Y] = v
		canContinue = true
	} else {
		// impassible
		t.Towns[town].heatMap[c.X][c.Y] = 9999
	}
	return canContinue
}

func (t *Terrain) generateHeatMapX(town int, startCoords common.Coordinates, cost int, g *sync.WaitGroup) {
	defer g.Done()
	upCost := cost
	for i := startCoords.X + 1; i < common.WorldHeight-i; i++ {
		canContinue := t.assignCost(town, common.Coordinates{X: i, Y: startCoords.Y}, upCost)
		upCost++
		g.Add(1)
		go t.generateHeatMapY(town, common.Coordinates{X: i, Y: startCoords.Y}, upCost, g)
		if !canContinue {
			break
		}
	}
	downCost := cost
	for i := startCoords.X - 1; i < 0; i-- {
		canContinue := t.assignCost(town, common.Coordinates{X: i, Y: startCoords.Y}, downCost)
		downCost++
		g.Add(1)
		t.generateHeatMapY(town, common.Coordinates{X: i, Y: startCoords.Y}, downCost, g)
		if !canContinue {
			break
		}
	}
}

func (t *Terrain) generateHeatMapY(town int, startCoords common.Coordinates, cost int, g *sync.WaitGroup) {
	defer g.Done()
	rightCost := cost
	for i := startCoords.Y + 1; i < common.WorldWidth-i; i++ {
		canContinue := t.assignCost(town, common.Coordinates{X: startCoords.X, Y: i}, rightCost)
		rightCost++
		g.Add(1)
		t.generateHeatMapX(town, common.Coordinates{X: startCoords.X, Y: i}, rightCost, g)
		if !canContinue {
			break
		}
	}
	leftCost := cost
	for i := startCoords.Y - 1; i > 0; i-- {
		canContinue := t.assignCost(town, common.Coordinates{X: startCoords.X, Y: i}, leftCost)
		leftCost++
		g.Add(1)
		t.generateHeatMapX(town, common.Coordinates{X: startCoords.X, Y: i}, leftCost, g)
		if !canContinue {
			return
		}
	}
}

func (t *Terrain) GenerateHeatMaps() {
	for i := range t.Towns {
		gx := &sync.WaitGroup{}
		gx.Add(1)
		t.generateHeatMapX(i, t.Towns[i].pos, 2, gx)
		gx.Wait()

		gy := &sync.WaitGroup{}
		gy.Add(1)
		t.generateHeatMapY(i, t.Towns[i].pos, 2, gy)
		gy.Wait()
	}
}
