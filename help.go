package main

import (
	"fmt"

	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/world"
)

// helpSection is one titled block of the help screen.
type helpSection struct {
	title string
	rows  [][2]string
}

// helpScreenSections is the game's manual.
//
// The action bar names every binding, but it cannot explain the two things a new
// captain has no way to discover: that A and D turn relative to the current
// heading rather than snapping to a compass point, and that furled sails are a
// dead stop rather than a slow crawl.
func helpScreenSections(cfg sailing.Config) []helpSection {
	return []helpSection{
		{
			title: "Sailing",
			rows: [][2]string{
				{"A / D", "Tack one octant to port / starboard"},
				{"", "Turns are relative to your heading, not the compass."},
				{"", "One turn per tick: a ship cannot spin on the spot."},
				{"W / S", "Set more / less sail"},
				{"1 / 2 / 3", "Jump to full / half / furled sail"},
				{"", "Furled sail means stopped, not slow."},
			},
		},
		{
			title: "Speed",
			rows: [][2]string{
				{"", "Speed = hull x sail x point of sail x wind."},
				{"", fmt.Sprintf("Running with the wind x%.2f, in irons x%.2f.",
					cfg.PointOfSailRun, cfg.PointOfSailIrons)},
				{"", "The compass shows your heading against the wind:"},
				{"", "needles together is fast, opposed is a standstill."},
			},
		},
		{
			title: "Ports and ships",
			rows: [][2]string{
				{"Enter / O", "Dock at a town on the next cell over"},
				{"X", "Examine ships and towns in sight"},
				{"M", "Open the world minimap"},
				{"H", "Hail a ship you have come alongside"},
			},
		},
		{
			title: "General",
			rows: [][2]string{
				{"?", "Show this screen"},
				{"Esc", "Back out of any screen"},
				{"F3", "Debug overlay"},
				{"Ctrl+Q", "Quit"},
			},
		},
	}
}

// helpOverlayScreen renders the manual as an overlay panel.
func (gs *GameState) helpOverlayScreen() overlayScreen {
	rows := []overlayRow{}
	for i, section := range helpScreenSections(gs.sailingCfg) {
		if i > 0 {
			rows = append(rows, overlayRow{text: ""})
		}
		rows = append(rows, overlayRow{text: section.title, heading: true})
		for _, r := range section.rows {
			rows = append(rows, overlayRow{
				text: fmt.Sprintf("  %-10s %s", r[0], r[1]),
				dim:  r[0] == "",
			})
		}
	}
	rows = append(rows, overlayRow{buttons: []overlayAction{
		{label: barLabelFor(helpKeyMap(), "Close help"), do: gs.closeHelp},
	}})
	return overlayScreen{title: "Ship's Articles", rows: rows}
}

// closeHelp returns to sailing from the help screen.
func (gs *GameState) closeHelp() {
	ViewType = world.ViewTypeMainMap
}

// quitConfirmScreen asks before ending the voyage. Nothing is saved between
// runs, so an unguarded Ctrl+Q next to Q on the keyboard throws the session away.
func (gs *GameState) quitConfirmScreen() overlayScreen {
	return overlayScreen{
		title: "Abandon the voyage?",
		rows: []overlayRow{
			{text: "Your cargo, gold and charts are not saved.", dim: true},
			{text: ""},
			{buttons: []overlayAction{
				{label: "Quit (Y)", do: gs.confirmQuit},
				{label: barLabelFor(quitConfirmKeyMap(), "Keep sailing"), do: gs.cancelQuit},
			}},
		},
	}
}

// cancelQuit dismisses the quit confirmation.
func (gs *GameState) cancelQuit() {
	ViewType = world.ViewTypeMainMap
}
