package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"pirate-wars/cmd/avatar"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/terrain"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

type model struct {
	logger       *zap.SugaredLogger
	terrain      terrain.Terrain
	avatar       avatar.Type
	printMiniMap bool
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) View() string {
	if m.printMiniMap {
		return m.terrain.MiniMap.Paint(&m.avatar)
	} else {
		return m.terrain.World.Paint(&m.avatar)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.printMiniMap = false

	switch msg := msg.(type) {
	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.avatar.GetX() > 0 {
				target := m.avatar.GetX() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: target,
					Y: m.avatar.GetY(),
				}) {
					m.avatar.SetX(target)
				}
			}

		case "right", "l":
			if m.avatar.GetX() < m.terrain.World.GetWidth()-1 {
				target := m.avatar.GetX() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: target,
					Y: m.avatar.GetY(),
				}) {
					m.avatar.SetX(target)
				}
			}

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.avatar.GetY() > 0 {
				target := m.avatar.GetY() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: m.avatar.GetX(),
					Y: target,
				}) {
					m.avatar.SetY(target)
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.avatar.GetY() < m.terrain.World.GetHeight()-1 {
				target := m.avatar.GetY() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: m.avatar.GetX(),
					Y: target,
				}) {
					m.avatar.SetY(target)
				}
			}

		// The "up+left" and "y" keys move the cursor diagonal up+left
		case "up+left", "y":
			if m.avatar.GetY() > 0 && m.avatar.GetX() > 0 {
				targetY := m.avatar.GetY() - 1
				targetX := m.avatar.GetX() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.avatar.SetXY(common.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "down+left" and "b" keys move the cursor diagonal down+left
		case "down+left", "b":
			if m.avatar.GetY() < m.terrain.World.GetHeight()-1 && m.avatar.GetX() > 0 {
				targetY := m.avatar.GetY() + 1
				targetX := m.avatar.GetX() - 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.avatar.SetXY(common.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "upright" and "u" keys move the cursor diagonal up+left
		case "up+right", "u":
			if m.avatar.GetY() > 0 && m.avatar.GetX() < m.terrain.World.GetWidth()-1 {
				targetY := m.avatar.GetY() - 1
				targetX := m.avatar.GetX() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.avatar.SetXY(common.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "downright" and "n" keys move the cursor diagonal down+left
		case "down+right", "n":
			if m.avatar.GetY() < m.terrain.World.GetHeight()-1 && m.avatar.GetX() < m.terrain.World.GetWidth()-1 {
				targetY := m.avatar.GetY() + 1
				targetX := m.avatar.GetX() + 1
				if m.terrain.World.IsPassableByBoat(common.Coordinates{
					X: targetX,
					Y: targetY,
				}) {
					m.avatar.SetXY(common.Coordinates{X: targetX, Y: targetY})
				}
			}

		// The "m" key displays the minimap
		case "m":
			m.printMiniMap = true

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			//_, ok := m.selected[m.cursor]
			//if ok {
			//	delete(m.selected, m.cursor)
			//} else {
			//	m.selected[m.cursor] = struct{}{}
			//}
		}
	}

	m.logger.Debug(fmt.Sprintf("moving to x:%v y:%v", m.avatar.GetX(), m.avatar.GetY()))

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func createLogger() *zap.SugaredLogger {
	// truncate file
	configFile, err := os.OpenFile(common.LogFile, os.O_TRUNC, 0664)
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

	// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
	if _, err := tea.NewProgram(model{
		logger:  logger,
		terrain: *t,
		avatar:  avatar.Create(t.RandomPositionDeepWater(), '⏏'),
	}, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Exiting...")
}
