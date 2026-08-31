package main

import (
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/gfx"
	"pirate-wars/cmd/world"

	"go.uber.org/zap"
)

func mainMapGameState() *GameState {
	avatar := entities.CreateAvatar(common.Coordinates{X: 5, Y: 5}, common.ShipWhite, entities.ColorPossibilities[0])
	return &GameState{player: &avatar}
}

// actionBarText is the whole bar as read by a player: the legend line plus every
// button label.
func actionBarText(gs *GameState) string {
	parts := append([]string{gs.actionBarCaption()}, gs.actionBarLabels()...)
	return strings.Join(parts, " | ")
}

func TestSailingActionBarLabelsIncludeKeys(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	gs := mainMapGameState()

	labels := gs.actionBarLabels()
	want := []string{
		"More sail (W)",
		"Less sail (S)",
		"Tack port (A)",
		"Tack starboard (D)",
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
		gs := mainMapGameState()

		seen := map[string]int{}
		for _, label := range gs.actionBarLabels() {
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
	gs := mainMapGameState()

	for _, label := range gs.actionBarLabels() {
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

// TestSailingControlKeysVisibleOnBar keeps the WASD steering and sail keys on the
// bar rather than hidden behind a help screen, and keeps the old compass-direction
// heading legend from coming back now that W means sail, not north.
func TestSailingControlKeysVisibleOnBar(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	gs := mainMapGameState()

	text := actionBarText(gs)
	for _, want := range []string{"Tack port (A)", "Tack starboard (D)", "More sail (W)", "Less sail (S)"} {
		if !strings.Contains(text, want) {
			t.Fatalf("action bar should show %q, got %q", want, text)
		}
	}
	if strings.Contains(text, "Heading") {
		t.Fatalf("compass heading legend should be gone, got %q", text)
	}
	for _, gone := range []string{"N:", "NE:", "SW:"} {
		if strings.Contains(text, gone) {
			t.Fatalf("action bar should not list compass heading %q, got %q", gone, text)
		}
	}
	if !strings.Contains(text, "Sail Full (1)") {
		t.Fatalf("sail presets should stay on the legend, got %q", text)
	}
}

func TestAdminKeysVisibleOnBar(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	gs := mainMapGameState()

	text := actionBarText(gs)
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
		gs := mainMapGameState()

		labels := gs.actionBarLabels()
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
	gs := mainMapGameState()

	var label string
	for _, l := range gs.actionBarLabels() {
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
	gs := mainMapGameState()

	labels := gs.actionBarLabels()
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

// TestActionBarFitsActionMenuArea keeps the cheat-sheet bar inside its allotted
// space so no button or key hint gets clipped.
func TestActionBarFitsActionMenuArea(t *testing.T) {
	for _, view := range []int{
		world.ViewTypeMainMap,
		world.ViewTypeDock,
		world.ViewTypeHail,
		world.ViewTypeMiniMap,
		world.ViewTypeExamine,
	} {
		ViewType = view
		gs := mainMapGameState()

		if w := gfx.TextWidth(gs.actionBarCaption()); w > actionBarRect.Dx() {
			t.Fatalf("view %d: legend width %d exceeds bar %d", view, w, actionBarRect.Dx())
		}

		width := buttonGap
		for _, label := range gs.actionBarLabels() {
			width += gfx.TextWidth(label) + buttonPadding*2 + buttonGap
		}
		if width > actionBarRect.Dx() {
			t.Fatalf("view %d: button row width %d exceeds bar %d", view, width, actionBarRect.Dx())
		}
	}
}

// TestActionBarButtonsAreLaidOut checks every bar command actually gets a hit
// rectangle. Ebiten hit-tests these rects each frame, which is what removed the
// Fyne widget-rebuild click bug.
func TestActionBarButtonsAreLaidOut(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	gs := mainMapGameState()
	gs.refreshActionBar()

	labels := gs.actionBarLabels()
	if len(gs.buttons) != len(labels) {
		t.Fatalf("laid out %d buttons for %d bar commands", len(gs.buttons), len(labels))
	}
	for i, b := range gs.buttons {
		if b.label != labels[i] {
			t.Fatalf("button %d label = %q, want %q", i, b.label, labels[i])
		}
		if b.rect.Dx() <= 0 || b.rect.Dy() <= 0 {
			t.Fatalf("button %q has empty rect %v", b.label, b.rect)
		}
		if !b.rect.Overlaps(actionBarRect) {
			t.Fatalf("button %q rect %v is outside the action bar %v", b.label, b.rect, actionBarRect)
		}
	}
}

// TestTapRunsButtonAction covers the click path end to end: a tap inside a button
// rect runs that command on the game state the button was built for.
func TestTapRunsButtonAction(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())

	dockPos, ok := waterTileBesideTown(gs)
	if !ok {
		t.Skip("generated world has no water-adjacent town")
	}
	gs.player.SetPos(dockPos)
	ViewType = world.ViewTypeMainMap
	gs.refreshActionBar()

	var dockBtn *button
	for i := range gs.buttons {
		if strings.HasPrefix(gs.buttons[i].label, "Dock") {
			dockBtn = &gs.buttons[i]
		}
	}
	if dockBtn == nil {
		t.Fatalf("no dock button beside a town, got %v", gs.actionBarLabels())
	}

	center := dockBtn.rect.Min.Add(dockBtn.rect.Max).Div(2)
	gs.tap(center.X, center.Y)

	if ViewType != world.ViewTypeDock {
		t.Fatalf("tapping Dock did not open the dock, ViewType = %d", ViewType)
	}
	if gs.dockTown == nil {
		t.Fatal("tapping Dock did not set dockTown on the state the button was built for")
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
	labels := gs.actionBarLabels()
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

	width := buttonGap
	for _, label := range labels {
		width += gfx.TextWidth(label) + buttonPadding*2 + buttonGap
	}
	if width > actionBarRect.Dx() {
		t.Fatalf("action bar width %d exceeds area %d with dock button", width, actionBarRect.Dx())
	}

	gs.handleKeyPress("Enter")
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

	dockPos, ok := waterTileBesideTown(gs)
	if !ok {
		t.Skip("generated world has no water-adjacent town")
	}
	gs.player.SetPos(dockPos)

	other := mainMapGameState()

	for _, key := range []string{"Enter", "O"} {
		ViewType = world.ViewTypeMainMap
		gs.dockTown = nil
		other.dockTown = nil

		gs.handleKeyPress(key)

		if gs.dockTown == nil {
			t.Fatalf("%s did not dock the game state it was called with", key)
		}
		if other.dockTown != nil {
			t.Fatalf("%s docked an unrelated game state", key)
		}
	}
}

// TestKeyboardViewChangeUpdatesActionBar covers the bar going stale until the next
// tick after a key changes the view.
func TestKeyboardViewChangeUpdatesActionBar(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())

	dockPos, ok := waterTileBesideTown(gs)
	if !ok {
		t.Skip("generated world has no water-adjacent town")
	}
	gs.player.SetPos(dockPos)

	ViewType = world.ViewTypeMainMap
	gs.refreshActionBar()

	gs.handleKeyPress("Enter")
	if ViewType != world.ViewTypeDock {
		t.Fatalf("Enter did not open dock, ViewType = %d", ViewType)
	}
	barLabels := []string{}
	for _, b := range gs.buttons {
		if b.rect.Overlaps(actionBarRect) {
			barLabels = append(barLabels, b.label)
		}
	}
	if len(barLabels) != 1 || barLabels[0] != "Leave dock (Esc)" {
		t.Fatalf("bar should show dock commands right after the key press, got %v", barLabels)
	}

	gs.handleKeyPress("Escape")
	if ViewType != world.ViewTypeMainMap {
		t.Fatalf("Escape did not leave dock, ViewType = %d", ViewType)
	}
	if labels := gs.actionBarLabels(); len(labels) == 0 || !strings.HasPrefix(labels[0], "Dock") {
		t.Fatalf("bar should show sailing commands right after leaving dock, got %v", labels)
	}
	if gs.currentOverlayScreen().title != "" {
		t.Fatal("overlay should be gone after leaving dock with a key")
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
