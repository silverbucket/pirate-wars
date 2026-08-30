package main

import (
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/sailing"
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

func (m *GameState) syncMinimap() {
	if ViewType == world.ViewTypeMiniMap {
		var viewable entities.ViewableEntities
		for _, t := range m.towns.GetTowns() {
			viewable = append(viewable, &t)
		}
		m.world.ShowMinimapPopup(m.player.GetPos(), viewable, m.window)
	} else {
		m.world.HideMinimapPopup()
	}
}

func (m *GameState) handleKeyPress(key *fyne.KeyEvent) {
	if string(key.Name) == "F3" {
		debugOverlayVisible = !debugOverlayVisible
		m.updateDebugOverlay()
		return
	}

	if ViewType == world.ViewTypeMainMap {
		m.processInput(key, sailingKeyMap)
	} else if ViewType == world.ViewTypeMiniMap {
		m.processInput(key, miniMapKeyMap)
	} else if ViewType == world.ViewTypeExamine {
		m.processInput(key, examineKeyMap)
	}
	m.syncMinimap()
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
			towns := m.towns.GetVisible(m.player.GetPos())
			ExamineData = user_action.Examine()
			if len(npcs.GetList()) > 0 || len(towns) > 0 {
				ViewType = world.ViewTypeExamine
				npcs.ForEach(func(n npc.Npc) {
					ExamineData.AddItem(&n)
				})
				for i := range towns {
					ExamineData.AddItem(&towns[i])
				}
			}
		},
	},
	{
		key:  []string{"Left", "H", "A"},
		help: "heading W",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: -1, Y: 0})
		},
	},
	{
		key:  []string{"Right", "L", "D"},
		help: "heading E",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: 1, Y: 0})
		},
	},
	{
		key:  []string{"Up", "K", "W"},
		help: "heading N",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: 0, Y: -1})
		},
	},
	{
		key:  []string{"Down", "J", "S"},
		help: "heading S",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: 0, Y: 1})
		},
	},
	{
		key:  []string{"Q", "Y"},
		help: "heading NW",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: -1, Y: -1})
		},
	},
	{
		key:  []string{"B", "Z"},
		help: "heading SW",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: -1, Y: 1})
		},
	},
	{
		key:  []string{"U", "E"},
		help: "heading NE",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: 1, Y: -1})
		},
	},
	{
		key:  []string{"N", "C"},
		help: "heading SE",
		cat:  KeyCatNav,
		exec: func(m GameState) {
			setPlayerHeading(&m, common.Coordinates{X: 1, Y: 1})
		},
	},
	{
		key:  []string{"1"},
		help: "full sail",
		cat:  KeyCatAction,
		exec: func(m GameState) {
			setPlayerSail(&m, sailing.SailFull)
		},
	},
	{
		key:  []string{"2"},
		help: "half sail",
		cat:  KeyCatAction,
		exec: func(m GameState) {
			setPlayerSail(&m, sailing.SailHalf)
		},
	},
	{
		key:  []string{"3"},
		help: "furled sail",
		cat:  KeyCatAction,
		exec: func(m GameState) {
			setPlayerSail(&m, sailing.SailFurled)
		},
	},
	{
		key:  []string{"V"},
		help: "cycle sail",
		cat:  KeyCatAction,
		exec: func(m GameState) {
			cyclePlayerSail(&m)
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
		elements = append(elements, widget.NewLabel(examineActionBarLabel(ExamineData.GetFocusedEntity())))
		keyMap = examineKeyMap
	} else if ViewType == world.ViewTypeMiniMap {
		elements = append(elements, widget.NewLabel("MiniMap"))
		keyMap = miniMapKeyMap
	} else if ViewType == world.ViewTypeMainMap {
		elements = append(elements, widget.NewLabel("Sailing"))
		keyMap = sailingKeyMap
		elements = append(elements, widget.NewButton("Full sail", func() {
			setPlayerSail(gs, sailing.SailFull)
			gs.syncMinimap()
		}))
		elements = append(elements, widget.NewButton("Half sail", func() {
			setPlayerSail(gs, sailing.SailHalf)
			gs.syncMinimap()
		}))
		elements = append(elements, widget.NewButton("Furled", func() {
			setPlayerSail(gs, sailing.SailFurled)
			gs.syncMinimap()
		}))
	}

	for _, k := range keyMap {
		if k.cat != KeyCatAdmin && k.cat != KeyCatNav {
			elements = append(elements, widget.NewButton(k.help, func() {
				k.exec(*gs)
				gs.syncMinimap()
			}))
		}
	}

	return container.NewHBox(elements...)
}
