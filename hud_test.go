package main

import (
	"strings"
	"testing"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/sailing"
)

func testShipStatus() shipStatus {
	cfg := sailing.DefaultConfig()
	wind := sailing.NewFixedWind(cfg, common.FacingN, 2)
	heading := common.FacingN
	point := sailing.ClassifyPointOfSail(heading, wind.Facing)
	return shipStatus{
		Heading:       heading,
		Sail:          sailing.SailFull,
		Wind:          wind,
		Point:         point,
		PointMult:     cfg.PointOfSailMultiplier(point),
		Speed:         0.55,
		TimeOfDay:     "12:00",
		Gold:          50,
		CargoTotal:    0,
		CargoCapacity: 20,
	}
}

func TestShipStatusTextSpeedTwoDecimals(t *testing.T) {
	text := shipStatusText(testShipStatus())
	if !strings.Contains(text, "Speed") || !strings.Contains(text, "0.55") {
		t.Fatalf("speed should display two decimals honestly, got:\n%s", text)
	}
	if strings.Contains(text, "0.6\n") {
		t.Fatal("speed must not round 0.55 up to 0.6")
	}
}

// TestShipStatusReportsTheSpeedInputs is the headline contract of the WASD
// change. Steering became relative, so heading, sail and point of sail stopped
// being inferable from the last key pressed and have to be on screen — they are
// what turn the Speed number from a mystery into a decision.
func TestShipStatusReportsTheSpeedInputs(t *testing.T) {
	text := shipStatusText(testShipStatus())
	for _, field := range []string{"Heading", "Sail", "Wind", "Point", "Speed", "Time", "Cargo", "Gold"} {
		if !strings.Contains(text, field) {
			t.Fatalf("ship status should include %s, got:\n%s", field, text)
		}
	}
	if !strings.Contains(text, "Galleon") {
		t.Fatal("ship status should include Galleon")
	}
	if strings.Contains(text, "Galeon") || strings.Contains(text, "Postion") {
		t.Fatal("ship status should not contain the old misspellings")
	}
	if strings.Contains(text, "{") {
		t.Fatal("ship status should not include raw coordinates")
	}
}

// TestShipStatusDropsThePlaceholderHealth removes a stat that was a constant.
func TestShipStatusDropsThePlaceholderHealth(t *testing.T) {
	if strings.Contains(shipStatusText(testShipStatus()), "Health") {
		t.Fatal("Health was always 100; a stat that never changes is not status")
	}
}

// TestFurledSailIsFlaggedNotJustReported covers the S trap: sail_furled is a 0.0
// multiplier, so two taps of the natural "slow down" key stop the ship dead. The
// panel has to name the cause, not just print Speed 0.00.
func TestFurledSailIsFlaggedNotJustReported(t *testing.T) {
	s := testShipStatus()
	s.Sail = sailing.SailFurled
	s.Speed = 0

	lines := s.statusLines()
	var sail statusLine
	for _, l := range lines {
		if l.label == "Sail" {
			sail = l
		}
	}
	if !sail.warn {
		t.Fatal("furled sail should be flagged as the reason the ship is stopped")
	}
	if reason := s.stallReason(); !strings.Contains(reason, "W") {
		t.Fatalf("stall reason should name the key that fixes it, got %q", reason)
	}
}

// TestInIronsIsExplained covers the other dead stop: pointing into the wind.
func TestInIronsIsExplained(t *testing.T) {
	cfg := sailing.DefaultConfig()
	wind := sailing.NewFixedWind(cfg, common.FacingN, 2)
	s := testShipStatus()
	s.Heading = common.OppositeFacing(wind.Facing)
	s.Point = sailing.ClassifyPointOfSail(s.Heading, wind.Facing)

	if s.Point != sailing.PointIrons {
		t.Fatalf("heading opposite the wind should be in irons, got %s", s.Point.Label())
	}
	if reason := s.stallReason(); !strings.Contains(reason, "irons") {
		t.Fatalf("in irons should be explained, got %q", reason)
	}
}

func TestExaminePanelTextHidesPosition(t *testing.T) {
	examine := entities.NewEmptyViewableEntity()
	text := examinePanelText(examine)
	if strings.Contains(text, "Position:") {
		t.Fatal("examine panel should not include raw position")
	}
}

func TestDebugOverlayDefaultOff(t *testing.T) {
	if debugOverlayVisible {
		t.Fatal("debug overlay should default to off")
	}
}

func TestExamineActionBarLabel(t *testing.T) {
	if got := examineActionBarLabel(entities.NewEmptyViewableEntity()); got != "Examine" {
		t.Fatalf("empty focus label = %q, want Examine", got)
	}
}

func TestDebugOverlayTextIncludesHiddenHUDFields(t *testing.T) {
	wind := sailing.NewWind(sailing.DefaultConfig())
	text := debugOverlayText(entities.NewEmptyViewableEntity().GetPos(), wind)
	for _, field := range []string{"Window:", "Viewport:", "Map:", "Player:", "wind:"} {
		if !strings.Contains(text, field) {
			t.Fatalf("debug overlay should include %s", field)
		}
	}
}
