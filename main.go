package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"io"
	"math/rand"
	"os"
	"pirate-wars/cmd/avatar"
	"pirate-wars/cmd/terrain"
	"time"
)

type model struct {
	world        terrain.World
	avatar       avatar.Type
	miniMap      terrain.World
	printMiniMap bool
	dump         io.Writer
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) View() string {
	if m.printMiniMap {
		return m.miniMap.Paint(m.avatar, true)
	} else {
		return m.world.Paint(m.avatar, false)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.printMiniMap = false

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:

		spew.Fdump(m.dump, msg)
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.avatar.GetX() > 0 {
				target := m.avatar.GetX() - 1
				if m.world[target][m.avatar.GetY()].IsPassableByBoat() {
					m.avatar.SetX(target)
				}
			}

		case "right", "l":
			if m.avatar.GetX() < len(m.world[m.avatar.GetY()])-1 {
				target := m.avatar.GetX() + 1
				if m.world[target][m.avatar.GetY()].IsPassableByBoat() {
					m.avatar.SetX(target)
				}
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.avatar.GetY() > 0 {
				target := m.avatar.GetY() - 1
				if m.world[m.avatar.GetX()][target].IsPassableByBoat() {
					m.avatar.SetY(target)
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.avatar.GetY() < len(m.world)-1 {
				target := m.avatar.GetY() + 1
				if m.world[m.avatar.GetX()][target].IsPassableByBoat() {
					m.avatar.SetY(target)
				}
			}

		// The "up+left" and "y" keys move the cursor diagonal up+left
		case "up+left", "y":
			if m.avatar.GetY() > 0 && m.avatar.GetX() > 0 {
				targetY := m.avatar.GetY() - 1
				targetX := m.avatar.GetX() - 1
				if m.world[targetX][targetY].IsPassableByBoat() {
					m.avatar.SetXY(avatar.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "down+left" and "b" keys move the cursor diagonal down+left
		case "down+left", "b":
			if m.avatar.GetY() < len(m.world)-1 && m.avatar.GetX() > 0 {
				targetY := m.avatar.GetY() + 1
				targetX := m.avatar.GetX() - 1
				if m.world[targetX][targetY].IsPassableByBoat() {
					m.avatar.SetXY(avatar.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "upright" and "u" keys move the cursor diagonal up+left
		case "up+right", "u":
			if m.avatar.GetY() > 0 && m.avatar.GetX() < len(m.world[m.avatar.GetY()])-1 {
				targetY := m.avatar.GetY() - 1
				targetX := m.avatar.GetX() + 1
				if m.world[targetX][targetY].IsPassableByBoat() {
					m.avatar.SetXY(avatar.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "downright" and "n" keys move the cursor diagonal down+left
		case "down+right", "n":
			if m.avatar.GetY() < len(m.world)-1 && m.avatar.GetX() < len(m.world[m.avatar.GetY()])-1 {
				targetY := m.avatar.GetY() + 1
				targetX := m.avatar.GetX() + 1
				if m.world[targetX][targetY].IsPassableByBoat() {
					m.avatar.SetXY(avatar.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "m" key displays the minimap
		case "m":
			m.printMiniMap = true

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			//_, ok := m.selected[m.cursor]
			//if ok {
			//	delete(m.selected, m.cursor)
			//} else {
			//	m.selected[m.cursor] = struct{}{}
			//}
		}
		//spew.Fdump(m.dump, fmt.Sprintf("x:%v y:%v", m.avatar.GetX(), m.avatar.GetY()))
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func findStartingPosition(world terrain.World) avatar.Coordinates {
	rand.Seed(time.Now().UnixNano())
	for {
		coords := avatar.Coordinates{X: rand.Intn(terrain.WorldWidth), Y: rand.Intn(terrain.WorldHeight)}
		if world[coords.X][coords.Y] == terrain.TypeDeepWater {
			return coords
		}
	}
}

func main() {
	var dump *os.File
	if _, ok := os.LookupEnv("DEBUG"); ok {
		var err error
		dump, err = os.OpenFile("pirate-wars.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}
	}

	t := terrain.Init()
	world := t.GenerateWorld()
	// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
	if _, err := tea.NewProgram(model{
		world:   world,
		miniMap: world.RenderMiniMap(),
		avatar:  avatar.Create(findStartingPosition(world), '⏏'),
		dump:    dump,
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
