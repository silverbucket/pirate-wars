package main

import (
	"fmt"
	"image/color"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

var ViewType = world.ViewTypeMainMap
var SidePanel *fyne.Container
var ActionMenu *fyne.Container

type GameState struct {
	paused         bool
	initialized    bool
	logger         *zap.SugaredLogger
	window         fyne.Window
	world          *world.MapView
	player         *entities.Avatar
	npcs           *npc.Npcs
	towns          *town.Towns
	debugOverlay   *widget.Label
	sailingCfg     sailing.Config
	wind           *sailing.Wind
	economyCfg     economy.Config
	clock          *economy.Clock
	hold           player.Hold
	dockTown       *town.Town
	dockPage       dockPage
	tavernRumor    string
	hailData       hail.Payload
	overlayRoot    *fyne.Container
	overlayPanel   *fyne.Container
	actionBarSig   string
	actionBarItems *fyne.Container
}

func initGameState(logger *zap.SugaredLogger) *GameState {
	sailingCfg := sailing.LoadConfig("sailing.cfg")
	economyCfg := economy.LoadConfig("economy.cfg")
	gs := GameState{
		paused:      true,
		initialized: false,
		sailingCfg:  sailingCfg,
		wind:        sailing.NewWind(sailingCfg),
		economyCfg:  economyCfg,
		clock:       economy.NewClock(economyCfg.TicksPerDay),
		hold:        player.NewHold(economyCfg),
	}
	gs.logger = logger
	gs.world = world.Init(gs.logger)
	gs.towns = town.Init(gs.world, gs.logger, economyCfg)
	gs.npcs = npc.Init(gs.towns, gs.world, gs.logger, economyCfg)
	gs.player = player.Create(gs.world)
	return &gs
}

func (gs *GameState) sidePanelContent(examine entities.ViewableEntity) *fyne.Container {
	shipStatusContent := widget.NewLabel(shipStatusText(
		gs.player.GetLastSpeed(),
		gs.wind,
		gs.clock.TimeOfDay(),
		gs.hold.Gold,
		gs.hold.Cargo.Total(),
		gs.hold.Cargo.Capacity(),
	))
	shipStatusContent.Wrapping = fyne.TextWrapWord
	examineContent := widget.NewLabel(examinePanelText(examine))
	examineContent.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		widget.NewLabel("Ship Status                        "),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		shipStatusContent,
		layout.NewSpacer(),
		widget.NewLabel("Examine"),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		examineContent,
	)
	content.Resize(fyne.NewSize(float32(window.SidePanel.Width), float32(window.SidePanel.Height)))
	return content
}

func (gs *GameState) updateDebugOverlay() {
	if gs.debugOverlay == nil {
		return
	}
	if debugOverlayVisible {
		gs.debugOverlay.SetText(debugOverlayText(gs.player.GetPos(), gs.wind))
		gs.debugOverlay.Show()
	} else {
		gs.debugOverlay.Hide()
	}
	gs.debugOverlay.Refresh()
}

func (gs *GameState) actionBarSignature() string {
	switch ViewType {
	case world.ViewTypeMainMap:
		if t := gs.adjacentDockTown(); t != nil {
			return fmt.Sprintf("main:%s", t.GetID())
		}
		return "main:"
	case world.ViewTypeExamine:
		return fmt.Sprintf("examine:%s", ExamineData.GetFocusedEntity().GetID())
	case world.ViewTypeMiniMap:
		return "minimap"
	case world.ViewTypeDock:
		return "dock"
	case world.ViewTypeHail:
		return "hail"
	default:
		return fmt.Sprintf("view:%d", ViewType)
	}
}

// updateActionBarIfNeeded rebuilds the action bar widgets only when the set of
// available commands changes, so buttons survive between ticks and stay clickable.
// The widgets are built even before ActionMenu exists, since createActionMenu
// needs them to assemble the container.
func (gs *GameState) updateActionBarIfNeeded() {
	sig := gs.actionBarSignature()
	if sig != gs.actionBarSig || gs.actionBarItems == nil {
		gs.actionBarSig = sig
		gs.actionBarItems = gs.ActionItems()
		if ActionMenu != nil {
			ActionMenu.Objects[1] = gs.actionBarItems
		}
		return
	}
	if ActionMenu != nil && ActionMenu.Objects[1] != gs.actionBarItems {
		ActionMenu.Objects[1] = gs.actionBarItems
	}
}

func (gs *GameState) updatePanels(examine entities.ViewableEntity) {
	if SidePanel != nil {
		SidePanel.Objects[1] = gs.sidePanelContent(examine)
	}
	gs.updateActionBarIfNeeded()
	fyne.Do(func() {
		if ActionMenu != nil {
			ActionMenu.Refresh()
		}
		if SidePanel != nil {
			SidePanel.Refresh()
		}
	})
}

func (gs *GameState) createSidePanel() *fyne.Container {
	// Create the sidebar content
	content := gs.sidePanelContent(entities.NewEmptyViewableEntity())
	viewportBg := canvas.NewRectangle(color.Black)

	// Create a fixed width container using layout.NewPadded
	sidePanel := container.NewStack(
		viewportBg,
		content,
	)

	// Set minimum size to enforce width
	sidePanel.Resize(fyne.NewSize(float32(window.SidePanel.Width), float32(window.SidePanel.Height)))
	return sidePanel
}

func (gs *GameState) createActionMenu() *fyne.Container {
	gs.actionBarSig = ""
	gs.actionBarItems = nil
	gs.updateActionBarIfNeeded()
	viewportBg := canvas.NewRectangle(color.Black)

	gs.actionBarItems.Resize(fyne.NewSize(float32(window.ActionMenu.Width), float32(window.ActionMenu.Height)))
	return container.NewStack(viewportBg, gs.actionBarItems)
}

func (m *GameState) processTick() {
	if m.paused {
		return
	}

	if ViewType == world.ViewTypeMainMap || ViewType == world.ViewTypeDock || ViewType == world.ViewTypeHail {
		m.clock.Tick()
	}

	if ViewType == world.ViewTypeMainMap {
		m.resolveSailingTick()
	}

	// get visible NPCs
	highlight := ExamineData.GetFocusedEntity()
	visible := []entities.AvatarReadOnly{}
	for _, n := range m.npcs.GetList() {
		visible = append(visible, &n)
	}

	m.updatePanels(highlight)
	m.updateDebugOverlay()

	m.world.AdvanceAnimation()
	m.world.Paint(m.player, visible, highlight, m.wind.Facing)
	m.syncMinimap()
}

// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
func main() {
	app := app.New()
	app.Settings().SetTheme(&customDarkTheme{})

	logger := createLogger()
	logger.Info("Starting...")

	w := app.NewWindow("Pirate Wars")
	w.Resize(fyne.NewSize(float32(window.Window.Width), float32(window.Window.Height)))
	w.SetFixedSize(true) // don't allow resizing for now

	// Create splash overlay
	splash := canvas.NewImageFromFile("./assets/pirate-wars.png")
	splash.Resize(fyne.NewSize(1024, 768))
	splash.FillMode = canvas.ImageFillOriginal

	// Show splash screen immediately
	w.SetContent(splash)
	w.Show()

	// Initialize game state in background
	go func() {
		// Create a channel to signal when initialization is complete
		initComplete := make(chan struct{})
		var gameState *GameState
		var gameContent fyne.CanvasObject

		// Start initialization in a separate goroutine
		go func() {
			logger.Info(fmt.Sprintf("Window Dimensions %+v", window.Window))
			logger.Info(fmt.Sprintf("Viewable Area %+v", window.ViewPort))

			gameState = initGameState(logger)
			gameState.window = w
			gameStateRef = gameState
			mainContent := gameState.world.GetViewPort()
			SidePanel = gameState.createSidePanel()
			ActionMenu = gameState.createActionMenu()

			debugOverlay := widget.NewLabel("")
			debugOverlay.Hide()
			gameState.debugOverlay = debugOverlay

			overlay := gameState.buildOverlayShell()

			// Main layout
			viewportBg := canvas.NewRectangle(color.Transparent)
			viewportBg.Resize(fyne.NewSize(float32(window.ViewPort.Dimensions.Width), float32(window.ViewPort.Dimensions.Height)))

			gameContent = container.NewBorder(
				nil,
				ActionMenu,
				nil,
				SidePanel,
				container.NewStack(viewportBg, mainContent, overlay, debugOverlay),
			)

			// Signal that initialization is complete
			close(initComplete)

			go gameState.gameLoop()

			w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
				gameState.handleKeyPress(key)
			})
		}()

		// Wait for both initialization and minimum splash screen time
		select {
		case <-initComplete:
			// Initialization complete, but still need to wait for minimum splash time
			time.Sleep(2 * time.Second)
		case <-time.After(2 * time.Second):
			// Minimum splash time reached, but initialization might still be in progress
			<-initComplete // Wait for initialization to complete
		}

		// Now switch to game content and unpause
		fyne.Do(func() {
			w.SetContent(gameContent)
			gameState.paused = false
		})
	}()

	w.ShowAndRun()
	logger.Info("Exiting...")
}

func (m *GameState) gameLoop() {
	tick := m.sailingCfg.TickDuration()
	for {
		time.Sleep(tick)
		// Use fyne.Do to ensure UI updates happen on the main thread
		fyne.Do(func() {
			m.processTick()
		})
	}
}

// Custom dark theme implementation
type customDarkTheme struct{}

func (t *customDarkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (t *customDarkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *customDarkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customDarkTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
