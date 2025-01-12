package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"pirate-wars/cmd/common"
)

var ExamineData = common.NewUserActionExamine()

func (m model) miniMapInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		// The "m" key toggles the minimap
		case "m", "enter":
			m.viewType = ViewTypeMainMap
		}
	}
	return m, nil
}

func (m model) sailingInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		// The "m" key displays the minimap
		case "m":
			m.viewType = ViewTypeMiniMap

		case "p":
			if m.viewType == ViewTypeHeatMap {
				m.viewType = ViewTypeMainMap
			} else {
				m.viewType = ViewTypeHeatMap
			}

		// examine something on the map
		case "x":
			m.action = common.UserActionIdExamine
			ExamineData = common.NewUserActionExamine()

		case "left", "h", "a":
			c := m.player.GetPos()
			if c.X > 0 {
				t := common.Coordinates{
					X: c.X - 1,
					Y: c.Y,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}

		case "right", "l", "d":
			c := m.player.GetPos()
			if c.X < m.terrain.World.GetWidth()-1 {
				t := common.Coordinates{
					X: c.X + 1,
					Y: c.Y,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k", "w":
			c := m.player.GetPos()
			if c.Y > 0 {
				t := common.Coordinates{
					X: c.X,
					Y: c.Y - 1,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j", "s":
			c := m.player.GetPos()
			if c.Y < m.terrain.World.GetHeight()-1 {
				t := common.Coordinates{
					X: c.X,
					Y: c.Y + 1,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}

		// The "up+left" and "y" keys move the cursor diagonal up+left
		case "up+left", "y", "q":
			c := m.player.GetPos()
			if c.Y > 0 && c.X > 0 {
				t := common.Coordinates{
					X: c.X - 1,
					Y: c.Y - 1,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}

		// The "down+left" and "b" keys move the cursor diagonal down+left
		case "down+left", "b", "z":
			c := m.player.GetPos()
			if c.Y < m.terrain.World.GetHeight()-1 && c.X > 0 {
				t := common.Coordinates{
					X: c.X - 1,
					Y: c.Y + 1,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}

		// The "upright" and "u" keys move the cursor diagonal up+left
		case "up+right", "u", "e":
			c := m.player.GetPos()
			if c.Y > 0 && c.X < m.terrain.World.GetWidth()-1 {
				t := common.Coordinates{
					X: c.X + 1,
					Y: c.Y - 1,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}

		// The "downright" and "n" keys move the cursor diagonal down+left
		case "down+right", "n", "c":
			c := m.player.GetPos()
			if c.Y < m.terrain.World.GetHeight()-1 && c.X < m.terrain.World.GetWidth()-1 {
				t := common.Coordinates{
					X: c.X + 1,
					Y: c.Y + 1,
				}
				if m.terrain.World.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
		}
	}
	m.logger.Debug(fmt.Sprintf("Player position %v", m.player.GetPos()))
	return m, nil
}

func (m model) actionInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		// exit examine mode
		case "x", "return":
			m.action = common.UserActionIdNone
			ExamineData = common.NewUserActionExamine()

		case "left", "h", "a":
			size := len(ExamineData.List)
			if ExamineData.Idx < 0 {
				// left
				ExamineData.Idx = size - 1
			} else {
				// right
				if ExamineData.Idx == size-1 {
					ExamineData.Idx = 0
				} else {
					ExamineData.Idx = ExamineData.Idx + 1
				}
			}

		case "right", "l", "d":

		}
	}
	return m, nil
}
