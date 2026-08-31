package main

import (
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

func actionBarButtonLabels(bar *fyne.Container) []string {
	labels := []string{}
	walkActionBar(bar, func(obj fyne.CanvasObject) {
		if btn, ok := obj.(*widget.Button); ok {
			labels = append(labels, btn.Text)
		}
	})
	return labels
}

func actionBarText(bar *fyne.Container) string {
	parts := []string{}
	walkActionBar(bar, func(obj fyne.CanvasObject) {
		switch w := obj.(type) {
		case *widget.Button:
			parts = append(parts, w.Text)
		case *widget.Label:
			parts = append(parts, w.Text)
		case *canvas.Text:
			parts = append(parts, w.Text)
		}
	})
	return strings.Join(parts, " | ")
}

func walkActionBar(obj fyne.CanvasObject, visit func(fyne.CanvasObject)) {
	if c, ok := obj.(*fyne.Container); ok {
		for _, child := range c.Objects {
			walkActionBar(child, visit)
		}
		return
	}
	visit(obj)
}

func mainMapGameState() *GameState {
	avatar := entities.CreateAvatar(common.Coordinates{X: 5, Y: 5}, common.ShipWhite, entities.ColorPossibilities[0])
	return &GameState{player: &avatar}
}

func TestSailingActionBarLabelsIncludeKeys(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	gs := mainMapGameState()

	labels := actionBarButtonLabels(gs.ActionItems())
	want := []string{
		"Full sail (1)",
		"Half sail (2)",
		"Furled (3)",
		"Cycle sail (V)",
		"Examine (X)",
		"Minimap (M)",
	}
	for _, w := range want {
		found := false
		for _, got := range labels {
			if got == w {
				found = true
			}
		}
		if !found {
			t.Fatalf("sailing action bar missing %q, got %v", w, labels)
		}
	}

	for _, label := range labels {
		if !strings.Contains(label, "(") {
			t.Fatalf("action bar button %q does not show its key", label)
		}
	}
}

func TestActionBarHasNoDuplicateCommands(t *testing.T) {
	for _, view := range []int{
		world.ViewTypeMainMap,
		world.ViewTypeDock,
		world.ViewTypeHail,
		world.ViewTypeMiniMap,
	} {
		ViewType = view
		ActionMenu = nil
		gs := mainMapGameState()

		seen := map[string]int{}
		for _, label := range actionBarButtonLabels(gs.ActionItems()) {
			action := strings.SplitN(label, " (", 2)[0]
			seen[action]++
			if seen[action] > 1 {
				t.Fatalf("view %d: duplicate button for %q", view, action)
			}
		}
	}
}

func TestDockButtonHiddenWhenNotAdjacent(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	gs := mainMapGameState()

	for _, label := range actionBarButtonLabels(gs.ActionItems()) {
		if strings.HasPrefix(label, "Dock") {
			t.Fatalf("dock button should be hidden away from a town, got %q", label)
		}
	}
}

func TestDockCommandLabelShowsBothKeys(t *testing.T) {
	sailing := sailingKeyMap()
	var dockItem *keyItem
	for i := range sailing {
		if sailing[i].label == "Dock" {
			dockItem = &sailing[i]
		}
	}
	if dockItem == nil {
		t.Fatal("sailing key map should contain a Dock command")
	}
	if got := dockItem.barLabel(); got != "Dock (Enter/O)" {
		t.Fatalf("dock bar label = %q, want Dock (Enter/O)", got)
	}

	keys := map[string]bool{}
	for _, k := range dockItem.key {
		keys[k] = true
	}
	if !keys["Enter"] || !keys["O"] {
		t.Fatalf("dock should be bound to Enter and O, got %v", dockItem.key)
	}
	if dockItem.barVisible == nil {
		t.Fatal("dock button should be gated on town adjacency")
	}
	if dockItem.barVisible(mainMapGameState()) {
		t.Fatal("dock button should not be visible without an adjacent town")
	}
}

// TestHeadingKeysVisibleOnBar keeps WASD and the bound diagonals on the bar legend
// rather than hidden behind a help screen.
func TestHeadingKeysVisibleOnBar(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	gs := mainMapGameState()

	text := actionBarText(gs.ActionItems())
	if !strings.Contains(text, "Heading") {
		t.Fatalf("action bar should show a heading legend, got %q", text)
	}
	for _, key := range []string{"W", "A", "S", "D", "Q", "E", "Z", "C"} {
		if !strings.Contains(text, key) {
			t.Fatalf("heading legend should show key %q, got %q", key, text)
		}
	}
	for _, dir := range []string{"N", "S", "E", "W", "NW", "NE", "SW", "SE"} {
		if !strings.Contains(text, dir) {
			t.Fatalf("heading legend should show direction %q, got %q", dir, text)
		}
	}
}

func TestAdminKeysVisibleOnBar(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	gs := mainMapGameState()

	text := actionBarText(gs.ActionItems())
	for _, want := range []string{"Help (?)", "Quit (Ctrl+Q)"} {
		if !strings.Contains(text, want) {
			t.Fatalf("action bar should list %q, got %q", want, text)
		}
	}
}

func TestOverlayActionBarLabels(t *testing.T) {
	cases := []struct {
		view int
		want string
	}{
		{world.ViewTypeDock, "Leave dock (Esc)"},
		{world.ViewTypeHail, "Dismiss (Enter/Esc/X)"},
		{world.ViewTypeMiniMap, "Exit minimap (M/Enter)"},
	}
	for _, c := range cases {
		ViewType = c.view
		ActionMenu = nil
		gs := mainMapGameState()

		labels := actionBarButtonLabels(gs.ActionItems())
		found := false
		for _, label := range labels {
			if label == c.want {
				found = true
			}
		}
		if !found {
			t.Fatalf("view %d: expected button %q, got %v", c.view, c.want, labels)
		}
	}
}

// TestBarLabelListsEveryBinding covers commands bound to more than one key: the bar
// must name them all rather than only the first.
func TestBarLabelListsEveryBinding(t *testing.T) {
	if got := barLabelFor(miniMapKeyMap(), "Exit minimap"); got != "Exit minimap (M/Enter)" {
		t.Fatalf("exit minimap label = %q, want Exit minimap (M/Enter)", got)
	}

	ViewType = world.ViewTypeMiniMap
	ActionMenu = nil
	gs := mainMapGameState()

	var label string
	for _, l := range actionBarButtonLabels(gs.ActionItems()) {
		if strings.HasPrefix(l, "Exit minimap") {
			label = l
		}
	}
	if !strings.Contains(label, "M") || !strings.Contains(label, "Enter") {
		t.Fatalf("exit minimap button %q should show both M and Enter", label)
	}
}

func TestExamineActionBarLabelsIncludeKeys(t *testing.T) {
	ViewType = world.ViewTypeExamine
	ActionMenu = nil
	gs := mainMapGameState()

	labels := actionBarButtonLabels(gs.ActionItems())
	for _, want := range []string{"Exit examine (X/Enter)", "Examine left (←/H/A)", "Examine right (→/L/D)"} {
		found := false
		for _, got := range labels {
			if got == want {
				found = true
			}
		}
		if !found {
			t.Fatalf("examine bar missing %q, got %v", want, labels)
		}
	}
}

// TestActionBarFitsActionMenuArea keeps the cheat-sheet bar inside its allotted space
// so no button or key hint gets clipped.
func TestActionBarFitsActionMenuArea(t *testing.T) {
	for _, view := range []int{
		world.ViewTypeMainMap,
		world.ViewTypeDock,
		world.ViewTypeHail,
		world.ViewTypeMiniMap,
		world.ViewTypeExamine,
	} {
		ViewType = view
		ActionMenu = nil
		gs := mainMapGameState()

		min := gs.ActionItems().MinSize()
		if min.Width > float32(window.ActionMenu.Width) {
			t.Fatalf("view %d: action bar width %.0f exceeds area %d", view, min.Width, window.ActionMenu.Width)
		}
		if min.Height > float32(window.ActionMenu.Height) {
			t.Fatalf("view %d: action bar height %.0f exceeds area %d", view, min.Height, window.ActionMenu.Height)
		}
	}
}

// TestDockButtonAndKeyOnGeneratedWorld exercises the real world: standing on water
// beside a town must show exactly one Dock button, keep the bar within its area, and
// let the Enter/O binding open the dock overlay.
func TestDockButtonAndKeyOnGeneratedWorld(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())

	dockPos, ok := waterTileBesideTown(gs)
	if !ok {
		t.Skip("generated world has no water-adjacent town")
	}
	gs.player.SetPos(dockPos)

	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	gs.createActionMenu()

	labels := actionBarButtonLabels(gs.actionBarItems)
	dockButtons := 0
	for _, label := range labels {
		if strings.HasPrefix(label, "Dock") {
			if label != "Dock (Enter/O)" {
				t.Fatalf("dock button label = %q, want Dock (Enter/O)", label)
			}
			dockButtons++
		}
	}
	if dockButtons != 1 {
		t.Fatalf("want exactly one dock button, got %d in %v", dockButtons, labels)
	}
	if w := gs.actionBarItems.MinSize().Width; w > float32(window.ActionMenu.Width) {
		t.Fatalf("action bar width %.0f exceeds area %d with dock button", w, window.ActionMenu.Width)
	}

	gs.window = test.NewWindow(nil)
	gs.buildOverlayShell()
	gs.handleKeyPress(&fyne.KeyEvent{Name: "Enter"})
	if ViewType != world.ViewTypeDock {
		t.Fatalf("dock key did not open the dock overlay, ViewType = %d", ViewType)
	}
	if gs.dockTown == nil {
		t.Fatal("dock key did not set dockTown on the game state it was called with")
	}
}

// TestDockKeyUsesCallbackGameState checks the Enter/O handler acts on the GameState it
// is invoked with rather than a package-level reference.
func TestDockKeyUsesCallbackGameState(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())
	gs.window = test.NewWindow(nil)
	gs.buildOverlayShell()

	dockPos, ok := waterTileBesideTown(gs)
	if !ok {
		t.Skip("generated world has no water-adjacent town")
	}
	gs.player.SetPos(dockPos)

	other := mainMapGameState()
	ViewType = world.ViewTypeMainMap
	ActionMenu = nil

	for _, key := range []string{"Enter", "O"} {
		ViewType = world.ViewTypeMainMap
		gs.dockTown = nil
		other.dockTown = nil

		gs.handleKeyPress(&fyne.KeyEvent{Name: fyne.KeyName(key)})

		if gs.dockTown == nil {
			t.Fatalf("%s did not dock the game state it was called with", key)
		}
		if other.dockTown != nil {
			t.Fatalf("%s docked an unrelated game state", key)
		}
	}
}

// TestKeyboardViewChangeUpdatesActionBar covers the bar going stale until the next tick
// after a key changes the view.
func TestKeyboardViewChangeUpdatesActionBar(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())
	gs.window = test.NewWindow(nil)
	gs.buildOverlayShell()

	dockPos, ok := waterTileBesideTown(gs)
	if !ok {
		t.Skip("generated world has no water-adjacent town")
	}
	gs.player.SetPos(dockPos)

	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	ActionMenu = gs.createActionMenu()

	gs.handleKeyPress(&fyne.KeyEvent{Name: "Enter"})
	if ViewType != world.ViewTypeDock {
		t.Fatalf("Enter did not open dock, ViewType = %d", ViewType)
	}
	labels := actionBarButtonLabels(gs.actionBarItems)
	if len(labels) != 1 || labels[0] != "Leave dock (Esc)" {
		t.Fatalf("bar should show dock commands right after the key press, got %v", labels)
	}
	if ActionMenu.Objects[1] != gs.actionBarItems {
		t.Fatal("rebuilt bar was not mounted into the action menu")
	}

	gs.handleKeyPress(&fyne.KeyEvent{Name: "Escape"})
	if ViewType != world.ViewTypeMainMap {
		t.Fatalf("Escape did not leave dock, ViewType = %d", ViewType)
	}
	labels = actionBarButtonLabels(gs.actionBarItems)
	if len(labels) == 0 || !strings.HasPrefix(labels[0], "Dock") {
		t.Fatalf("bar should show sailing commands right after leaving dock, got %v", labels)
	}
	if gs.overlayRoot != nil && !gs.overlayRoot.Hidden {
		t.Fatal("overlay should be hidden after leaving dock with a key")
	}
}

// TestActionBarWidgetsReusedAcrossKeyPresses keeps the signature cache doing its job:
// repeated key presses that do not change the commands must not rebuild the bar.
func TestActionBarWidgetsReusedAcrossKeyPresses(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())
	gs.window = test.NewWindow(nil)
	gs.buildOverlayShell()

	ViewType = world.ViewTypeMainMap
	ActionMenu = nil
	ActionMenu = gs.createActionMenu()
	first := gs.actionBarItems

	for i := 0; i < 5; i++ {
		gs.handleKeyPress(&fyne.KeyEvent{Name: "1"})
		if gs.actionBarItems != first {
			t.Fatalf("press %d rebuilt the action bar despite unchanged commands", i)
		}
	}
}

func waterTileBesideTown(gs *GameState) (common.Coordinates, bool) {
	for _, tw := range gs.towns.GetTowns() {
		for _, tp := range tw.GetPositions() {
			for _, adj := range gs.world.GetAdjacentCoords(tp) {
				if gs.world.IsPassableByBoat(adj) {
					return adj, true
				}
			}
		}
	}
	return common.Coordinates{}, false
}

// TestOverlayButtonsMatchBarLabels keeps overlay buttons in sync with the key maps
// so a key hint never drifts from its binding.
func TestOverlayButtonsMatchBarLabels(t *testing.T) {
	if got := barLabelFor(dockKeyMap(), "Leave dock"); got != "Leave dock (Esc)" {
		t.Fatalf("dock overlay label = %q", got)
	}
	if got := barLabelFor(hailKeyMap(), "Dismiss"); got != "Dismiss (Enter/Esc/X)" {
		t.Fatalf("hail overlay label = %q", got)
	}
	if got := barLabelFor(dockKeyMap(), "Back"); got != "Back" {
		t.Fatalf("unbound action should keep a plain label, got %q", got)
	}
}
