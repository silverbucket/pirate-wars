package main

import (
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/sailing"
	"strings"
	"testing"
)

func TestShipStatusTextUsesGalleonSpelling(t *testing.T) {
	wind := sailing.NewWind(sailing.DefaultConfig())
	text := shipStatusText(2.5, wind)
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
	for _, field := range []string{"Health:", "Speed:", "Wind:", "Cargo:", "Gold:"} {
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
