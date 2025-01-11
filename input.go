package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"pirate-wars/cmd/common"
)

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
			m.action = common.UserActionExamine

		case "left", "h", "a":
			if m.player.GetX() > 0 {
				target := m.player.GetX() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: target,
					Y: m.player.GetY(),
				}) {
					m.player.SetX(target)
				}
			}

		case "right", "l", "d":
			if m.player.GetX() < m.terrain.World.GetWidth()-1 {
				target := m.player.GetX() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: target,
					Y: m.player.GetY(),
				}) {
					m.player.SetX(target)
				}
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k", "w":
			if m.player.GetY() > 0 {
				target := m.player.GetY() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: m.player.GetX(),
					Y: target,
				}) {
					m.player.SetY(target)
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j", "s":
			if m.player.GetY() < m.terrain.World.GetHeight()-1 {
				target := m.player.GetY() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: m.player.GetX(),
					Y: target,
				}) {
					m.player.SetY(target)
				}
			}

		// The "up+left" and "y" keys move the cursor diagonal up+left
		case "up+left", "y", "q":
			if m.player.GetY() > 0 && m.player.GetX() > 0 {
				targetY := m.player.GetY() - 1
				targetX := m.player.GetX() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.player.SetPos(common.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "down+left" and "b" keys move the cursor diagonal down+left
		case "down+left", "b", "z":
			if m.player.GetY() < m.terrain.World.GetHeight()-1 && m.player.GetX() > 0 {
				targetY := m.player.GetY() + 1
				targetX := m.player.GetX() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.player.SetPos(common.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "upright" and "u" keys move the cursor diagonal up+left
		case "up+right", "u", "e":
			if m.player.GetY() > 0 && m.player.GetX() < m.terrain.World.GetWidth()-1 {
				targetY := m.player.GetY() - 1
				targetX := m.player.GetX() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.player.SetPos(common.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "downright" and "n" keys move the cursor diagonal down+left
		case "down+right", "n", "c":
			if m.player.GetY() < m.terrain.World.GetHeight()-1 && m.player.GetX() < m.terrain.World.GetWidth()-1 {
				targetY := m.player.GetY() + 1
				targetX := m.player.GetX() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.player.SetPos(common.Coordinates{X: targetX, Y: targetY})
				}
			}
		}
	}
	m.logger.Debug(fmt.Sprintf("Player position x:%v y:%v", m.player.GetX(), m.player.GetY()))
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
			m.action = common.UserActionNone

		case "left", "h", "a":

		case "right", "l", "d":

		}
	}
	return m, nil
}
