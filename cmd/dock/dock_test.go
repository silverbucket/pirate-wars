package dock

import (
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/town"
)

type stubWorld map[int]bool

func (s stubWorld) IsPassableByBoat(pos common.Coordinates) bool {
	return s[common.CoordToKey(pos)]
}

func TestCanDockDelegatesToAdjacentTown(t *testing.T) {
	cfg := economy.DefaultConfig()
	townPos := common.Coordinates{X: 20, Y: 20}
	waterPos := common.Coordinates{X: 19, Y: 20}
	landPos := common.Coordinates{X: 21, Y: 20}

	world := stubWorld{
		common.CoordToKey(waterPos): true,
		common.CoordToKey(landPos):  false,
	}
	towns := town.TestTownsWith(town.NewTownForTest(townPos, cfg))

	if CanDock(landPos, world, towns) {
		t.Fatal("should not dock from land")
	}
	if !CanDock(waterPos, world, towns) {
		t.Fatal("should dock from adjacent water")
	}
}
