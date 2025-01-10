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
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "m" key toggles the minimap
		case "m", "enter":
			m.printMiniMap = false
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
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "m" key displays the minimap
		case "m":
			m.printMiniMap = true

		case "p":
			m.settings.heatMapEnabled = !m.settings.heatMapEnabled

		case "left", "h":
			if m.player.GetX() > 0 {
				target := m.player.GetX() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: target,
					Y: m.player.GetY(),
				}) {
					m.player.SetX(target)
				}
			}

		case "right", "l":
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
		case "up", "k":
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
		case "down", "j":
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
		case "up+left", "y":
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
		case "down+left", "b":
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
		case "up+right", "u":
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
		case "down+right", "n":
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
