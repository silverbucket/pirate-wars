package main

import (
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/world"

	"fyne.io/fyne/v2/container"
)

type testBoatWorld map[int]bool

func (s testBoatWorld) IsPassableByBoat(pos common.Coordinates) bool {
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

	if !openDockIfAdjacent(gs, waterPos, bw, towns) {
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

	if openDockIfAdjacent(gs, landPos, bw, towns) {
		t.Fatal("should not open dock from land tile")
	}
	if ViewType != world.ViewTypeMainMap {
		t.Fatalf("ViewType = %d, want main map", ViewType)
	}
	if gs.dockTown != nil {
		t.Fatal("dockTown should remain nil on land")
	}
}

func TestActionBarSignatureStableAcrossTicks(t *testing.T) {
	ViewType = world.ViewTypeDock
	gs := &GameState{}

	sig := gs.actionBarSignature()
	for i := 0; i < 10; i++ {
		if got := gs.actionBarSignature(); got != sig {
			t.Fatalf("tick %d: signature changed %q -> %q without state change", i, sig, got)
		}
	}
}

func TestActionBarNotRebuiltWhenSignatureUnchanged(t *testing.T) {
	ViewType = world.ViewTypeDock
	gs := &GameState{
		actionBarSig:   "dock",
		actionBarItems: container.NewHBox(),
	}
	ActionMenu = container.NewStack(container.NewHBox(), gs.actionBarItems)
	first := gs.actionBarItems

	gs.updateActionBarIfNeeded()
	if gs.actionBarItems != first {
		t.Fatal("action bar widget should be reused when signature is unchanged")
	}
}
