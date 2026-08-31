package main

import (
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/harbor"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/world"
)

type testBoatWorld map[int]bool

func (s testBoatWorld) IsPassableByBoat(pos common.Coordinates) bool {
	return s[common.CoordToKey(pos)]
}

// testHarborDock reports a single cell as green parking water.
type testHarborDock map[int]bool

func (s testHarborDock) IsDock(pos common.Coordinates) bool {
	return s[common.CoordToKey(pos)]
}

func testGameStateAt(pos common.Coordinates, towns *town.Towns) *GameState {
	avatar := entities.CreateAvatar(pos, common.ShipWhite, entities.ColorPossibilities[0])
	return &GameState{
		player: &avatar,
		towns:  towns,
	}
}

func TestOpenDockFromAdjacentWater(t *testing.T) {
	cfg := economy.DefaultConfig()
	townPos := common.Coordinates{X: 20, Y: 20}
	waterPos := common.Coordinates{X: 19, Y: 20}
	bw := testBoatWorld{common.CoordToKey(waterPos): true}
	towns := town.TestTownsWith(town.NewTownForTest(townPos, cfg))

	ViewType = world.ViewTypeMainMap
	gs := testGameStateAt(waterPos, towns)

	if !openDockIfAdjacent(gs, waterPos, bw, towns, nil) {
		t.Fatal("expected dock to open from adjacent water")
	}
	if ViewType != world.ViewTypeDock {
		t.Fatalf("ViewType = %d, want ViewTypeDock", ViewType)
	}
	if gs.dockTown == nil {
		t.Fatal("dockTown should be set")
	}
}

func TestDockKeyMapped(t *testing.T) {
	found := false
	for _, item := range sailingKeyMap() {
		for _, k := range item.key {
			if k == "Enter" || k == "O" {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("sailing key map should include Enter or O for dock")
	}
}

func TestCannotDockFromLand(t *testing.T) {
	cfg := economy.DefaultConfig()
	townPos := common.Coordinates{X: 20, Y: 20}
	landPos := common.Coordinates{X: 21, Y: 20}
	bw := testBoatWorld{
		common.CoordToKey(landPos): false,
	}
	towns := town.TestTownsWith(town.NewTownForTest(townPos, cfg))

	ViewType = world.ViewTypeMainMap
	gs := testGameStateAt(landPos, towns)

	if openDockIfAdjacent(gs, landPos, bw, towns, nil) {
		t.Fatal("should not open dock from land tile")
	}
	if ViewType != world.ViewTypeMainMap {
		t.Fatalf("ViewType = %d, want main map", ViewType)
	}
	if gs.dockTown != nil {
		t.Fatal("dockTown should remain nil on land")
	}
}

// TestDockOnGreenParkingWater is the harbor rule: the ship docks when it is ON
// green, with no adjacent town required.
func TestDockOnGreenParkingWater(t *testing.T) {
	cfg := economy.DefaultConfig()
	greenPos := harbor.TownPos
	towns := town.TestHarborTowns(town.NewTownForTest(greenPos, cfg))
	bw := testBoatWorld{common.CoordToKey(greenPos): true}
	green := testHarborDock{common.CoordToKey(greenPos): true}

	ViewType = world.ViewTypeMainMap
	gs := testGameStateAt(greenPos, towns)

	if !openDockIfAdjacent(gs, greenPos, bw, towns, green) {
		t.Fatal("expected dock to open on green parking water")
	}
	if gs.dockTown == nil {
		t.Fatal("dockTown should be the harbor settlement")
	}
}

// TestNoDockOnHarborBlueWater keeps adjacent-blue from opening the dock inside the
// painted harbor: only green parking water docks.
func TestNoDockOnHarborBlueWater(t *testing.T) {
	cfg := economy.DefaultConfig()
	greenPos := harbor.TownPos
	bluePos := common.Coordinates{X: greenPos.X + 4, Y: greenPos.Y + 4}
	towns := town.TestHarborTowns(town.NewTownForTest(greenPos, cfg))
	bw := testBoatWorld{common.CoordToKey(bluePos): true}
	green := testHarborDock{common.CoordToKey(greenPos): true}

	ViewType = world.ViewTypeMainMap
	gs := testGameStateAt(bluePos, towns)

	if openDockIfAdjacent(gs, bluePos, bw, towns, green) {
		t.Fatal("blue harbor water should not open the dock")
	}
	if gs.dockTown != nil {
		t.Fatal("dockTown should remain nil on blue water")
	}
}
