package main

import (
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/user_action"
	"pirate-wars/cmd/world"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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
	} else if ViewType == world.ViewTypeExamine {
		m.processInput(key, examineKeyMap)
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
		help: "(Ctrl+Q) quit",
		exec: keyQuit,
	},
	{
		key:  []string{"M", "Enter"},
		cat:  KeyCatAux,
		help: "(M) exit minimap",
		exec: func(m GameState) {
			ViewType = world.ViewTypeMainMap
		},
	},
}

var sailingKeyMap = KeyMap{
	{
		key:  []string{"?"},
		cat:  KeyCatAdmin,
		help: "(?) Help",
		exec: func(m GameState) {
			Action = user_action.UserActionIdHelp
		},
	},
	{
		key:  []string{"M"},
		help: "(M) minimap",
		cat:  KeyCatAux,
		exec: func(m GameState) {
			ViewType = world.ViewTypeMiniMap
		},
	},
	{
		key:  []string{"X"},
		help: "(X) examine",
		cat:  KeyCatAction,
		exec: func(m GameState) {
			Action = user_action.UserActionIdExamine
			npcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
			ExamineData = user_action.Examine()
			if len(npcs.GetList()) > 0 {
				ViewType = world.ViewTypeExamine
				npcs.ForEach(func(n npc.Npc) {
					ExamineData.AddItem(&n)
				})
			}
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
		help: "(Ctrl+Q) quit",
		cat:  KeyCatAdmin,
		exec: keyQuit,
	},
}

var examineKeyMap = KeyMap{
	{
		key:  []string{"X", "Enter"},
		help: "(X) exit examine mode",
		cat:  KeyCatAction,
		exec: func(m GameState) {
			Action = user_action.UserActionIdNone
			ViewType = world.ViewTypeMainMap
			ExamineData = user_action.Examine()
		},
	},
	{
		key:  []string{"Left", "H", "A"},
		help: "(←) examine left",
		cat:  KeyCatAux,
		exec: func(m GameState) {
			ExamineData.FocusLeft()
		},
	},
	{
		key:  []string{"Right", "L", "D"},
		help: "(→) examine right",
		cat:  KeyCatAux,
		exec: func(m GameState) {
			ExamineData.FocusRight()
		},
	},
	{
		key:  []string{"ctrl+q"},
		help: "(Ctrl+Q) quit",
		cat:  KeyCatAdmin,
		exec: keyQuit,
	},
}

func (gs *GameState) ActionItems() *fyne.Container {
	elements := []fyne.CanvasObject{}

	var keyMap KeyMap
	if ViewType == world.ViewTypeExamine {
		elements = append(elements, widget.NewLabel("Examine"))
		keyMap = examineKeyMap
	} else if ViewType == world.ViewTypeMiniMap {
		elements = append(elements, widget.NewLabel("MiniMap"))
		keyMap = miniMapKeyMap
	} else if ViewType == world.ViewTypeMainMap {
		elements = append(elements, widget.NewLabel("Sailing"))
		keyMap = sailingKeyMap
	}

	for _, k := range keyMap {
		if k.cat != KeyCatAdmin && k.cat != KeyCatNav {
			elements = append(elements, widget.NewButton(k.help, func() {
				k.exec(*gs)
			}))
		}
	}

	return container.NewHBox(elements...)
}
