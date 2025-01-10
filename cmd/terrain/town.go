package terrain

import (
	"fmt"
	"math/rand"
	"pirate-wars/cmd/common"
)

type HeatMap = [][]int

type Town struct {
	id      string
	pos     []common.Coordinates
	heatMap HeatMap
}

func (t *Town) GetId() string {
	return t.id
}

func (t *Town) GetX() int {
	return t.pos[0].X
}
func (t *Town) GetMiniMapX() int {
	return t.pos[0].X / common.MiniMapFactor
}

func (t *Town) GetY() int {
	return t.pos[0].Y
}

func (t *Town) GetMiniMapY() int {
	return t.pos[0].Y / common.MiniMapFactor
}

func (t *Town) GetMiniMapPos() common.Coordinates {
	return common.Coordinates{
		X: t.pos[0].X / common.MiniMapFactor,
		Y: t.pos[0].Y / common.MiniMapFactor,
	}
}
func (t *Town) GetPos() common.Coordinates {
	return t.pos[0]
}

func (t *Town) SetHeatmapCost(c common.Coordinates, v int) {
	t.heatMap[c.X][c.Y] = v
}

func (t *Town) GetHeatmapCost(c common.Coordinates) int {
	return t.heatMap[c.X][c.Y]
}

func (t *Terrain) MakeGhostTown(town *Town) {
	t.Logger.Info(fmt.Sprintf("[%v] Town turns to ghost town at %v, %v", town.id, town.GetX(), town.GetY()))
	for _, c := range town.pos {
		t.World.SetPositionType(c, TypeGhostTown)
	}
}

func (t *Terrain) CreateTown(c common.Coordinates) Town {
	var heatMap = make([][]int, common.WorldHeight)

	for i := range heatMap {
		heatMap[i] = make([]int, common.WorldWidth)
		for j := range heatMap[i] {
			heatMap[i][j] = -1
		}
	}

	town := Town{
		id:      common.GenID(c),
		pos:     []common.Coordinates{c},
		heatMap: heatMap,
	}

	t.World.SetPositionType(c, TypeTown)
	heatMap[c.X][c.Y] = 0

	// grow towns
	for _, a := range t.World.GetAdjacentCoords(c) {
		p := t.World.GetPositionType(a)
		if (p == TypeLowland || p == TypeBeach) && t.World.isAdjacentToWater(a) {
			t.World.SetPositionType(a, TypeTown)
			heatMap[a.X][a.Y] = 0
			town.pos = append(town.pos, a)
		}
	}

	return town
}

func (t *Terrain) generateTowns(fn func() common.Coordinates) {
	t.Logger.Info(fmt.Sprintf("Initializing %v towns", common.TotalTowns))
	for i := 0; i < common.TotalTowns; i++ {
		for {
			c := fn()
			if c.X > 1 && c.Y > 1 &&
				c.X < common.WorldWidth-1 && c.Y < common.WorldHeight &&
				t.World.GetPositionType(c) == TypeBeach {

				if t.World.isAdjacentToWater(c) {
					town := t.CreateTown(c)
					if t.GenerateTownHeatMap(&town) {
						t.Logger.Info(fmt.Sprintf("[%v] Town created at %v", town.id, c))
						t.Towns = append(t.Towns, town)
						break
					} else {
						t.MakeGhostTown(&town)
					}
				}
			}
		}
	}
}

func (t *Terrain) GenerateTowns() {
	t.generateTowns(t.RandomPosition)
}

func (t *Terrain) GetRandomTown() Town {
	return t.Towns[rand.Intn(len(t.Towns))]
}
