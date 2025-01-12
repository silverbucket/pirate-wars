package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/terrain"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

const ViewTypeMainMap = 0
const ViewTypeHeatMap = 1
const ViewTypeMiniMap = 2

type model struct {
	logger   *zap.SugaredLogger
	terrain  *terrain.Terrain
	player   terrain.Avatar
	viewType int
	action   int
}

var ScreenInitialized = false
var WorldInitialized = false

var borderStyle = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┘",
	BottomRight: "└",
}

var SidebarStyle = lipgloss.NewStyle().
	Align(lipgloss.Left).
	Border(borderStyle).
	Foreground(lipgloss.Color("#FAFAFA")).
	BorderForeground(lipgloss.Color("33")).
	Background(lipgloss.Color("0")).
	Margin(1, 3, 0, 0).
	Padding(1, 2).
	Height(19).
	Width(20)

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.logger.Info(fmt.Sprintf("Window size: %v", msg))
		common.SetWindowSize(msg.Width, msg.Height)
		ScreenInitialized = true
		if !WorldInitialized {
			m.terrain.GenerateWorld()
			m.terrain.GenerateTowns()
			m.terrain.InitNpcs()
			m.player = player.Create(m.terrain)
			WorldInitialized = true
			m.logger.Info(fmt.Sprintf("Player initialized at: %v, %v",
				m.player.GetPos(), m.terrain.World.GetPositionType(m.player.GetPos())))
		}
		m.terrain.GenerateMiniMap()
		return m, nil
	}
	if !ScreenInitialized || !WorldInitialized {
		return m, nil
	}

	if m.viewType == ViewTypeMiniMap {
		return m.miniMapInput(msg)
	} else if m.action == common.UserActionIdExamine {
		return m.actionInput(msg)
	} else {
		return m.sailingInput(msg)
	}
}

func (m model) View() string {
	if !WorldInitialized || !ScreenInitialized {
		return "Loading..."
	}

	highlight := ExamineData.GetFocusedEntity()
	npcs := m.terrain.GetVisibleNpcs(m.player.GetPos())
	visible := []common.AvatarReadOnly{}
	for _, npc := range npcs {
		visible = append(visible, &npc)
	}

	bottomText := ""
	sidePanel := ""

	if m.viewType == ViewTypeMiniMap {
		return m.terrain.MiniMap.Paint(&m.player, []common.AvatarReadOnly{}, highlight)
	} else if m.viewType == ViewTypeHeatMap {
		return m.terrain.Towns[0].HeatMap.Paint(&m.player, visible, highlight)
	} else {
		if m.action == common.UserActionIdNone {
			// user is not doing some meta-action, NPCs can move
			m.terrain.CalcNpcMovements()
		}

		// display main map
		paint := m.terrain.World.Paint(&m.player, visible, highlight)

		if m.action == common.UserActionIdExamine {
			bottomText += fmt.Sprintf("examining %v", highlight.GetID())
			sidePanel += fmt.Sprintf("NPC: %v", highlight.GetID())
		}
		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			paint,
			SidebarStyle.MarginRight(0).Render(sidePanel),
		)
		content += "\n" + lipgloss.JoinHorizontal(lipgloss.Top,
			bottomText)
		return content
	}
}

func main() {
	logger := createLogger()
	logger.Info("Starting...")

	t := terrain.Init(logger)

	// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
	if _, err := tea.NewProgram(model{
		logger:   logger,
		terrain:  t,
		player:   terrain.Avatar{},
		viewType: ViewTypeMainMap,
		action:   common.UserActionIdNone,
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Exiting...")
}
