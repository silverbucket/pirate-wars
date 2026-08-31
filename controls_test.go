package main

import (
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/world"
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

// TestTackStarboardWalksRingClockwise checks D turns one octant to starboard all
// the way around, wrapping NW → N.
func TestTackStarboardWalksRingClockwise(t *testing.T) {
	gs := mainMapGameState()
	gs.player.SetHeading(common.FacingN)

	for i := 1; i <= len(sailingRing); i++ {
		pressSailing(gs, "D")
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
		pressSailing(gs, "A")
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
	pressSailing(gs, "D")
	if got := gs.player.GetFacing(); got != common.FacingNE {
		t.Fatalf("D from N = %s, want NE", common.FacingLabel(got))
	}

	gs.player.SetHeading(common.FacingN)
	pressSailing(gs, "A")
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
	pressSailing(gs, "A")
	pressSailing(gs, "D")
	if got := gs.player.GetSail(); got != sailing.SailHalf {
		t.Fatalf("tack keys changed sail to %s", got.Label())
	}
}

// TestHeadingSnapBindingsGone checks the old compass keys no longer steer.
func TestHeadingSnapBindingsGone(t *testing.T) {
	dropped := []string{"Q", "E", "Z", "C", "Y", "U", "B", "N", "H", "J", "K", "L", "Left", "Right", "Up", "Down"}

	bound := map[string]string{}
	for _, item := range sailingKeyMap() {
		for _, k := range item.key {
			bound[k] = item.label
		}
	}
	for _, key := range dropped {
		if label, ok := bound[key]; ok {
			t.Fatalf("sailing key %q should be unbound, still runs %q", key, label)
		}
	}

	gs := mainMapGameState()
	for _, key := range dropped {
		gs.player.SetHeading(common.FacingN)
		pressSailing(gs, key)
		if got := gs.player.GetFacing(); got != common.FacingN {
			t.Fatalf("%q still set a heading: %s", key, common.FacingLabel(got))
		}
	}
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
