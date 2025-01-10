package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/terrain"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

type model struct {
	logger       *zap.SugaredLogger
	terrain      terrain.Terrain
	player       *terrain.Avatar
	printMiniMap bool
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.printMiniMap {
		return m.miniMapInput(msg)
	} else {
		return m.sailingInput(msg)
	}
}

func (m model) View() string {
	if m.printMiniMap {
		return m.terrain.MiniMap.Paint(m.player, []terrain.Avatar{})
	} else {
		// calc AI stuff
		m.terrain.CalcNpcMovements()
		return m.terrain.World.Paint(m.player, m.terrain.GetNpcAvatars())
	}
}

func createLogger() *zap.SugaredLogger {
	// truncate file
	configFile, err := os.OpenFile(common.LogFile, os.O_TRUNC|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}
	if err = configFile.Close(); err != nil {
		panic(err)
	}
	// create logger
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{common.LogFile}
	cfg.Level = zap.NewAtomicLevelAt(BASE_LOG_LEVEL)
	cfg.Development = DEV_MODE
	cfg.DisableCaller = false
	cfg.DisableStacktrace = false
	cfg.Encoding = "console"
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig = encoderConfig
	logger := zap.Must(cfg.Build())
	defer logger.Sync()
	return logger.Sugar()
}

func main() {
	logger := createLogger()
	logger.Info("Starting...")

	t := terrain.Init(logger)
	t.GenerateWorld()
	t.GenerateTowns()
	t.InitNpcs()
	t.GenerateMiniMap()

	// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
	if _, err := tea.NewProgram(model{
		logger:  logger,
		terrain: *t,
		player:  player.Create(t),
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Exiting...")
}
