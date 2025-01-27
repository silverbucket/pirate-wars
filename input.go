package main

import (
	"github.com/charmbracelet/bubbletea"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/user_action"
	"pirate-wars/cmd/world"
)

var ExamineData = user_action.Examine()

const KeyCatAdmin = 0
const KeyCatNav = 1
const KeyCatAction = 3
const KeyCatAux = 4

type keyItem struct {
	key  []string
	cat  int
	help string
	exec func(m model) (tea.Model, tea.Cmd)
}

type KeyMap []keyItem

func (m model) getInput(msg tea.Msg, km KeyMap) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		input := msg.String()
		for _, e := range km {
			for _, k := range e.key {
				if input == k {
					return e.exec(m)
				}
			}
		}
	}
	return m, nil
}

func keyQuit(m model) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

var miniMapKeyMap = KeyMap{
	{
		key:  []string{"ctrl+q"},
		cat:  KeyCatAdmin,
		help: "quit",
		exec: keyQuit,
	},
	{
		key:  []string{"m", "enter"},
		cat:  KeyCatAux,
		help: "exit minimap",
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.viewType = world.ViewTypeMainMap
			return m, nil
		},
	},
}

var sailingKeyMap = KeyMap{
	{
		key:  []string{"?"},
		cat:  KeyCatAdmin,
		help: "help",
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.action = user_action.UserActionIdHelp
			return m, nil
		},
	},
	{
		key:  []string{"m"},
		help: "minimap",
		cat:  KeyCatAux,
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.viewType = world.ViewTypeMiniMap
			return m, nil
		},
	},
	{
		key:  []string{"x"},
		help: "examine",
		cat:  KeyCatAction,
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.action = user_action.UserActionIdExamine
			npcs := m.npcs.GetVisible(m.player.GetPos())
			ExamineData = user_action.Examine()
			npcs.ForEach(func(n npc.Npc) {
				ExamineData.AddItem(&n)
			})
			return m, nil
		},
	},
	{
		key:  []string{"left", "h", "a"},
		help: "left",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.X > 0 {
				t := common.Coordinates{
					X: c.X - 1,
					Y: c.Y,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"right", "l", "d"},
		help: "right",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.X < m.world.GetWidth()-1 {
				t := common.Coordinates{
					X: c.X + 1,
					Y: c.Y,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"up", "k", "w"},
		help: "up",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.Y > 0 {
				t := common.Coordinates{
					X: c.X,
					Y: c.Y - 1,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"down", "j", "s"},
		help: "down",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.Y < m.world.GetHeight()-1 {
				t := common.Coordinates{
					X: c.X,
					Y: c.Y + 1,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"q", "y"},
		help: "up & left",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.Y > 0 && c.X > 0 {
				t := common.Coordinates{
					X: c.X - 1,
					Y: c.Y - 1,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"b", "z"},
		help: "down & left",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.Y < m.world.GetHeight()-1 && c.X > 0 {
				t := common.Coordinates{
					X: c.X - 1,
					Y: c.Y + 1,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"u", "e"},
		help: "up & right",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.Y > 0 && c.X < m.world.GetWidth()-1 {
				t := common.Coordinates{
					X: c.X + 1,
					Y: c.Y - 1,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"n", "c"},
		help: "down & right",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			c := m.player.GetPos()
			if c.Y < m.world.GetHeight()-1 && c.X < m.world.GetWidth()-1 {
				t := common.Coordinates{
					X: c.X + 1,
					Y: c.Y + 1,
				}
				if m.world.IsPassableByBoat(t) {
					m.player.SetPos(t)
				}
			}
			return m, nil
		},
	},
	{
		key:  []string{"ctrl+q"},
		help: "quit",
		cat:  KeyCatAdmin,
		exec: keyQuit,
	},
}

//case "p":
//	if m.viewType == world.ViewTypeHeatMap {
//		m.viewType = world.ViewTypeMainMap
//	} else {
//		m.viewType = world.ViewTypeHeatMap
//	}

var examineKeyMap = KeyMap{
	{
		key:  []string{"ctrl+q"},
		help: "quit",
		cat:  KeyCatAdmin,
		exec: keyQuit,
	},
	{
		key:  []string{"x", "enter"},
		help: "exit examine mode",
		cat:  KeyCatAux,
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.action = user_action.UserActionIdNone
			ExamineData = user_action.Examine()
			return m, nil
		},
	},
	{
		key:  []string{"left", "h", "a"},
		help: "examine item left",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			ExamineData.FocusLeft()
			return m, nil
		},
	},
	{
		key:  []string{"right", "l", "d"},
		help: "examine item right",
		cat:  KeyCatNav,
		exec: func(m model) (tea.Model, tea.Cmd) {
			ExamineData.FocusRight()
			return m, nil
		},
	},
}
