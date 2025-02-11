package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/dialog"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/screen"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/user_action"
	"pirate-wars/cmd/world"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

type model struct {
	logger      *zap.SugaredLogger
	world       *world.MapView
	player      *npc.Avatar
	npcs        *npc.Npcs
	towns       *town.Towns
	screen      lipgloss.Style
	initialized bool
	viewType    int
	action      int
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.logger.Info(fmt.Sprintf("Window size: %v", msg))
		screen.SetWindowSize(msg.Width, msg.Height)
		m.logger.Info(fmt.Sprintf("Info pane size %v", screen.InfoPaneSize))
		m.screen = dialog.SetScreenStyle(msg.Width, msg.Height)
		if !m.initialized {
			// only run once at startup
			m.world = world.Init(m.logger)
			m.towns = town.Init(m.world, m.logger)
			m.npcs = npc.Init(m.towns, m.world, m.logger)
			m.player = player.Create(m.world)
			m.initialized = true
			m.logger.Info(fmt.Sprintf("Player initialized at: %+v, %v",
				m.player.GetPos(), m.world.GetPositionType(m.player.GetPos())))
		}
		// redrew minimap every time screen resizes
		m.world.GenerateMiniMap()
		return m, nil
	}
	if !m.initialized {
		return m, nil
	}

	if m.viewType == world.ViewTypeMiniMap {
		return m.getInput(msg, miniMapKeyMap)
	} else if m.action == user_action.UserActionIdExamine {
		return m.getInput(msg, examineKeyMap)
	} else {
		return m.getInput(msg, sailingKeyMap)
	}
}

func (m model) View() string {
	if !m.initialized {
		return "Loading..."
	}

	highlight := ExamineData.GetFocusedEntity()
	npcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
	visible := []common.AvatarReadOnly{}
	for _, n := range npcs.GetList() {
		visible = append(visible, &n)
	}

	bottomText := ""
	sidePanel := ""

	if m.viewType == world.ViewTypeMiniMap {
		paint := m.world.Paint(m.player, []common.AvatarReadOnly{}, highlight, world.ViewTypeMiniMap)
		paint += helpText(miniMapKeyMap, KeyCatAux)
		return paint
	} else {
		if m.action == user_action.UserActionIdNone {
			// user is not doing some meta-action, NPCs can move
			m.npcs.CalcMovements()
		}

		// display main map
		paint := m.world.Paint(m.player, visible, highlight, world.ViewTypeMainMap)

		if m.action == user_action.UserActionIdExamine {
			bottomText += helpText(examineKeyMap, KeyCatAction)
			sidePanel = lipgloss.JoinVertical(lipgloss.Left,
				dialog.ListHeader(fmt.Sprintf("%v", highlight.GetName())),
				dialog.ListItem(fmt.Sprintf("Flag: %v", highlight.GetFlag())),
				dialog.ListItem(fmt.Sprintf("ID: %v", highlight.GetID())),
				dialog.ListItem(fmt.Sprintf("Type: %v", highlight.GetType())),
				dialog.ListItem(fmt.Sprintf("Color: %v", highlight.GetForegroundColor())),
			)
		} else {
			bottomText += lipgloss.JoinHorizontal(
				lipgloss.Top,
				helpText(sailingKeyMap, KeyCatAction),
				helpText(sailingKeyMap, KeyCatAux),
				helpText(sailingKeyMap, KeyCatAdmin),
			)
		}
		s := dialog.GetSidebarStyle()
		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			paint,
			s.Background(lipgloss.Color("0")).Render(sidePanel),
		)
		content += "\n" + bottomText
		return m.screen.Render(content)
	}
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

func main() {
	logger := createLogger()
	logger.Info("Starting...")

	// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
	if _, err := tea.NewProgram(model{
		logger:   logger,
		viewType: world.ViewTypeMainMap,
		action:   user_action.UserActionIdNone,
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Exiting...")
}
