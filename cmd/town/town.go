package town

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/terrain"
	"pirate-wars/cmd/world"
)

type Towns struct {
	logger *zap.SugaredLogger
	list   []Town
}

type Town struct {
	id          string
	pos         []common.Coordinates
	terrainType terrain.Type
	logger      *zap.SugaredLogger
	HeatMap     HeatMap
}

var townList []Town

func (t *Town) GetID() string {
	return t.id
}

func (t *Town) GetPos() common.Coordinates {
	return t.pos[0]
}

func (t *Town) GetTerrainType() terrain.Type {
	return t.terrainType
}

func (t *Town) SetTerrainType(tt terrain.Type) {
	t.terrainType = tt
}

func (t *Town) AccessibleFrom(c common.Coordinates) bool {
	for _, d := range common.Directions {
		n := common.AddDirection(c, d)
		if n.X < 0 || n.Y < 0 {
			return false
		}
	}
	return true
}

func (t *Town) MakeGhostTown(world *world.MapView) {
	t.logger.Info(fmt.Sprintf("[%v] Town turns to ghost town at %v", t.id, t.GetPos()))
	for _, c := range t.pos {
		t.SetTerrainType(terrain.TypeGhostTown)
		world.SetPositionType(c, terrain.TypeGhostTown)
	}
}

func (ts *Towns) CreateTown(c common.Coordinates, world *world.MapView) Town {
	var heatMap = make([][]HeatMapCost, common.WorldHeight)

	for i := range heatMap {
		heatMap[i] = make([]HeatMapCost, common.WorldWidth)
		for j := range heatMap[i] {
			heatMap[i][j] = -1
		}
	}

	town := Town{
		id:          common.GenID(c),
		pos:         []common.Coordinates{c},
		terrainType: terrain.TypeTown,
		logger:      ts.logger,
		HeatMap: HeatMap{
			grid: heatMap,
		},
	}

	world.SetPositionType(c, terrain.TypeTown)
	heatMap[c.X][c.Y] = 0

	// grow towns
	for _, a := range world.GetAdjacentCoords(c) {
		p := world.GetPositionType(a)
		if (p == terrain.TypeLowland || p == terrain.TypeBeach) && world.IsAdjacentToWater(a) {
			world.SetPositionType(a, terrain.TypeTown)
			//HeatMap[a.X][a.Y] = 0
			town.pos = append(town.pos, a)
		}
	}
	world.SetMapItem(&town)
	return town
}

func (ts *Towns) initializeTowns(fn func() common.Coordinates, world *world.MapView) []Town {
	ts.logger.Info(fmt.Sprintf("Initializing %v towns", common.TotalTowns))
	for i := 0; i < common.TotalTowns; i++ {
		for {
			c := fn()
			if c.X > 1 && c.Y > 1 &&
				c.X < common.WorldWidth-1 && c.Y < common.WorldHeight &&
				world.GetPositionType(c) == terrain.TypeBeach {

				if world.IsAdjacentToWater(c) {
					town := ts.CreateTown(c, world)
					if town.generateHeatMap(world) {
						ts.logger.Info(fmt.Sprintf("[%v] Town created at %v", town.id, c))
						townList = append(townList, town)
						break
					} else {
						town.MakeGhostTown(world)
					}
				}
			}
		}
	}
	return townList
}

func Init(world *world.MapView, logger *zap.SugaredLogger) *Towns {
	ts := Towns{
		logger: logger,
		list:   []Town{},
	}
	ts.list = ts.initializeTowns(common.RandomPosition, world)
	ts.logger.Info(fmt.Sprintf("Created %v towns", len(ts.list)))
	return &ts
}

func (ts *Towns) GetRandomTown() (Town, error) {
	if len(ts.list) == 0 {
		return Town{}, errors.New("no towns found")
	}
	return ts.list[rand.Intn(len(ts.list))], nil
}
