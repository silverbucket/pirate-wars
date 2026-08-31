package main

import (
	"testing"

	"pirate-wars/cmd/harbor"
	"pirate-wars/cmd/world"

	"go.uber.org/zap"
)

// TestInitGameStateNoHarborAssets is the headless spawn smoke path: without the
// harbor PNGs the harbor world must stay nil and the player must still spawn.
// This is the regression test for the nil *harbor.World crash in ClampSpawn.
func TestInitGameStateNoHarborAssets(t *testing.T) {
	if _, err := harbor.LoadAssets(""); err == nil {
		t.Skip("harbor assets installed; this test covers the missing-asset path")
	}

	logger, _ := zap.NewDevelopment()
	gs := initGameState(logger.Sugar())

	if gs.harborWorld != nil {
		t.Fatal("harbor world should be nil when assets are absent")
	}
	if gs.harborRenderer != nil {
		t.Fatal("harbor renderer should be nil when assets are absent")
	}
	if gs.player == nil {
		t.Fatal("player was not created")
	}
	if !gs.world.IsPassableByBoat(gs.player.GetPos()) {
		t.Fatalf("player spawned on impassable cell %+v", gs.player.GetPos())
	}
}

// TestNilHarborWorldIsNilSafe locks the typed-nil interface hazard: a nil
// *harbor.World stored in an interface is not == nil, so its methods must not
// dereference the receiver.
func TestNilHarborWorldIsNilSafe(t *testing.T) {
	var hw *harbor.World
	if _, ok := hw.ClampSpawn(harbor.TownPos); ok {
		t.Fatal("nil harbor world should not report a spawn")
	}
	if hw.IsDock(harbor.TownPos) {
		t.Fatal("nil harbor world should not report a dock")
	}
	if hw.IsPassableByBoat(harbor.TownPos) {
		t.Fatal("nil harbor world should not report passable water")
	}
}

// TestCombinedWorldWithoutHarbor checks sailing passability falls through to the
// tilemap when the harbor world is missing, including inside the harbor rect.
func TestCombinedWorldWithoutHarbor(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	cw := newCombinedWorld(world.Init(logger.Sugar()), nil)

	if cw.IsDock(harbor.TownPos) {
		t.Fatal("no harbor mask means no dock-on-green")
	}
	// Must not panic; result depends on generated terrain.
	_ = cw.IsPassableByBoat(harbor.TownPos)
}
