package main

import (
	"fmt"
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/player"
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
	paused      bool
	initialized bool
	logger      *zap.SugaredLogger
	world       *world.MapView
	player      *entities.Avatar
	npcs        *npc.Npcs
	towns       *town.Towns
}

func initGameState(logger *zap.SugaredLogger) *GameState {
	gs := GameState{
		paused:      true,
		initialized: false,
	}
	gs.logger = logger
	gs.world = world.Init(gs.logger)
	gs.towns = town.Init(gs.world, gs.logger)
	gs.npcs = npc.Init(gs.towns, gs.world, gs.logger)
	gs.player = player.Create(gs.world)
	return &gs
}

func (gs *GameState) sidePanelContent(examine entities.ViewableEntity) *fyne.Container {
	shipStatusContent := widget.NewLabel(
		fmt.Sprintf("Galeon\nPostion %+v\nHealth: %d\nSpeed: %d\nCargo: %d\n", gs.player.GetPos(), 100, 5, 250),
	)
	shipStatusContent.Wrapping = fyne.TextWrapWord
	examineContent := widget.NewLabel(
		fmt.Sprintf("Captain: %s\nType: %s\nFlag: %s\nPosition: %+v\n",
			examine.GetName(), examine.GetType(), examine.GetFlag(), examine.GetPos()),
	)
	examineContent.Wrapping = fyne.TextWrapWord

	windowContent := widget.NewLabel(fmt.Sprintf("Window: %dx%dpx\nViewport: %dx%dpx\nSide Panel: %dx%dpx\nAction Menu: %dx%dpx\n",
		window.Window.Width, window.Window.Height,
		window.ViewPort.Dimensions.Width, window.ViewPort.Dimensions.Height,
		window.SidePanel.Width, window.SidePanel.Height,
		window.ActionMenu.Width, window.ActionMenu.Height),
	)
	windowContent.Wrapping = fyne.TextWrapWord

	mapContent := widget.NewLabel(
		fmt.Sprintf("Map: %dx%d\nViewport: %dx%d\n",
			common.WorldCols, common.WorldRows, window.ViewPort.Region.Cols, window.ViewPort.Region.Rows),
	)
	mapContent.Wrapping = fyne.TextWrapWord

	content := container.NewVBox(
		widget.NewLabel("Ship Status                        "),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		shipStatusContent,
		layout.NewSpacer(),
		widget.NewLabel("Examine"),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		examineContent,
		layout.NewSpacer(),
		widget.NewLabel("Map Info"),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		mapContent,
		layout.NewSpacer(),
		widget.NewLabel("Window"),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		windowContent,
	)
	content.Resize(fyne.NewSize(float32(window.SidePanel.Width), float32(window.SidePanel.Height)))
	return content
}

func (gs *GameState) updatePanels(examine entities.ViewableEntity) {
	SidePanel.Objects[1] = gs.sidePanelContent(examine)
	ActionMenu.Objects[1] = gs.ActionItems()
	fyne.Do(func() {
		ActionMenu.Refresh()
		SidePanel.Refresh()
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
	// action menu
	actionMenu := gs.ActionItems()
	viewportBg := canvas.NewRectangle(color.Black)

	actionMenu.Resize(fyne.NewSize(float32(window.ActionMenu.Width), float32(window.ActionMenu.Height)))
	return container.NewStack(viewportBg, actionMenu)
}

func (m *GameState) processTick() {
	if m.paused {
		return
	}

	if ViewType == world.ViewTypeMainMap {
		m.npcs.CalcMovements()
	}

	// get visible NPCs
	highlight := ExamineData.GetFocusedEntity()
	visible := []entities.AvatarReadOnly{}
	for _, n := range m.npcs.GetList() {
		visible = append(visible, &n)
	}

	m.updatePanels(highlight)

	m.world.Paint(m.player, visible, highlight)
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
			mainContent := gameState.world.GetViewPort()
			SidePanel = gameState.createSidePanel()
			ActionMenu = gameState.createActionMenu()

			// Main layout
			viewportBg := canvas.NewRectangle(color.Transparent)
			viewportBg.Resize(fyne.NewSize(float32(window.ViewPort.Dimensions.Width), float32(window.ViewPort.Dimensions.Height)))

			gameContent = container.NewBorder(
				nil,
				ActionMenu,
				nil,
				SidePanel,
				container.NewStack(viewportBg, mainContent),
			)

			// Signal that initialization is complete
			close(initComplete)

			go gameState.gameLoop()

			w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
				gameState.handleKeyPress(key)
				if ViewType == world.ViewTypeMiniMap {
					towns := gameState.towns.GetTowns()
					var entities entities.ViewableEntities
					for _, t := range towns {
						entities = append(entities, &t)
					}
					gameState.world.ShowMinimapPopup(gameState.player.GetPos(), entities, w)
				} else {
					gameState.world.HideMinimapPopup()
				}
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
	for {
		time.Sleep(500 * time.Millisecond)
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
