package main

import (
	"github.com/charmbracelet/bubbletea"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/user_action"
	"pirate-wars/cmd/world"
)

var ExamineData = user_action.Examine()

type keyItem struct {
	key  []string
	help string
	exec func(m model) (tea.Model, tea.Cmd)
}

type KeyMap []keyItem

func keyQuit(m model) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

var miniMapKeyMap = KeyMap{
	{
		key:  []string{"ctrl+q"},
		help: "quit",
		exec: keyQuit,
	},
	{
		key:  []string{"m", "enter"},
		help: "exit minimap",
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.viewType = world.ViewTypeMainMap
			return m, nil
		},
	},
}

func (m model) miniMapInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		input := msg.String()
		for _, e := range miniMapKeyMap {
			for _, k := range e.key {
				if input == k {
					return e.exec(m)
				}
			}
		}
	}
	return m, nil
}

var sailingKeyMap = KeyMap{

	{
		key:  []string{"m"},
		help: "minimap",
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.viewType = world.ViewTypeMiniMap
			return m, nil
		},
	},
	{
		key:  []string{"x"},
		help: "examine",
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
		help: "move left",
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
		help: "move right",
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
		help: "move up",
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
		help: "move down",
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
		help: "move up & left",
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
		help: "move down & left",
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
		help: "move up & right",
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
		help: "move down & right",
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
		exec: keyQuit,
	},
}

//case "p":
//	if m.viewType == world.ViewTypeHeatMap {
//		m.viewType = world.ViewTypeMainMap
//	} else {
//		m.viewType = world.ViewTypeHeatMap
//	}

func (m model) sailingInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		input := msg.String()
		for _, e := range sailingKeyMap {
			for _, k := range e.key {
				if input == k {
					return e.exec(m)
				}
			}
		}
	}
	return m, nil
}

var examineKeyMap = KeyMap{
	{
		key:  []string{"ctrl+q"},
		help: "quit",
		exec: keyQuit,
	},
	{
		key:  []string{"x", "enter"},
		help: "exit examine mode",
		exec: func(m model) (tea.Model, tea.Cmd) {
			m.action = user_action.UserActionIdNone
			ExamineData = user_action.Examine()
			return m, nil
		},
	},
	{
		key:  []string{"left", "h", "a"},
		help: "examine item left",
		exec: func(m model) (tea.Model, tea.Cmd) {
			ExamineData.FocusLeft()
			return m, nil
		},
	},
	{
		key:  []string{"right", "l", "d"},
		help: "examine item right",
		exec: func(m model) (tea.Model, tea.Cmd) {
			ExamineData.FocusRight()
			return m, nil
		},
	},
}

func (m model) examineInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		input := msg.String()
		for _, e := range examineKeyMap {
			for _, k := range e.key {
				if input == k {
					return e.exec(m)
				}
			}
		}
	}
	return m, nil
}
