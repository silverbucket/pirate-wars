package main

import (
	"os"
	"testing"

	"pirate-wars/cmd/world"

	"fyne.io/fyne/v2/test"
)

// TestMain installs a headless Fyne app so widgets can be measured and laid out.
func TestMain(m *testing.M) {
	test.NewApp()
	os.Exit(m.Run())
}

// TestCreateActionMenuBeforeActionMenuAssigned guards the startup crash: createActionMenu
// runs before the global ActionMenu exists, so the bar widgets must still get built.
func TestCreateActionMenuBeforeActionMenuAssigned(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	gs := mainMapGameState()

	menu := gs.createActionMenu()
	if menu == nil {
		t.Fatal("createActionMenu returned nil")
	}
	if gs.actionBarItems == nil {
		t.Fatal("action bar widgets should be built even when ActionMenu is still nil")
	}
}

func TestUpdateActionBarIfNeededBuildsWithoutActionMenu(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	gs := mainMapGameState()

	gs.updateActionBarIfNeeded()
	if gs.actionBarItems == nil {
		t.Fatal("action bar widgets should be built when ActionMenu is nil")
	}
}
