package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
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
	terrain  terrain.Terrain
	player   terrain.Avatar
	viewType int
	action   int
}

var ScreenInitialized = false
var WorldInitialized = false

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
			m.terrain.GenerateMiniMap()
			WorldInitialized = true
		}
		return m, nil
	}
	if !ScreenInitialized {
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
	if !WorldInitialized {
		return ""
	}
	highlight := ExamineData.GetFocusedEntity()
	npcs := m.terrain.GetVisibleNpcs(m.player.GetPos())
	visible := []common.AvatarReadOnly{}
	for _, npc := range npcs {
		visible = append(visible, &npc)
	}
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
		return m.terrain.World.Paint(&m.player, visible, highlight)
	}
}

func main() {
	logger := createLogger()
	logger.Info("Starting...")

	t := terrain.Init(logger)

	// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
	if _, err := tea.NewProgram(model{
		logger:   logger,
		terrain:  *t,
		player:   player.Create(t),
		viewType: ViewTypeMainMap,
		action:   common.UserActionIdNone,
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Exiting...")
}
