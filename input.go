package main

import (
	"fmt"
	"image/color"
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/user_action"
	"pirate-wars/cmd/world"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
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
	key []string
	cat int
	// label is the action name shown on the action bar; the bound key is appended
	// automatically so the bar doubles as the key cheat-sheet.
	label string
	// barVisible gates context-dependent commands such as dock.
	barVisible func(gs *GameState) bool
	exec       func(m *GameState)
}

type KeyMap []keyItem

var keyDisplayNames = map[string]string{
	"Escape": "Esc",
	"Left":   "←",
	"Right":  "→",
	"Up":     "↑",
	"Down":   "↓",
	"ctrl+q": "Ctrl+Q",
}

func keyDisplay(key string) string {
	if name, ok := keyDisplayNames[key]; ok {
		return name
	}
	return key
}

// keyList renders every binding for the command, e.g. "←/H/A".
func (k keyItem) keyList() string {
	shown := make([]string, 0, len(k.key))
	for _, key := range k.key {
		shown = append(shown, keyDisplay(key))
	}
	return strings.Join(shown, "/")
}

// barLabel is the action bar button text, e.g. "Full sail (1)". Every binding is
// listed so no way of running the command is hidden from the player.
func (k keyItem) barLabel() string {
	if len(k.key) == 0 {
		return k.label
	}
	return fmt.Sprintf("%s (%s)", k.label, k.keyList())
}

// isBarButton reports whether the command gets its own button. Heading and admin
// keys are listed in the bar legend instead, to keep the bar from growing too wide.
func (k keyItem) isBarButton() bool {
	return k.label != "" && k.cat != KeyCatNav && k.cat != KeyCatAdmin
}

// legend lists the commands that have no button of their own: heading keys and
// admin keys. It stays on one line so the bar fits the action menu area.
func (km KeyMap) legend() string {
	var headings, admin []string
	for _, k := range km {
		if k.label == "" {
			continue
		}
		switch k.cat {
		case KeyCatNav:
			dir := strings.TrimPrefix(k.label, "Heading ")
			headings = append(headings, fmt.Sprintf("%s:%s", dir, k.keyList()))
		case KeyCatAdmin:
			admin = append(admin, k.barLabel())
		}
	}

	sections := []string{}
	if len(headings) > 0 {
		sections = append(sections, "Heading "+strings.Join(headings, "  "))
	}
	if len(admin) > 0 {
		sections = append(sections, strings.Join(admin, "  "))
	}
	return strings.Join(sections, "   ·   ")
}

// barLabelFor looks up a command's bar label by action name, so overlay buttons
// show the same key hints as the action bar without duplicating the bindings.
func barLabelFor(km KeyMap, label string) string {
	for _, k := range km {
		if k.label == label {
			return k.barLabel()
		}
	}
	return label
}

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

	previousView := ViewType

	if ViewType == world.ViewTypeMainMap {
		m.processInput(key, sailingKeyMap())
	} else if ViewType == world.ViewTypeMiniMap {
		m.processInput(key, miniMapKeyMap())
	} else if ViewType == world.ViewTypeExamine {
		m.processInput(key, examineKeyMap())
	} else if ViewType == world.ViewTypeDock {
		m.processInput(key, dockKeyMap())
	} else if ViewType == world.ViewTypeHail {
		m.processInput(key, hailKeyMap())
	}
	m.syncMinimap()

	// The action bar and overlay are otherwise only updated on the next tick, which
	// would leave them showing the previous view's commands after a keyboard change.
	if ViewType != previousView {
		m.syncOverlay()
	}
	m.refreshActionBar()
}

// syncOverlay matches the modal overlay to the current view.
func (m *GameState) syncOverlay() {
	if ViewType == world.ViewTypeDock || ViewType == world.ViewTypeHail {
		m.showOverlay()
	} else {
		m.hideOverlay()
	}
}

func (m *GameState) refreshActionBar() {
	m.updateActionBarIfNeeded()
	if ActionMenu != nil {
		fyne.Do(ActionMenu.Refresh)
	}
}

func (m *GameState) processInput(key *fyne.KeyEvent, km KeyMap) {
	for _, e := range km {
		for _, k := range e.key {
			if string(key.Name) == k {
				e.exec(m)
			}
		}
	}
}

func keyQuit(m *GameState) {
	os.Exit(0)
}

// The key maps are functions rather than package variables: a variable holding
// handlers that reach the action bar, which reads the key maps back, is an
// initialization cycle, and working around that meant routing handlers through a
// global game state instead of the one they are called with.
func miniMapKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"M", "Enter"},
			cat:   KeyCatAux,
			label: "Exit minimap",
			exec: func(m *GameState) {
				ViewType = world.ViewTypeMainMap
			},
		},
		{
			key:   []string{"ctrl+q"},
			cat:   KeyCatAdmin,
			label: "Quit",
			exec:  keyQuit,
		},
	}
}

func sailingKeyMap() KeyMap {
	return KeyMap{
		{
			key:        []string{"Enter", "O"},
			label:      "Dock",
			cat:        KeyCatAction,
			barVisible: func(gs *GameState) bool { return gs.adjacentDockTown() != nil },
			exec: func(m *GameState) {
				m.tryOpenDock()
			},
		},
		{
			key:   []string{"1"},
			label: "Full sail",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				setPlayerSail(m, sailing.SailFull)
			},
		},
		{
			key:   []string{"2"},
			label: "Half sail",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				setPlayerSail(m, sailing.SailHalf)
			},
		},
		{
			key:   []string{"3"},
			label: "Furled",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				setPlayerSail(m, sailing.SailFurled)
			},
		},
		{
			key:   []string{"V"},
			label: "Cycle sail",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				cyclePlayerSail(m)
			},
		},
		{
			key:   []string{"X"},
			label: "Examine",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
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
			key:   []string{"M"},
			label: "Minimap",
			cat:   KeyCatAux,
			exec: func(m *GameState) {
				ViewType = world.ViewTypeMiniMap
			},
		},
		{
			key:   []string{"Up", "K", "W"},
			label: "Heading N",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: 0, Y: -1})
			},
		},
		{
			key:   []string{"Down", "J", "S"},
			label: "Heading S",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: 0, Y: 1})
			},
		},
		{
			key:   []string{"Left", "H", "A"},
			label: "Heading W",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: -1, Y: 0})
			},
		},
		{
			key:   []string{"Right", "L", "D"},
			label: "Heading E",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: 1, Y: 0})
			},
		},
		{
			key:   []string{"Q", "Y"},
			label: "Heading NW",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: -1, Y: -1})
			},
		},
		{
			key:   []string{"E", "U"},
			label: "Heading NE",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: 1, Y: -1})
			},
		},
		{
			key:   []string{"Z", "B"},
			label: "Heading SW",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: -1, Y: 1})
			},
		},
		{
			key:   []string{"C", "N"},
			label: "Heading SE",
			cat:   KeyCatNav,
			exec: func(m *GameState) {
				setPlayerHeading(m, common.Coordinates{X: 1, Y: 1})
			},
		},
		{
			key:   []string{"?"},
			label: "Help",
			cat:   KeyCatAdmin,
			exec: func(m *GameState) {
				Action = user_action.UserActionIdHelp
			},
		},
		{
			key:   []string{"ctrl+q"},
			label: "Quit",
			cat:   KeyCatAdmin,
			exec:  keyQuit,
		},
	}
}

func hailKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"Enter", "Escape", "X"},
			label: "Dismiss",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.hailData = hail.Payload{}
				ViewType = world.ViewTypeMainMap
				m.hideOverlay()
			},
		},
		{
			key:   []string{"ctrl+q"},
			label: "Quit",
			cat:   KeyCatAdmin,
			exec:  keyQuit,
		},
	}
}

func dockKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"Escape"},
			label: "Leave dock",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.dockTown = nil
				m.dockPage = dockPageMenu
				m.tavernRumor = ""
				ViewType = world.ViewTypeMainMap
				m.hideOverlay()
			},
		},
		{
			key:   []string{"ctrl+q"},
			label: "Quit",
			cat:   KeyCatAdmin,
			exec:  keyQuit,
		},
	}
}

func examineKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"X", "Enter"},
			label: "Exit examine",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				Action = user_action.UserActionIdNone
				ViewType = world.ViewTypeMainMap
				ExamineData = user_action.Examine()
			},
		},
		{
			key:   []string{"Left", "H", "A"},
			label: "Examine left",
			cat:   KeyCatAux,
			exec: func(m *GameState) {
				ExamineData.FocusLeft()
			},
		},
		{
			key:   []string{"Right", "L", "D"},
			label: "Examine right",
			cat:   KeyCatAux,
			exec: func(m *GameState) {
				ExamineData.FocusRight()
			},
		},
		{
			key:   []string{"ctrl+q"},
			label: "Quit",
			cat:   KeyCatAdmin,
			exec:  keyQuit,
		},
	}
}

// actionBarContext returns the bar title and the key map driving the current view.
func actionBarContext() (string, KeyMap) {
	switch ViewType {
	case world.ViewTypeExamine:
		return examineActionBarLabel(ExamineData.GetFocusedEntity()), examineKeyMap()
	case world.ViewTypeMiniMap:
		return "MiniMap", miniMapKeyMap()
	case world.ViewTypeDock:
		return "Dock", dockKeyMap()
	case world.ViewTypeHail:
		return "Hail", hailKeyMap()
	case world.ViewTypeMainMap:
		return "Sailing", sailingKeyMap()
	}
	return "", nil
}

// ActionItems builds the bottom bar as a row of buttons, one per available command
// and each labelled with its key, above a legend line listing the view name plus the
// heading and admin keys that have no button of their own. Every binding is visible.
func (gs *GameState) ActionItems() *fyne.Container {
	title, keyMap := actionBarContext()

	buttons := []fyne.CanvasObject{}
	for _, k := range keyMap {
		item := k
		if !item.isBarButton() {
			continue
		}
		if item.barVisible != nil && !item.barVisible(gs) {
			continue
		}
		buttons = append(buttons, widget.NewButton(item.barLabel(), func() {
			item.exec(gs)
			gs.syncMinimap()
		}))
	}

	caption := title
	if legend := keyMap.legend(); legend != "" {
		if caption != "" {
			caption += "   ·   "
		}
		caption += legend
	}

	rows := []fyne.CanvasObject{container.NewHBox(buttons...)}
	if caption != "" {
		rows = append(rows, actionBarCaption(caption))
	}
	return container.NewVBox(rows...)
}

// actionBarCaption renders the legend line. canvas.Text is used instead of a Label
// so the legend and the button row both fit the action menu height.
func actionBarCaption(text string) *canvas.Text {
	caption := canvas.NewText(text, color.RGBA{R: 200, G: 200, B: 200, A: 255})
	caption.TextSize = 13
	return caption
}
