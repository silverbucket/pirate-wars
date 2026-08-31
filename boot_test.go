package main

import (
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/terrain"
	"pirate-wars/cmd/world"

	"go.uber.org/zap"
)

// TestBootIsTheOceanGame is the scope lock for this PR: the Ebiten port must boot
// into the same 800×800 ocean game as master — full world, towns, NPC traders,
// wind — with the player on open water and the main map as the live view.
func TestBootIsTheOceanGame(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())

	if ViewType != world.ViewTypeMainMap {
		t.Fatalf("boot view = %d, want the ocean main map", ViewType)
	}

	if w, h := gs.world.GetWidth(), gs.world.GetHeight(); w != common.WorldCols || h != common.WorldRows {
		t.Fatalf("world is %dx%d, want %dx%d", w, h, common.WorldCols, common.WorldRows)
	}

	if gs.player == nil {
		t.Fatal("no player")
	}
	if !gs.world.IsPassableByBoat(gs.player.GetPos()) {
		t.Fatalf("player spawned on impassable cell %+v", gs.player.GetPos())
	}

	if got := len(gs.towns.GetTowns()); got != common.TotalTowns {
		t.Fatalf("boot has %d towns, want the master count %d", got, common.TotalTowns)
	}
	if got := len(gs.npcs.GetList()); got != common.TotalNpcs {
		t.Fatalf("boot has %d NPC traders, want the master count %d", got, common.TotalNpcs)
	}
	if gs.wind == nil {
		t.Fatal("no wind")
	}
}

// TestBootWorldHasEveryTerrainType guards against a rendering change quietly
// dropping terrain: the generated ocean must still contain the full range from
// deep water up to peaks, plus the towns stamped onto it.
func TestBootWorldHasEveryTerrainType(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())

	seen := map[common.TerrainType]int{}
	for x := 0; x < gs.world.GetWidth(); x++ {
		for y := 0; y < gs.world.GetHeight(); y++ {
			seen[gs.world.GetPositionType(common.Coordinates{X: x, Y: y})]++
		}
	}

	want := []common.TerrainType{
		common.TerrainTypeDeepWater,
		common.TerrainTypeOpenWater,
		common.TerrainTypeShallowWater,
		common.TerrainTypeBeach,
		common.TerrainTypeLowland,
		common.TerrainTypeHighland,
		common.TerrainTypeRock,
		common.TerrainTypePeak,
		common.TerrainTypeTown,
	}
	for _, tt := range want {
		if seen[tt] == 0 {
			t.Fatalf("terrain type %d missing from the booted world", tt)
		}
		if _, ok := terrain.TypeLookup[tt]; !ok {
			t.Fatalf("terrain type %d has no tileset entry", tt)
		}
	}
}

// TestBootTicksLikeMaster runs the real loop for a few sailing ticks to catch a
// panic in the ocean path (NPC pathing, wind drift, animation).
func TestBootTicksLikeMaster(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())
	gs.paused = false
	ViewType = world.ViewTypeMainMap

	startPos := gs.player.GetPos()
	for i := 0; i < 25; i++ {
		gs.processTick()
	}

	if !gs.world.IsPassableByBoat(gs.player.GetPos()) {
		t.Fatalf("player left navigable water: %+v -> %+v", startPos, gs.player.GetPos())
	}
	for _, n := range gs.npcs.GetList() {
		if !common.Inbounds(n.GetPos()) {
			t.Fatalf("NPC %s left the map at %+v", n.GetID(), n.GetPos())
		}
	}
}
