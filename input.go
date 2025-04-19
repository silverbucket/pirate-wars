package main

import (
	"fmt"
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/dialog"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/user_action"
	"pirate-wars/cmd/world"

	"fyne.io/fyne/v2"
)

var ExamineData = user_action.Examine()
var Action = user_action.UserActionIdNone

const KeyCatAdmin = 0
const KeyCatNav = 1
const KeyCatAction = 3
const KeyCatAux = 4

type keyItem struct {
	key  []string
	cat  int
	help string
	exec func(m GameState)
}

type KeyMap []keyItem

func (m *GameState) handleKeyPress(key *fyne.KeyEvent) {
	if ViewType == world.ViewTypeMainMap {
		m.processInput(key, sailingKeyMap)
	} else if ViewType == world.ViewTypeMiniMap {
		m.processInput(key, miniMapKeyMap)
	}
}

func (m *GameState) processInput(key *fyne.KeyEvent, km KeyMap) {
	for _, e := range km {
		for _, k := range e.key {
			if string(key.Name) == k {
				e.exec(*m)
			}
		}
	}
}

func keyQuit(m GameState) {
	os.Exit(0)
}

var miniMapKeyMap = KeyMap{
	{
		key:  []string{"ctrl+q"},
		cat:  KeyCatAdmin,
		help: "quit",
		exec: keyQuit,
	},
	{
		key:  []string{"M", "Enter"},
		cat:  KeyCatAux,
		help: "exit minimap",
		exec: func(m GameState) {
			fmt.Println("-- minimap exit set")
			ViewType = world.ViewTypeMainMap
		},
	},
}

var sailingKeyMap = KeyMap{
	{
		key:  []string{"?"},
		cat:  KeyCatAdmin,
		help: "help",
		exec: func(m GameState) {
			Action = user_action.UserActionIdHelp
		},
	},
	{
		key:  []string{"M"},
		help: "minimap",
		cat:  KeyCatAux,
		exec: func(m GameState) {
			fmt.Printf("minimap called %v\n", ViewType)
			ViewType = world.ViewTypeMiniMap
			fmt.Printf("minimap set %v\n", ViewType)
		},
	},
	{
		key:  []string{"X"},
		help: "examine",
		cat:  KeyCatAction,
		exec: func(m GameState) {
			Action = user_action.UserActionIdExamine
			npcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
			ExamineData = user_action.Examine()
			npcs.ForEach(func(n npc.Npc) {
				ExamineData.AddItem(&n)
			})
		},
	},
	{
		key:  []string{"Left", "H", "A"},
		help: "left",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
		},
	},
	{
		key:  []string{"Right", "L", "D"},
		help: "right",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
		},
	},
	{
		key:  []string{"Up", "K", "W"},
		help: "up",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
		},
	},
	{
		key:  []string{"Down", "J", "S"},
		help: "down",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
		},
	},
	{
		key:  []string{"Q", "Y"},
		help: "up & left",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
		},
	},
	{
		key:  []string{"B", "Z"},
		help: "down & left",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
		},
	},
	{
		key:  []string{"U", "E"},
		help: "up & right",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
		},
	},
	{
		key:  []string{"N", "C"},
		help: "down & right",
		cat:  KeyCatNav,
		exec: func(m GameState) {
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
//	if ViewType == world.ViewTypeHeatMap {
//		ViewType = world.ViewTypeMainMap
//	} else {
//		ViewType = world.ViewTypeHeatMap
//	}

var examineKeyMap = KeyMap{
	{
		key:  []string{"ctrl+q"},
		help: "quit",
		cat:  KeyCatAdmin,
		exec: keyQuit,
	},
	{
		key:  []string{"X", "Enter"},
		help: "exit examine mode",
		cat:  KeyCatAux,
		exec: func(m GameState) {
			Action = user_action.UserActionIdNone
			ExamineData = user_action.Examine()
		},
	},
	{
		key:  []string{"Left", "H", "A"},
		help: "examine item left",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			ExamineData.FocusLeft()
		},
	},
	{
		key:  []string{"Right", "L", "D"},
		help: "examine item right",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			ExamineData.FocusRight()
		},
	},
}

func helpText(km KeyMap, cat int) string {
	r := ""
	f := true
	for _, k := range km {
		if k.cat != cat {
			continue
		}
		s := ""
		t := true
		for _, i := range k.key {
			if t {
				t = false
			} else {
				s += "/"
			}
			if i == "up" {
				i = "↑"
			} else if i == "down" {
				i = "↓"
			} else if i == "left" {
				i = "←"
			} else if i == "right" {
				i = "→"
			}
			s += fmt.Sprintf("%v", i)
		}

		if f {
			f = false
		} else {
			r += " • "
		}
		r += fmt.Sprintf("%v: %v", s, k.help)
	}
	return dialog.HelpStyle(r)
}
