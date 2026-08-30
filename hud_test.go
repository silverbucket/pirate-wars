package main

import (
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/sailing"
	"strings"
	"testing"
)

func defaultHUDArgs() (float64, *sailing.Wind, string, int, int, int) {
	wind := sailing.NewWind(sailing.DefaultConfig())
	return 0.55, wind, "12:00", 50, 0, 20
}

func TestShipStatusTextSpeedTwoDecimals(t *testing.T) {
	speed, wind, tod, gold, cargo, cap := defaultHUDArgs()
	text := shipStatusText(speed, wind, tod, gold, cargo, cap)
	if !strings.Contains(text, "Speed: 0.55") {
		t.Fatalf("speed should display two decimals honestly, got:\n%s", text)
	}
	if strings.Contains(text, "Speed: 0.6") {
		t.Fatal("speed must not round 0.55 up to 0.6")
	}
}

func TestShipStatusTextUsesGalleonSpelling(t *testing.T) {
	speed, wind, tod, gold, cargo, cap := defaultHUDArgs()
	speed = 2.5
	text := shipStatusText(speed, wind, tod, gold, cargo, cap)
	if strings.Contains(text, "Galeon") {
		t.Fatal("ship status should use Galleon spelling")
	}
	if !strings.Contains(text, "Galleon") {
		t.Fatal("ship status should include Galleon")
	}
	if strings.Contains(text, "Postion") {
		t.Fatal("ship status should not include misspelled Postion")
	}
	if strings.Contains(text, "{") {
		t.Fatal("ship status should not include raw coordinates")
	}
	for _, field := range []string{"Time:", "Health:", "Speed:", "Wind:", "Cargo:", "Gold:"} {
		if !strings.Contains(text, field) {
			t.Fatalf("ship status should include %s", field)
		}
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
