package resources

import (
	"pirate-wars/cmd/common"
	"testing"
)

func TestGetShipTileFallsBackToNorth(t *testing.T) {
	north := GetShipTile(common.ShipWhite, common.FacingN)
	if north == nil {
		t.Fatal("north ship tile should not be nil")
	}
	if isTileNearlyEmpty(north) {
		t.Fatal("north ship tile should not be empty")
	}

	for _, facing := range []common.Facing{
		common.FacingNE,
		common.FacingE,
		common.FacingSE,
		common.FacingS,
		common.FacingSW,
		common.FacingW,
		common.FacingNW,
	} {
		got := GetShipTile(common.ShipWhite, facing)
		if got == nil {
			t.Fatalf("ship tile for facing %v should not be nil", facing)
		}
		if got != north {
			t.Fatalf("ship tile for facing %v should fall back to north tile", facing)
		}
	}
}

func TestGetShipTilePerFaction(t *testing.T) {
	ships := []common.ShipType{
		common.ShipWhite,
		common.ShipPirate,
		common.ShipRed,
		common.ShipGreen,
		common.ShipBlue,
		common.ShipYellow,
	}

	for _, ship := range ships {
		tile := GetShipTile(ship, common.FacingN)
		if tile == nil {
			t.Fatalf("ship tile for %v should not be nil", ship)
		}
		if isTileNearlyEmpty(tile) {
			t.Fatalf("ship tile for %v should not be empty", ship)
		}
	}
}
