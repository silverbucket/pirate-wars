package main

import (
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/world"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

// sailingRing is the octant order, clockwise = starboard.
var sailingRing = []common.Facing{
	common.FacingN, common.FacingNE, common.FacingE, common.FacingSE,
	common.FacingS, common.FacingSW, common.FacingW, common.FacingNW,
}

func pressSailing(gs *GameState, key string) {
	ViewType = world.ViewTypeMainMap
	gs.handleKeyPress(key)
}

// tickHelm clears the per-tick steering ration, standing in for a sailing tick.
// The helm answers once per tick, so a test that turns twice has to tick between.
func tickHelm(gs *GameState) {
	gs.turnedThisTick = false
}

// pressSailingAcrossTicks presses a key with a tick in front of it.
func pressSailingAcrossTicks(gs *GameState, key string) {
	tickHelm(gs)
	pressSailing(gs, key)
}

// TestTackStarboardWalksRingClockwise checks D turns one octant to starboard all
// the way around, wrapping NW → N.
func TestTackStarboardWalksRingClockwise(t *testing.T) {
	gs := mainMapGameState()
	gs.player.SetHeading(common.FacingN)

	for i := 1; i <= len(sailingRing); i++ {
		pressSailingAcrossTicks(gs, "D")
		want := sailingRing[i%len(sailingRing)]
		if got := gs.player.GetFacing(); got != want {
			t.Fatalf("press %d: facing = %s, want %s",
				i, common.FacingLabel(got), common.FacingLabel(want))
		}
	}
	if gs.player.GetFacing() != common.FacingN {
		t.Fatal("eight starboard tacks should return to the starting heading")
	}
}

// TestTackPortWalksRingAnticlockwise checks A turns one octant to port, wrapping
// N → NW.
func TestTackPortWalksRingAnticlockwise(t *testing.T) {
	gs := mainMapGameState()
	gs.player.SetHeading(common.FacingN)

	for i := 1; i <= len(sailingRing); i++ {
		pressSailingAcrossTicks(gs, "A")
		want := sailingRing[(len(sailingRing)-i)%len(sailingRing)]
		if got := gs.player.GetFacing(); got != want {
			t.Fatalf("press %d: facing = %s, want %s",
				i, common.FacingLabel(got), common.FacingLabel(want))
		}
	}
	if gs.player.GetFacing() != common.FacingN {
		t.Fatal("eight port tacks should return to the starting heading")
	}
}

// TestTackFromNorth locks the two cases named in the control spec.
func TestTackFromNorth(t *testing.T) {
	gs := mainMapGameState()

	gs.player.SetHeading(common.FacingN)
	pressSailingAcrossTicks(gs, "D")
	if got := gs.player.GetFacing(); got != common.FacingNE {
		t.Fatalf("D from N = %s, want NE", common.FacingLabel(got))
	}

	gs.player.SetHeading(common.FacingN)
	pressSailingAcrossTicks(gs, "A")
	if got := gs.player.GetFacing(); got != common.FacingNW {
		t.Fatalf("A from N = %s, want NW", common.FacingLabel(got))
	}
}

// TestSailStepsAndClamps checks W adds canvas up to full and stops there, and S
// takes it off down to furled and stops there.
func TestSailStepsAndClamps(t *testing.T) {
	gs := mainMapGameState()

	gs.player.SetSail(sailing.SailFurled)
	pressSailing(gs, "W")
	if got := gs.player.GetSail(); got != sailing.SailHalf {
		t.Fatalf("W from furled = %s, want Half", got.Label())
	}
	pressSailing(gs, "W")
	if got := gs.player.GetSail(); got != sailing.SailFull {
		t.Fatalf("W from half = %s, want Full", got.Label())
	}
	pressSailing(gs, "W")
	if got := gs.player.GetSail(); got != sailing.SailFull {
		t.Fatalf("W at full should clamp, got %s", got.Label())
	}

	pressSailing(gs, "S")
	if got := gs.player.GetSail(); got != sailing.SailHalf {
		t.Fatalf("S from full = %s, want Half", got.Label())
	}
	pressSailing(gs, "S")
	if got := gs.player.GetSail(); got != sailing.SailFurled {
		t.Fatalf("S from half = %s, want Furled", got.Label())
	}
	pressSailing(gs, "S")
	if got := gs.player.GetSail(); got != sailing.SailFurled {
		t.Fatalf("S at furled should clamp, got %s", got.Label())
	}
}

// TestSailKeysDoNotChangeHeading and TestTackKeysDoNotChangeSail keep the two
// controls independent.
func TestSailKeysDoNotChangeHeading(t *testing.T) {
	gs := mainMapGameState()
	gs.player.SetHeading(common.FacingE)
	pressSailing(gs, "W")
	pressSailing(gs, "S")
	if got := gs.player.GetFacing(); got != common.FacingE {
		t.Fatalf("sail keys changed heading to %s", common.FacingLabel(got))
	}
}

func TestTackKeysDoNotChangeSail(t *testing.T) {
	gs := mainMapGameState()
	gs.player.SetSail(sailing.SailHalf)
	pressSailingAcrossTicks(gs, "A")
	pressSailingAcrossTicks(gs, "D")
	if got := gs.player.GetSail(); got != sailing.SailHalf {
		t.Fatalf("tack keys changed sail to %s", got.Label())
	}
}

// TestHeadingSnapBindingsGone checks the old compass keys no longer steer.
//
// H is on this list as a former heading key but is bound again, to Hail, so it
// is checked for steering only. The rest must be unbound outright.
func TestHeadingSnapBindingsGone(t *testing.T) {
	dropped := []string{"Q", "E", "Z", "C", "Y", "U", "B", "N", "H", "J", "K", "L", "Left", "Right", "Up", "Down"}
	reboundElsewhere := map[string]bool{"H": true}

	bound := map[string]string{}
	for _, item := range sailingKeyMap() {
		for _, k := range item.key {
			bound[k] = item.label
		}
	}
	for _, key := range dropped {
		if reboundElsewhere[key] {
			continue
		}
		if label, ok := bound[key]; ok {
			t.Fatalf("sailing key %q should be unbound, still runs %q", key, label)
		}
	}

	gs := mainMapGameState()
	for _, key := range dropped {
		gs.player.SetHeading(common.FacingN)
		tickHelm(gs)
		pressSailing(gs, key)
		if got := gs.player.GetFacing(); got != common.FacingN {
			t.Fatalf("%q still set a heading: %s", key, common.FacingLabel(got))
		}
	}
}

// TestHelmAnswersOncePerTick is the cost that makes the wind model matter.
//
// Keys are read at the render rate and ticks are 250ms apart, so without a
// ration the player could spin through all eight octants between two ticks at no
// cost — and the choice between a slow close-hauled line and a longer fast reach
// never arises if the best heading is always one free keypress away.
func TestHelmAnswersOncePerTick(t *testing.T) {
	gs := mainMapGameState()
	gs.player.SetHeading(common.FacingN)

	tickHelm(gs)
	pressSailing(gs, "D")
	pressSailing(gs, "D")
	pressSailing(gs, "D")
	if got := gs.player.GetFacing(); got != common.FacingNE {
		t.Fatalf("three tacks inside one tick turned to %s, want one octant to NE",
			common.FacingLabel(got))
	}

	tickHelm(gs)
	pressSailing(gs, "D")
	if got := gs.player.GetFacing(); got != common.FacingE {
		t.Fatalf("the helm should answer again on the next tick, got %s",
			common.FacingLabel(got))
	}
}

// TestSailTrimIsNotRationed keeps the ration on the helm alone: trimming sail is
// not steering, and clamping it would make W/S feel unresponsive.
func TestSailTrimIsNotRationed(t *testing.T) {
	gs := mainMapGameState()
	gs.player.SetSail(sailing.SailFurled)

	tickHelm(gs)
	pressSailing(gs, "W")
	pressSailing(gs, "W")
	if got := gs.player.GetSail(); got != sailing.SailFull {
		t.Fatalf("two W presses in one tick should reach full sail, got %s", got.Label())
	}
}

// TestHelpScreenOpensAndExplainsRelativeSteering covers the key that used to do
// nothing: "?" set a package variable no code ever read, while the action bar
// advertised it. The screen has to carry what the bar cannot — that A and D turn
// relative to the heading, and that furled sail is a dead stop.
func TestHelpScreenOpensAndExplainsRelativeSteering(t *testing.T) {
	gs := mainMapGameState()
	gs.sailingCfg = sailing.DefaultConfig()

	pressSailing(gs, "?")
	if ViewType != world.ViewTypeHelp {
		t.Fatalf("? should open the help screen, ViewType = %d", ViewType)
	}

	var text strings.Builder
	screen := gs.helpOverlayScreen()
	text.WriteString(screen.title)
	for _, r := range screen.rows {
		text.WriteString(" " + r.text)
	}
	body := text.String()

	for _, want := range []string{"relative", "Furled", "A / D", "W / S"} {
		if !strings.Contains(body, want) {
			t.Fatalf("help screen should explain %q, got:\n%s", want, body)
		}
	}

	gs.handleKeyPress("Escape")
	if ViewType != world.ViewTypeMainMap {
		t.Fatalf("Escape should close help, ViewType = %d", ViewType)
	}
}

// TestQuitAsksFirst covers the most destructive key on the board. Nothing is
// saved between runs, and Ctrl+Q used to call os.Exit straight from the handler.
func TestQuitAsksFirst(t *testing.T) {
	gs := mainMapGameState()

	pressSailing(gs, "ctrl+q")
	if ViewType != world.ViewTypeQuitConfirm {
		t.Fatalf("Ctrl+Q should ask before quitting, ViewType = %d", ViewType)
	}
	if gs.quitting {
		t.Fatal("Ctrl+Q alone should not end the run")
	}

	gs.handleKeyPress("N")
	if ViewType != world.ViewTypeMainMap || gs.quitting {
		t.Fatal("declining should return to sailing")
	}

	pressSailing(gs, "ctrl+q")
	gs.handleKeyPress("Y")
	if !gs.quitting {
		t.Fatal("confirming should end the run")
	}
}

// TestCancelQuitReturnsToTheScreenItWasRaisedFrom covers a review finding.
// Ctrl+Q is bound in every view, so sending the player to the main map on cancel
// would throw away an open dock, hail, minimap or examine for a mis-key.
func TestCancelQuitReturnsToTheScreenItWasRaisedFrom(t *testing.T) {
	for _, from := range []int{
		world.ViewTypeMainMap,
		world.ViewTypeMiniMap,
		world.ViewTypeExamine,
		world.ViewTypeDock,
		world.ViewTypeHail,
	} {
		gs := mainMapGameState()
		ViewType = from

		gs.handleKeyPress("ctrl+q")
		if ViewType != world.ViewTypeQuitConfirm {
			t.Fatalf("Ctrl+Q from view %d did not ask, ViewType = %d", from, ViewType)
		}

		gs.handleKeyPress("Escape")
		if ViewType != from {
			t.Fatalf("cancelling quit from view %d landed in view %d", from, ViewType)
		}
		if gs.quitting {
			t.Fatalf("cancelling quit from view %d ended the run", from)
		}
	}
}

// TestExamineWithNothingInSightSaysSo covers a command that used to fail in
// silence, which leaves the player unable to tell a broken key from an unmet
// condition.
func TestExamineWithNothingInSightSaysSo(t *testing.T) {
	gs := initGameState(zap.NewNop().Sugar())
	gs.npcs = npc.TestNpcsWith()
	ViewType = world.ViewTypeMainMap

	gs.player.SetPos(emptyOceanAwayFromTowns(gs))
	gs.handleKeyPress("X")

	if ViewType == world.ViewTypeExamine {
		t.Fatal("examine should not open with nothing in sight")
	}
	if got := gs.activeNotice(); !strings.Contains(got, "Nothing in sight") {
		t.Fatalf("examine with nothing in sight should say so, got %q", got)
	}
}

// emptyOceanAwayFromTowns finds navigable water with no town in view.
func emptyOceanAwayFromTowns(gs *GameState) common.Coordinates {
	for x := 1; x < 200; x++ {
		for y := 1; y < 200; y++ {
			c := common.Coordinates{X: x, Y: y}
			if gs.world.IsPassableByBoat(c) && len(gs.towns.GetVisible(c)) == 0 {
				return c
			}
		}
	}
	return gs.player.GetPos()
}

// TestSailPresetsStillJump keeps the cheap 1/2/3 shortcuts working.
func TestSailPresetsStillJump(t *testing.T) {
	gs := mainMapGameState()
	for _, c := range []struct {
		key  string
		want sailing.SailSetting
	}{
		{"1", sailing.SailFull},
		{"2", sailing.SailHalf},
		{"3", sailing.SailFurled},
	} {
		pressSailing(gs, c.key)
		if got := gs.player.GetSail(); got != c.want {
			t.Fatalf("%s set sail to %s, want %s", c.key, got.Label(), c.want.Label())
		}
	}
}

// TestSailingBarShowsSteeringKeys is the player-facing contract: the bar names the
// keys that steer and trim.
func TestSailingBarShowsSteeringKeys(t *testing.T) {
	ViewType = world.ViewTypeMainMap
	gs := mainMapGameState()

	labels := strings.Join(gs.actionBarLabels(), " | ")
	for _, want := range []string{"Tack port (A)", "Tack starboard (D)", "More sail (W)", "Less sail (S)"} {
		if !strings.Contains(labels, want) {
			t.Fatalf("bar buttons missing %q, got %s", want, labels)
		}
	}
}

// TestSlashOnlyReportsQuestionMarkWithShift covers a review finding.
//
// keyName mapped ebiten.KeySlash to "?" unconditionally, so plain "/" opened the
// help screen while the action bar, the README and the help screen all named the
// binding as "?". A key that does something nothing advertises is exactly the
// class of problem the self-documenting action bar exists to prevent.
func TestSlashOnlyReportsQuestionMarkWithShift(t *testing.T) {
	if got := keyName(ebiten.KeySlash, true); got != "?" {
		t.Fatalf("shift+slash = %q, want ?", got)
	}
	if got := keyName(ebiten.KeySlash, false); got != "/" {
		t.Fatalf("plain slash = %q, want /", got)
	}

	bound := map[string]bool{}
	for _, item := range sailingKeyMap() {
		for _, k := range item.key {
			bound[k] = true
		}
	}
	if !bound["?"] {
		t.Fatal("help should be bound to ?")
	}
	if bound["/"] {
		t.Fatal("plain / should be unbound: the bar advertises ? and F1")
	}

	gs := mainMapGameState()
	ViewType = world.ViewTypeMainMap
	gs.handleKeyPress("/")
	if ViewType == world.ViewTypeHelp {
		t.Fatal("plain / should not open help")
	}
	gs.handleKeyPress("?")
	if ViewType != world.ViewTypeHelp {
		t.Fatalf("? should open help, ViewType = %d", ViewType)
	}
}

// TestEveryBarBindingIsReachable is the general form of the finding above: every
// key named on the action bar must be a name keyName can actually produce.
func TestEveryBarBindingIsReachable(t *testing.T) {
	producible := map[string]bool{"ctrl+q": true}
	for _, shift := range []bool{false, true} {
		for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
			if n := keyName(k, shift); n != "" {
				producible[n] = true
			}
		}
	}

	maps := map[string]KeyMap{
		"sailing": sailingKeyMap(), "minimap": miniMapKeyMap(), "examine": examineKeyMap(),
		"dock": dockKeyMap(), "hail": hailKeyMap(), "help": helpKeyMap(), "quit": quitConfirmKeyMap(),
	}
	for name, km := range maps {
		for _, item := range km {
			for _, k := range item.key {
				if !producible[k] {
					t.Fatalf("%s key map binds %q to %q, but no key press produces that name",
						name, k, item.label)
				}
			}
		}
	}
}
