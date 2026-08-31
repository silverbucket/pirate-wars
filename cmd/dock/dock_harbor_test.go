package dock

import (
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/harbor"
	"pirate-wars/cmd/town"
)

type stubHarborDock struct {
	dock map[int]bool
}

func (s stubHarborDock) IsDock(pos common.Coordinates) bool {
	return s.dock[common.CoordToKey(pos)]
}

func TestCanDockOnHarborGreen(t *testing.T) {
	cfg := economy.DefaultConfig()
	greenPos := harbor.TownPos
	openPos := common.Coordinates{X: greenPos.X + 5, Y: greenPos.Y + 5}

	world := stubWorld{}
	ht := town.NewTownForTest(harbor.TownPos, cfg)
	ht.SetTerrainType(common.TerrainTypeShallowWater)
	towns := town.TestHarborTowns(ht)
	harborDock := stubHarborDock{
		dock: map[int]bool{
			common.CoordToKey(greenPos): true,
		},
	}

	if CanDock(openPos, world, towns, harborDock) {
		t.Fatal("should not dock off green")
	}
	if !CanDock(greenPos, world, towns, harborDock) {
		t.Fatal("should dock on green parking water")
	}
}
