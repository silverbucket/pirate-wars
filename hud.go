package main

import (
	"fmt"
	"strings"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/window"
)

var debugOverlayVisible = false

// shipStatus is everything the side panel reports about the player's ship.
//
// Steering is relative — A and D turn one octant from the current heading — so
// heading, sail and point of sail are state the player can no longer infer from
// the key they last pressed. Speed is hull × sail × point-of-sail × wind, a
// range of roughly 45× between the best and worst combinations, so all four
// inputs have to be on screen for the number to mean anything.
type shipStatus struct {
	Heading       common.Facing
	Sail          sailing.SailSetting
	Wind          *sailing.Wind
	Point         sailing.PointOfSail
	PointMult     float64
	Speed         float64
	TimeOfDay     string
	Gold          int
	CargoTotal    int
	CargoCapacity int
}

// newShipStatus reads the live sailing state off the game.
func (gs *GameState) newShipStatus() shipStatus {
	s := shipStatus{
		TimeOfDay:     gs.clock.TimeOfDay(),
		Gold:          gs.hold.Gold,
		CargoTotal:    gs.hold.Cargo.Total(),
		CargoCapacity: gs.hold.Cargo.Capacity(),
		Wind:          gs.wind,
	}
	if gs.player != nil {
		s.Heading = gs.player.GetFacing()
		s.Sail = gs.player.GetSail()
		s.Speed = gs.player.GetLastSpeed()
	}
	if gs.wind != nil {
		s.Point = sailing.ClassifyPointOfSail(s.Heading, gs.wind.Facing)
		s.PointMult = gs.sailingCfg.PointOfSailMultiplier(s.Point)
	}
	return s
}

// statusLine is one label/value row of the side panel. warn marks a value that
// explains why the ship is slow or stopped, so the panel can call it out.
type statusLine struct {
	label string
	value string
	warn  bool
}

// statusLines is the side panel readout. The order runs down the speed formula —
// sail, then wind, then the angle between them — so a stalled ship reads as a
// cause rather than just a number.
func (s shipStatus) statusLines() []statusLine {
	windLabel, windStrength := "--", 0
	if s.Wind != nil {
		windLabel, windStrength = s.Wind.Label(), s.Wind.Strength
	}

	sailValue := s.Sail.Label()
	if s.Sail == sailing.SailFurled {
		sailValue = "Furled (no way on)"
	}

	return []statusLine{
		{label: "Heading", value: common.FacingLabel(s.Heading)},
		{label: "Sail", value: sailValue, warn: s.Sail == sailing.SailFurled},
		{label: "Wind", value: fmt.Sprintf("%s (%d)", windLabel, windStrength)},
		{
			label: "Point",
			value: fmt.Sprintf("%s x%.2f", s.Point.ShortLabel(), s.PointMult),
			warn:  s.Point == sailing.PointIrons,
		},
		{label: "Speed", value: fmt.Sprintf("%.2f", s.Speed), warn: s.Speed <= 0},
		{label: "Time", value: s.TimeOfDay},
		{label: "Cargo", value: fmt.Sprintf("%d/%d", s.CargoTotal, s.CargoCapacity)},
		{label: "Gold", value: fmt.Sprintf("%d", s.Gold)},
	}
}

// stallReason explains a stopped ship in the player's own terms, or returns ""
// when the ship is making way. Furling the sails sets the multiplier to zero, so
// two taps of S leave the ship dead in the water with nothing else to say why.
func (s shipStatus) stallReason() string {
	switch {
	case s.Sail == sailing.SailFurled:
		return "Sails furled. Press W to set sail."
	case s.Point == sailing.PointIrons:
		return "In irons. Tack (A/D) off the wind."
	case s.Wind != nil && s.Wind.Strength <= 0:
		return "Becalmed. Wait for the wind."
	}
	return ""
}

// shipStatusText renders the side panel readout as plain text. The drawing code
// uses statusLines directly so it can colour warnings; this stays the readable
// form for tests and the debug overlay.
func shipStatusText(s shipStatus) string {
	var b strings.Builder
	b.WriteString("Galleon\n")
	for _, line := range s.statusLines() {
		fmt.Fprintf(&b, "%-9s%s\n", line.label+":", line.value)
	}
	return b.String()
}

// examinePanelText describes the entity the player is examining.
func examinePanelText(examine entities.ViewableEntity) string {
	return fmt.Sprintf(
		"Captain: %s\nType: %s\nFlag: %s\n",
		examine.GetName(),
		examine.GetType(),
		examine.GetFlag(),
	)
}

// debugOverlayText dumps window, viewport and world geometry behind F3.
func debugOverlayText(playerPos common.Coordinates, wind *sailing.Wind) string {
	windInfo := "wind: n/a"
	if wind != nil {
		windInfo = fmt.Sprintf("wind: %s strength %d", wind.Label(), wind.Strength)
	}
	return fmt.Sprintf(
		"Window: %dx%d\nViewport: %dx%d (%dx%d tiles)\nMap: %dx%d\nPlayer: %+v\n%s\n",
		window.Window.Width,
		window.Window.Height,
		window.ViewPort.Dimensions.Width,
		window.ViewPort.Dimensions.Height,
		window.ViewPort.Region.Cols,
		window.ViewPort.Region.Rows,
		common.WorldCols,
		common.WorldRows,
		playerPos,
		windInfo,
	)
}

// examineActionBarLabel titles the action bar with the examined entity's name.
func examineActionBarLabel(focused entities.ViewableEntity) string {
	if name := focused.GetName(); name != "" {
		return name
	}
	return "Examine"
}
