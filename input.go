package main

import (
	"fmt"
	"strings"

	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/user_action"
	"pirate-wars/cmd/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var ExamineData = user_action.Examine()

const KeyCatAdmin = 0

// KeyCatSailPreset marks the 1/2/3 shortcuts that jump straight to a sail state.
// They live on the bar legend rather than getting a button each, so the bar stays
// narrow now that W/S are the primary sail controls.
const KeyCatSailPreset = 1
const KeyCatAction = 3
const KeyCatAux = 4

type keyItem struct {
	key []string
	cat int
	// label is the action name shown on the action bar; the bound key is appended
	// automatically so the bar doubles as the key cheat-sheet.
	label string
	// barEnabled gates context-dependent commands such as dock. A disabled
	// command keeps its place on the bar and greys out, rather than vanishing:
	// removing a button reflows every button to its right, so the bar shifted
	// under the pointer whenever the player drifted past a town.
	barEnabled func(gs *GameState) bool
	// disabledReason is shown as a notice when the player runs the command
	// anyway, so the key never fails silently.
	disabledReason string
	exec           func(m *GameState)
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

// keyDisplay renders a key name for the player, e.g. "Escape" as "Esc".
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

// isBarButton reports whether the command gets its own button. Sail presets and
// admin keys are listed in the bar legend instead, to keep the bar from growing
// too wide.
func (k keyItem) isBarButton() bool {
	return k.label != "" && k.cat != KeyCatSailPreset && k.cat != KeyCatAdmin
}

// legend lists the commands that have no button of their own: the sail presets and
// the admin keys. It stays on one line so the bar fits the action menu area.
func (km KeyMap) legend() string {
	var presets, admin []string
	for _, k := range km {
		if k.label == "" {
			continue
		}
		switch k.cat {
		case KeyCatSailPreset:
			presets = append(presets, k.barLabel())
		case KeyCatAdmin:
			admin = append(admin, k.barLabel())
		}
	}

	sections := []string{}
	if len(presets) > 0 {
		sections = append(sections, "Sail "+strings.Join(presets, "  "))
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

// handleInput dispatches this frame's key presses and pointer taps.
func (m *GameState) handleInput() {
	for _, name := range pressedKeyNames() {
		m.handleKeyPress(name)
	}
	m.handlePointer()
}

// handleKeyPress routes one key through the current view's key map.
func (m *GameState) handleKeyPress(name string) {
	if name == "F3" {
		debugOverlayVisible = !debugOverlayVisible
		return
	}

	switch ViewType {
	case world.ViewTypeMainMap:
		m.processInput(name, sailingKeyMap())
	case world.ViewTypeMiniMap:
		m.processInput(name, miniMapKeyMap())
	case world.ViewTypeExamine:
		m.processInput(name, examineKeyMap())
	case world.ViewTypeDock:
		m.processInput(name, dockKeyMap())
	case world.ViewTypeHail:
		m.processInput(name, hailKeyMap())
	case world.ViewTypeHelp:
		m.processInput(name, helpKeyMap())
	case world.ViewTypeQuitConfirm:
		m.processInput(name, quitConfirmKeyMap())
	}

	// The bar and overlay would otherwise show the previous view's commands for
	// the rest of this frame after a keyboard view change.
	m.refreshActionBar()
}

// refreshActionBar rebuilds the tap targets for the current view. Ebiten hit-tests
// these rectangles every frame, so rebuilding cannot orphan a live button the way
// swapping Fyne widgets under the pointer did.
func (m *GameState) refreshActionBar() {
	m.buttons = m.buildButtons()
}

// processInput runs the command bound to name, or reports why it cannot run.
func (m *GameState) processInput(name string, km KeyMap) {
	for _, e := range km {
		for _, k := range e.key {
			if name != k {
				continue
			}
			if e.barEnabled != nil && !e.barEnabled(m) {
				if e.disabledReason != "" {
					m.setNotice(e.disabledReason)
				}
				return
			}
			e.exec(m)
			return
		}
	}
}

// handlePointer resolves mouse clicks and touch taps against the current buttons.
func (m *GameState) handlePointer() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		m.tap(x, y)
		return
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		m.tap(x, y)
	}
}

// tap runs the command under a pointer press, or explains a disabled one.
func (m *GameState) tap(x, y int) {
	for _, b := range m.buttons {
		if x >= b.rect.Min.X && x < b.rect.Max.X && y >= b.rect.Min.Y && y < b.rect.Max.Y {
			if !b.enabled {
				if b.disabledReason != "" {
					m.setNotice(b.disabledReason)
				}
				return
			}
			if b.action != nil {
				b.action()
			}
			m.refreshActionBar()
			return
		}
	}
}

// pressedKeyNames maps this frame's just-pressed Ebiten keys to the key names
// used by the KeyMap tables.
func pressedKeyNames() []string {
	ctrl := ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)

	var names []string
	for _, k := range inpututil.AppendJustPressedKeys(nil) {
		if ctrl {
			if k == ebiten.KeyQ {
				names = append(names, "ctrl+q")
			}
			continue
		}
		if name := keyName(k); name != "" {
			names = append(names, name)
		}
	}
	return names
}

// keyName maps an Ebiten key to the name used in the key map tables.
func keyName(k ebiten.Key) string {
	switch k {
	case ebiten.KeyArrowLeft:
		return "Left"
	case ebiten.KeyArrowRight:
		return "Right"
	case ebiten.KeyArrowUp:
		return "Up"
	case ebiten.KeyArrowDown:
		return "Down"
	case ebiten.KeyEnter, ebiten.KeyNumpadEnter:
		return "Enter"
	case ebiten.KeyEscape:
		return "Escape"
	case ebiten.KeySlash:
		return "?"
	case ebiten.KeyF1:
		return "F1"
	case ebiten.KeyF3:
		return "F3"
	case ebiten.KeyDigit1, ebiten.KeyNumpad1:
		return "1"
	case ebiten.KeyDigit2, ebiten.KeyNumpad2:
		return "2"
	case ebiten.KeyDigit3, ebiten.KeyNumpad3:
		return "3"
	}
	if k >= ebiten.KeyA && k <= ebiten.KeyZ {
		return string(rune('A' + int(k-ebiten.KeyA)))
	}
	return ""
}

// keyQuit opens the confirmation rather than exiting. Nothing persists between
// runs, so quitting is the most destructive key on the board.
func keyQuit(m *GameState) {
	m.quitReturnView = ViewType
	ViewType = world.ViewTypeQuitConfirm
}

// confirmQuit ends the game loop. Returning ebiten.Termination from Update lets
// Ebiten shut the window down and the logger flush, which os.Exit skipped.
func (gs *GameState) confirmQuit() {
	gs.quitting = true
}

// The key maps are functions rather than package variables: a variable holding
// handlers that reach the action bar, which reads the key maps back, is an
// initialization cycle, and working around that meant routing handlers through a
// global game state instead of the one they are called with.

// quitItem is the Ctrl+Q entry shared by every view.
func quitItem() keyItem {
	return keyItem{
		key:   []string{"ctrl+q"},
		cat:   KeyCatAdmin,
		label: "Quit",
		exec:  keyQuit,
	}
}

// miniMapKeyMap drives the world chart view.
func miniMapKeyMap() KeyMap {
	return KeyMap{
		{
			// Escape backs out of every screen, so it is bound here as well as on
			// the modal overlays rather than only on some of them.
			key:   []string{"M", "Enter", "Escape"},
			cat:   KeyCatAux,
			label: "Exit minimap",
			exec: func(m *GameState) {
				ViewType = world.ViewTypeMainMap
			},
		},
		quitItem(),
	}
}

// sailingKeyMap drives the main map: steering, sail, and the world commands.
func sailingKeyMap() KeyMap {
	return KeyMap{
		{
			key:            []string{"Enter", "O"},
			label:          "Dock",
			cat:            KeyCatAction,
			barEnabled:     func(gs *GameState) bool { return gs.adjacentDockTown() != nil },
			disabledReason: "No town within reach. Sail alongside one first.",
			exec: func(m *GameState) {
				m.tryOpenDock()
			},
		},
		{
			key:            []string{"H"},
			label:          "Hail",
			cat:            KeyCatAction,
			barEnabled:     func(gs *GameState) bool { return gs.hailTarget() != nil },
			disabledReason: "No ship alongside to hail.",
			exec: func(m *GameState) {
				m.tryHailAdjacent()
			},
		},
		{
			key:   []string{"W"},
			label: "More sail",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				trimPlayerSailMore(m)
			},
		},
		{
			key:   []string{"S"},
			label: "Less sail",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				trimPlayerSailLess(m)
			},
		},
		{
			key:   []string{"A"},
			label: "Tack port",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				tackPlayerPort(m)
			},
		},
		{
			key:   []string{"D"},
			label: "Tack starboard",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				tackPlayerStarboard(m)
			},
		},
		{
			key:   []string{"X"},
			label: "Examine",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.openExamine()
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
			key:   []string{"1"},
			label: "Full",
			cat:   KeyCatSailPreset,
			exec: func(m *GameState) {
				setPlayerSail(m, sailing.SailFull)
			},
		},
		{
			key:   []string{"2"},
			label: "Half",
			cat:   KeyCatSailPreset,
			exec: func(m *GameState) {
				setPlayerSail(m, sailing.SailHalf)
			},
		},
		{
			key:   []string{"3"},
			label: "Furled",
			cat:   KeyCatSailPreset,
			exec: func(m *GameState) {
				setPlayerSail(m, sailing.SailFurled)
			},
		},
		{
			key:   []string{"?", "F1"},
			label: "Help",
			cat:   KeyCatAdmin,
			exec: func(m *GameState) {
				ViewType = world.ViewTypeHelp
			},
		},
		quitItem(),
	}
}

// helpKeyMap drives the controls screen.
func helpKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"Escape", "Enter", "?", "F1"},
			label: "Close help",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.closeHelp()
			},
		},
		quitItem(),
	}
}

// quitConfirmKeyMap drives the "abandon the voyage?" guard on Ctrl+Q.
func quitConfirmKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"Escape", "N"},
			label: "Keep sailing",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.cancelQuit()
			},
		},
		{
			key:   []string{"Y"},
			label: "Quit",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.confirmQuit()
			},
		},
	}
}

// hailKeyMap drives the hail overlay.
func hailKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"Enter", "Escape", "X"},
			label: "Dismiss",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.closeHail()
			},
		},
		quitItem(),
	}
}

// dockKeyMap drives the dock overlay and its pages.
func dockKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"Escape"},
			label: "Leave dock",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.closeDock()
			},
		},
		quitItem(),
	}
}

// examineKeyMap drives the examine view and its focus cycling.
func examineKeyMap() KeyMap {
	return KeyMap{
		{
			key:   []string{"X", "Enter", "Escape"},
			label: "Exit examine",
			cat:   KeyCatAction,
			exec: func(m *GameState) {
				m.closeExamine()
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
		quitItem(),
	}
}

// openExamine focuses the ships and towns in sight, or says why it cannot.
func (m *GameState) openExamine() {
	npcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
	towns := m.towns.GetVisible(m.player.GetPos())
	if len(npcs.GetList()) == 0 && len(towns) == 0 {
		m.setNotice("Nothing in sight to examine.")
		return
	}

	ExamineData = user_action.Examine()
	ViewType = world.ViewTypeExamine
	npcs.ForEach(func(n npc.Npc) {
		ExamineData.AddItem(&n)
	})
	for i := range towns {
		ExamineData.AddItem(&towns[i])
	}
}

// closeExamine drops the focus and returns to sailing.
func (m *GameState) closeExamine() {
	ViewType = world.ViewTypeMainMap
	ExamineData = user_action.Examine()
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
	case world.ViewTypeHelp:
		return "Help", helpKeyMap()
	case world.ViewTypeQuitConfirm:
		return "Quit", quitConfirmKeyMap()
	case world.ViewTypeMainMap:
		return "Sailing", sailingKeyMap()
	}
	return "", nil
}
