package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap"
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/screen"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/world"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

type model struct {
	logger   *zap.SugaredLogger
	world    *world.MapView
	player   *npc.Avatar
	npcs     *npc.Npcs
	towns    *town.Towns
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
	BottomLeft:  "└",
	BottomRight: "┘",
}

var screenStyle lipgloss.Style

func SetScreenStyle(width int, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("0")).
		Width(width).
		Height(height)
}

func getSidebarStyle() lipgloss.Style {
	var SidebarWidth = (screen.InfoPaneSize * 3)
	if SidebarWidth > 25 {
		SidebarWidth += 1
	} else if SidebarWidth > 18 {
		SidebarWidth += 2
	} else {
		SidebarWidth += 3
	}
	return lipgloss.NewStyle().
		Align(lipgloss.Left).
		Border(borderStyle).
		Foreground(lipgloss.Color("#FAFAFA")).
		BorderForeground(lipgloss.Color("33")).
		BorderBackground(lipgloss.Color("0")).
		Background(lipgloss.Color("0")).
		//Margin(1, 1, 0, 0).
		Padding(1).
		Height(20).
		Width(SidebarWidth)
}

var base = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("0"))

var bullet = lipgloss.NewStyle().SetString("·").
	Foreground(lipgloss.Color("#43BF6D")).
	Background(lipgloss.Color("0")).
	PaddingRight(1).
	String()

var listItem = func(s string) string {
	return bullet + lipgloss.NewStyle().
		Foreground(lipgloss.Color("#969B86")).
		Background(lipgloss.Color("0")).
		Render(s)
}

var listHeader = base.
	//BorderStyle(lipgloss.NormalBorder()).
	//BorderBottom(true).
	Background(lipgloss.Color("0")).
	BorderBackground(lipgloss.Color("0")).
	PaddingBottom(1).
	Width(100).
	Render

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
		screenStyle = SetScreenStyle(msg.Width, msg.Height)
		ScreenInitialized = true
		if !WorldInitialized {
			m.world = world.Init(m.logger)
			m.towns = town.Init(m.world, m.logger)
			m.npcs = npc.Init(m.towns, m.world, m.logger)
			m.player = player.Create(m.world)
			WorldInitialized = true
			m.logger.Info(fmt.Sprintf("Player initialized at: %v, %v",
				m.player.GetPos(), m.world.GetPositionType(m.player.GetPos())))
		}
		m.world.GenerateMiniMap()
		return m, nil
	}
	if !ScreenInitialized || !WorldInitialized {
		return m, nil
	}

	if m.viewType == world.ViewTypeMiniMap {
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
	npcs := m.npcs.GetVisible(m.player.GetPos())
	visible := []common.AvatarReadOnly{}
	for _, npc := range npcs.GetList() {
		visible = append(visible, &npc)
	}

	bottomText := ""
	sidePanel := ""

	if m.viewType == world.ViewTypeMiniMap {
		return m.world.Paint(m.player, []common.AvatarReadOnly{}, highlight, world.ViewTypeMiniMap)
	} else {
		if m.action == common.UserActionIdNone {
			// user is not doing some meta-action, NPCs can move
			m.npcs.CalcMovements()
		}

		// display main map
		paint := m.world.Paint(m.player, visible, highlight, world.ViewTypeMainMap)

		if m.action == common.UserActionIdExamine {
			bottomText += fmt.Sprintf("examining %v", highlight.GetID())
			sidePanel = lipgloss.JoinVertical(lipgloss.Left,
				listHeader(fmt.Sprintf("%v", highlight.GetName())),
				listItem(fmt.Sprintf("Flag: %v", highlight.GetFlag())),
				listItem(fmt.Sprintf("ID: %v", highlight.GetID())),
				listItem(fmt.Sprintf("Type: %v", highlight.GetType())),
				listItem(fmt.Sprintf("Color: %v", highlight.GetForegroundColor())),
			)
		}
		s := getSidebarStyle()
		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			paint,
			s.Background(lipgloss.Color("0")).Render(sidePanel),
		)
		content += "\n" + lipgloss.JoinHorizontal(lipgloss.Top,
			bottomText)
		return screenStyle.Render(content)
	}
}

func main() {
	logger := createLogger()
	logger.Info("Starting...")

	//t := terrain.Init(logger)

	// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
	if _, err := tea.NewProgram(model{
		logger:   logger,
		viewType: world.ViewTypeMainMap,
		action:   common.UserActionIdNone,
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Exiting...")
}
