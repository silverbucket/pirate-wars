package town

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/resources"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"

	"go.uber.org/zap"
)

type Towns struct {
	logger *zap.SugaredLogger
	list   []Town
}

type Town struct {
	id          string
	pos         []common.Coordinates
	terrainType common.TerrainType
	logger      *zap.SugaredLogger
	color       color.Color
	HeatMap     HeatMap
	market      economy.TownMarket
	blink       bool
	alternate   bool
}

var townList []Town

func (t *Town) GetID() string {
	return t.id
}

func (t *Town) GetPos() common.Coordinates {
	return t.pos[0]
}

func (t *Town) GetPositions() []common.Coordinates {
	return t.pos
}

func (t *Town) Market() *economy.TownMarket {
	return &t.market
}

func (t *Town) GetPreviousPos() common.Coordinates {
	return t.pos[0]
}

func (t *Town) GetTerrainType() common.TerrainType {
	return t.terrainType
}

func (t *Town) GetType() string {
	return "Town"
}

func (t *Town) GetViewableRange() window.Dimensions {
	return window.Dimensions{Width: 20, Height: 20}
}

func (t *Town) Highlight(b bool) {
	t.blink = b
	t.alternate = b

}

func (t *Town) IsHighlighted() bool {
	return t.blink
}

func (t *Town) SetTerrainType(tt common.TerrainType) {
	t.terrainType = tt
}

func (t *Town) GetName() string {
	return t.id
}

func (t *Town) GetColor() color.Color {
	if t.blink {
		if !t.alternate {
			t.alternate = true
			return color.RGBA{0, 0, 0, 0}
		}
	}
	t.alternate = false
	return t.color
}

func (t *Town) GetTileImage() image.Image {
	return resources.GetTerrainTile(t.terrainType)
}

func (t *Town) GetFlag() string {
	return "NA"
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
		t.SetTerrainType(common.TerrainTypeGhostTown)
		world.SetPositionType(c, common.TerrainTypeGhostTown)
	}
}

func (ts *Towns) CreateTown(c common.Coordinates, world *world.MapView, cfg economy.Config) Town {
	var heatMap = make([][]HeatMapCost, common.WorldRows)

	for i := range heatMap {
		heatMap[i] = make([]HeatMapCost, common.WorldCols)
		for j := range heatMap[i] {
			heatMap[i][j] = -1
		}
	}

	town := Town{
		id:          common.GenID(c),
		pos:         []common.Coordinates{c},
		terrainType: common.TerrainTypeTown,
		logger:      ts.logger,
		color:       color.RGBA{189, 55, 31, 255},
		HeatMap: HeatMap{
			grid: heatMap,
		},
		market: economy.NewTownMarket(cfg, rand.New(rand.NewSource(int64(c.X*1000+c.Y)))),
	}

	world.SetPositionType(c, common.TerrainTypeTown)
	heatMap[c.X][c.Y] = 0

	// grow towns
	for _, a := range world.GetAdjacentCoords(c) {
		p := world.GetPositionType(a)
		if (p == common.TerrainTypeLowland || p == common.TerrainTypeBeach) && world.IsAdjacentToWater(a) {
			world.SetPositionType(a, common.TerrainTypeTown)
			//HeatMap[a.X][a.Y] = 0
			town.pos = append(town.pos, a)
		}
	}
	world.SetMapItem(&town)
	return town
}

func (ts *Towns) initializeTowns(fn func() common.Coordinates, world *world.MapView, cfg economy.Config) []Town {
	ts.logger.Info(fmt.Sprintf("Initializing %v towns", common.TotalTowns))
	for i := 0; i < common.TotalTowns; i++ {
		for {
			c := fn()
			if c.X > 1 && c.Y > 1 &&
				c.X < common.WorldCols-1 && c.Y < common.WorldRows &&
				world.GetPositionType(c) == common.TerrainTypeBeach {

				if world.IsAdjacentToWater(c) {
					town := ts.CreateTown(c, world, cfg)
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

func Init(world *world.MapView, logger *zap.SugaredLogger, cfg economy.Config) *Towns {
	ts := Towns{
		logger: logger,
		list:   []Town{},
	}
	ts.list = ts.initializeTowns(common.RandomPosition, world, cfg)
	ts.logger.Info(fmt.Sprintf("Created %v towns", len(ts.list)))
	return &ts
}

func (ts *Towns) GetByID(id string) *Town {
	for i := range ts.list {
		if ts.list[i].GetID() == id {
			return &ts.list[i]
		}
	}
	return nil
}

func (ts *Towns) AdjacentTown(pos common.Coordinates, world interface{ IsPassableByBoat(common.Coordinates) bool }) *Town {
	if world == nil || !world.IsPassableByBoat(pos) {
		return nil
	}
	for i := range ts.list {
		for _, tp := range ts.list[i].GetPositions() {
			if common.IsPositionAdjacent(pos, tp) {
				return &ts.list[i]
			}
		}
	}
	return nil
}

func (ts *Towns) GetRandomTown() (Town, error) {
	if len(ts.list) == 0 {
		return Town{}, errors.New("no towns found")
	}
	return ts.list[rand.Intn(len(ts.list))], nil
}

func (ts *Towns) GetTowns() []Town {
	return ts.list
}

// NewTownForTest builds a minimal town for unit tests.
func NewTownForTest(pos common.Coordinates, cfg economy.Config) Town {
	return Town{
		id:          common.GenID(pos),
		pos:         []common.Coordinates{pos},
		terrainType: common.TerrainTypeTown,
		color:       color.RGBA{189, 55, 31, 255},
		market:      economy.NewTownMarket(cfg, rand.New(rand.NewSource(int64(pos.X*1000+pos.Y)))),
	}
}

// TestTownsWith returns a Towns collection for unit tests.
func TestTownsWith(towns ...Town) *Towns {
	return &Towns{list: towns}
}

func (ts *Towns) GetVisible(playerPos common.Coordinates) []Town {
	vp := window.GetViewportRegion(playerPos)
	visible := []Town{}
	for _, town := range ts.list {
		if vp.IsPositionWithin(town.GetPos()) {
			visible = append(visible, town)
		}
	}
	return visible
}
