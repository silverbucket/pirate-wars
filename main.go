package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
	"os"
	"pirate-wars/cmd/action"
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
			m.logger.Info(fmt.Sprintf("Player initialized at: %v, %v",
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
		return m.miniMapInput(msg)
	} else if m.action == user_action.UserActionIdExamine {
		return m.actionInput(msg)
	} else {
		return m.sailingInput(msg)
	}
}

func (m model) View() string {
	if !m.initialized {
		return "Loading..."
	}

	highlight := ExamineData.GetFocusedEntity()
	npcs := m.npcs.GetVisible(m.player.GetPos())
	visible := []common.AvatarReadOnly{}
	for _, n := range npcs.GetList() {
		visible = append(visible, &n)
	}

	bottomText := ""
	sidePanel := ""

	if m.viewType == world.ViewTypeMiniMap {
		return m.world.Paint(m.player, []common.AvatarReadOnly{}, highlight, world.ViewTypeMiniMap)
	} else {
		if m.action == user_action.UserActionIdNone {
			// user is not doing some meta-action, NPCs can move
			m.npcs.CalcMovements()
		}

		// display main map
		paint := m.world.Paint(m.player, visible, highlight, world.ViewTypeMainMap)

		if m.action == user_action.UserActionIdExamine {
			bottomText += fmt.Sprintf("examining %v", highlight.GetID())
			sidePanel = lipgloss.JoinVertical(lipgloss.Left,
				dialog.ListHeader(fmt.Sprintf("%v", highlight.GetName())),
				dialog.ListItem(fmt.Sprintf("Flag: %v", highlight.GetFlag())),
				dialog.ListItem(fmt.Sprintf("ID: %v", highlight.GetID())),
				dialog.ListItem(fmt.Sprintf("Type: %v", highlight.GetType())),
				dialog.ListItem(fmt.Sprintf("Color: %v", highlight.GetForegroundColor())),
			)
		}
		s := dialog.GetSidebarStyle()
		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			paint,
			s.Background(lipgloss.Color("0")).Render(sidePanel),
		)
		content += "\n" + lipgloss.JoinHorizontal(lipgloss.Top,
			bottomText)
		return m.screen.Render(content)
	}
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
