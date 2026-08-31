package main

import (
	"strings"
	"testing"

	"pirate-wars/cmd/world"

	"go.uber.org/zap"
)

// TestBumpingAShipOffersAHailRatherThanForcingOne covers a review finding.
//
// Sailing into an occupied cell used to force the hail modal open. Movement is
// accumulator-driven, so the player cannot choose to stop short of a ship, and
// dismissing the modal re-triggered on the next step — about every 400ms at full
// sail — until the player happened to turn away.
func TestBumpingAShipOffersAHailRatherThanForcingOne(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())
	ViewType = world.ViewTypeMainMap

	npcs := gs.npcs.GetList()
	if len(npcs) == 0 {
		t.Skip("generated world has no NPCs")
	}
	target := npcs[0]

	gs.alongsideNpcID = target.GetID()
	gs.refreshActionBar()

	if ViewType != world.ViewTypeMainMap {
		t.Fatal("coming alongside must not open a modal by itself")
	}

	var hailBtn *button
	for i := range gs.buttons {
		if strings.HasPrefix(gs.buttons[i].label, "Hail") {
			hailBtn = &gs.buttons[i]
		}
	}
	if hailBtn == nil {
		t.Fatalf("Hail should be offered on the bar, got %v", gs.actionBarLabels())
	}
	if !hailBtn.enabled {
		t.Fatal("Hail should be enabled with a ship alongside")
	}

	gs.handleKeyPress("H")
	if ViewType != world.ViewTypeHail {
		t.Fatalf("H should open the hail, ViewType = %d", ViewType)
	}
}

// TestHailIsDisabledWithNoShipAlongside checks the key explains itself rather
// than doing nothing.
func TestHailIsDisabledWithNoShipAlongside(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	gs := mainMapGameState()

	gs.handleKeyPress("H")
	if ViewType == world.ViewTypeHail {
		t.Fatal("H should not open a hail with nothing alongside")
	}
	if got := gs.activeNotice(); !strings.Contains(got, "alongside") {
		t.Fatalf("H with nothing alongside should say so, got %q", got)
	}
}

// TestHailPromptDoesNotRepeatForTheSameShip covers the re-fire loop directly.
func TestHailPromptDoesNotRepeatForTheSameShip(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())
	ViewType = world.ViewTypeMainMap

	npcs := gs.npcs.GetList()
	if len(npcs) == 0 {
		t.Skip("generated world has no NPCs")
	}
	gs.alongsideNpcID = npcs[0].GetID()

	gs.handleKeyPress("H")
	if ViewType != world.ViewTypeHail {
		t.Fatalf("H should open the hail, ViewType = %d", ViewType)
	}
	gs.handleKeyPress("Escape")

	if gs.lastHailedNpcID != npcs[0].GetID() {
		t.Fatal("a hailed ship should be remembered so the prompt does not repeat")
	}
}
