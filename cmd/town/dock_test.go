package town

import (
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
)

type stubBoatWorld map[int]bool

func (s stubBoatWorld) IsPassableByBoat(pos common.Coordinates) bool {
	return s[common.CoordToKey(pos)]
}

func TestAdjacentTownWaterOnly(t *testing.T) {
	cfg := economy.DefaultConfig()
	townPos := common.Coordinates{X: 20, Y: 20}
	waterPos := common.Coordinates{X: 19, Y: 20}
	landPos := common.Coordinates{X: 21, Y: 20}
	farWater := common.Coordinates{X: 5, Y: 5}

	world := stubBoatWorld{
		common.CoordToKey(waterPos):   true,
		common.CoordToKey(landPos):    false,
		common.CoordToKey(farWater):   true,
	}
	towns := TestTownsWith(NewTownForTest(townPos, cfg))

	if towns.AdjacentTown(landPos, world) != nil {
		t.Fatal("should not find town from non-boat tile")
	}
	if towns.AdjacentTown(waterPos, world) == nil {
		t.Fatal("should find town from adjacent water")
	}
	if towns.AdjacentTown(farWater, world) != nil {
		t.Fatal("should not find town from distant water")
	}
}
