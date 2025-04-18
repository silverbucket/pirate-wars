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
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

var ViewType = world.ViewTypeMainMap

type GameState struct {
	logger *zap.SugaredLogger
	world  *world.MapView
	player *entities.Avatar
	npcs   *npc.Npcs
	towns  *town.Towns
}

//if ViewType == world.ViewTypeMiniMap {
//	return m.getInput(msg, miniMapKeyMap)
//} else if m.action == user_action.UserActionIdExamine {
//	return m.getInput(msg, examineKeyMap)
//} else {
//	return m.getInput(msg, sailingKeyMap)
//}

//	highlight := ExamineData.GetFocusedEntity()
//	npcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
//	visible := []entities.AvatarReadOnly{}
//	for _, n := range npcs.GetList() {
//		visible = append(visible, &n)
//	}
//
//	bottomText := ""
//	sidePanel := ""
//
//	if ViewType == world.ViewTypeMiniMap {
//		//paint := m.world.Paint(m.player, []entities.AvatarReadOnly{}, highlight, world.ViewTypeMiniMap)
//		//paint += helpText(miniMapKeyMap, KeyCatAux)
//		//return paint
//	} else {
//		if m.action == user_action.UserActionIdNone {
//			// user is not doing some meta-action, NPCs can move
//			m.npcs.CalcMovements()
//		}
//
//		// display main map
//		paint := m.world.Paint(m.player, visible, highlight, world.ViewTypeMainMap)
//
//		if m.action == user_action.UserActionIdExamine {
//			//bottomText += helpText(examineKeyMap, KeyCatAction)
//			//sidePanel = lipgloss.JoinVertical(lipgloss.Left,
//			//	dialog.ListHeader(fmt.Sprintf("%v", highlight.GetName())),
//			//	dialog.ListItem(fmt.Sprintf("Flag: %v", highlight.GetFlag())),
//			//	dialog.ListItem(fmt.Sprintf("ID: %v", highlight.GetID())),
//			//	dialog.ListItem(fmt.Sprintf("Type: %v", highlight.GetType())),
//			//	dialog.ListItem(fmt.Sprintf("Color: %v", highlight.GetForegroundColor())),
//			//)
//		} else {
//			bottomText += lipgloss.JoinHorizontal(
//				lipgloss.Top,
//				helpText(sailingKeyMap, KeyCatAction),
//				helpText(sailingKeyMap, KeyCatAux),
//				helpText(sailingKeyMap, KeyCatAdmin),
//			)
//		}
//		s := dialog.GetSidebarStyle()
//		content := lipgloss.JoinHorizontal(
//			lipgloss.Top,
//			paint,
//			s.Background(lipgloss.Color("0")).Render(sidePanel),
//		)
//		content += "\n" + bottomText
//		return m.screen.Render(content)
//	}
//}

func initGameState(logger *zap.SugaredLogger) *GameState {
	gs := GameState{}
	gs.logger = logger
	gs.world = world.Init(gs.logger)
	gs.towns = town.Init(gs.world, gs.logger)
	gs.npcs = npc.Init(gs.towns, gs.world, gs.logger)
	gs.player = player.Create(gs.world)
	return &gs
}

func createSidePanel(pos common.Coordinates) *fyne.Container {
	// Create the sidebar
	createSidePanel := container.NewVBox(
		widget.NewLabel("Ship Status"),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		widget.NewLabel(fmt.Sprintf("Postion (x: %d, y: %d)", pos.X, pos.Y)),
		widget.NewLabel(fmt.Sprintf("Health: %d", 100)),
		widget.NewLabel(fmt.Sprintf("Speed: %d", 5)),
		widget.NewLabel(fmt.Sprintf("Cargo: %d", 250)),
		layout.NewSpacer(),
		widget.NewLabel("Map Info"),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		widget.NewLabel(fmt.Sprintf("Map: %dx%d", common.WorldCols, common.WorldRows)),
		widget.NewLabel(fmt.Sprintf("Viewport: %dx%d", window.ViewPort.Width/window.CellSize, window.ViewPort.Height/window.CellSize)),
		layout.NewSpacer(),
		widget.NewLabel("Window"),
		canvas.NewRectangle(color.RGBA{R: 200, G: 200, B: 200, A: 255}),
		widget.NewLabel(fmt.Sprintf("Window: %dx%dpx", window.Window.Width, window.Window.Height)),
		widget.NewLabel(fmt.Sprintf("Viewport: %dx%dpx", window.ViewPort.Width, window.ViewPort.Height)),
		widget.NewLabel(fmt.Sprintf("Side Panel: %dx%dpx", window.SidePanel.Width, window.SidePanel.Height)),
		widget.NewLabel(fmt.Sprintf("Action Menu: %dx%dpx", window.ActionMenu.Width, window.ActionMenu.Height)),
	)
	//createSidePanel.Resize(fyne.NewSize(float32(layout.InfoPane.Width), float32(layout.InfoPane.Height)))
	//return sidebar

	// Create rectSide with MinSize width 100px
	rectSide := canvas.NewRectangle(color.Transparent)
	rectSide.SetMinSize(fyne.NewSize(100, 0))
	// Create fixedSidePanel
	return container.NewStack(rectSide, createSidePanel)
}

func createActionMenu() *fyne.Container {
	// action menu
	actionMenu := container.NewHBox(
		widget.NewLabel("Action Menu"),
		widget.NewButton("Right", func() {
		}),
		widget.NewButton("Left", func() {
		}),
		widget.NewButton("Up", func() {
		}),
		widget.NewButton("Down", func() {
		}),
		widget.NewButton("Settings", func() {
			// Add settings logic
		}),
		canvas.NewRectangle(color.RGBA{R: 180, G: 180, B: 180, A: 255}),
	)
	//actionMenu.Resize(fyne.NewSize(float32(layout.ActionMenu.Width), float32(layout.ActionMenu.Height)))
	//return actionMenu

	// Create rectAction with MinSize height 50px
	rectAction := canvas.NewRectangle(color.Transparent)
	rectAction.SetMinSize(fyne.NewSize(0, 50))
	// Create fixedActionMenu
	return container.NewStack(rectAction, actionMenu)
}

func (m *GameState) processTick() {
	m.npcs.CalcMovements()
	m.updateWorld()
	// // convert to readonly type for display
	// visibleNpcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
	// visible := []entities.AvatarReadOnly{}
	// for _, n := range visibleNpcs.GetList() {
	// 	visible = append(visible, &n)
	// }
}

func (m *GameState) updateWorld() fyne.CanvasObject {
	// get visible NPCs
	highlight := ExamineData.GetFocusedEntity()
	visible := []entities.AvatarReadOnly{}
	for _, n := range m.npcs.GetList() {
		visible = append(visible, &n)
	}

	return m.world.Paint(m.player, visible, highlight)
}

// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
func main() {
	logger := createLogger()
	logger.Info("Starting...")

	w := app.New().NewWindow("Pirate Wars")

	logger.Info(fmt.Sprintf("Window Dimensions %+v", window.Window))
	logger.Info(fmt.Sprintf("Viewable Area %+v", window.ViewPort))

	gameState := initGameState(logger)

	// redrew minimap every time screen resizes
	mainContent := gameState.updateWorld()

	sidePanel := createSidePanel(gameState.player.GetPos())
	actionMenu := createActionMenu()

	// Separators
	//vertLine := canvas.NewRectangle(color.Gray{Y: 128})
	//vertLine.Resize(fyne.NewSize(2, 718))
	//horizLine := canvas.NewRectangle(color.Gray{Y: 128})
	//horizLine.Resize(fyne.NewSize(924, 2))

	// Main layout
	//content := container.NewBorder(
	//	nil,
	//	container.NewVBox(horizLine, actionMenu),
	//	nil,
	//	container.NewHBox(sidebar, vertLine),
	//	g,
	//)
	//content := container.NewBorder(
	//	nil,
	//	actionMenu,
	//	nil,
	//	sidebar,
	//	g,
	//)

	// Create main content
	//mainContent := container.NewVBox(
	//	widget.NewLabel("Main Content"),
	//	canvas.NewRectangle(color.RGBA{R: 220, G: 220, B: 220, A: 255}),
	//)

	w.SetContent(
		container.NewBorder(
			nil,         // top
			actionMenu,  // bottom
			nil,         // left
			sidePanel,   // right
			mainContent, // center
		),
	)

	w.Resize(fyne.NewSize(float32(window.Window.Width), float32(window.Window.Height)))
	w.SetFixedSize(true) // don't allow resizing for now

	// Handle refresh signals from the goroutine in the main thread
	go gameState.gameLoop()

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		gameState.handleKeyPress(key, w)
	})

	w.ShowAndRun()

	logger.Info("Exiting...")
}

func (m *GameState) gameLoop() {
	for {
		time.Sleep(500 * time.Millisecond)
		// This runs on the main thread because ShowAndRun() processes it
		if ViewType == world.ViewTypeMainMap {
			m.processTick()
			// Use fyne.Do to ensure UI updates happen on the main thread
			fyne.Do(func() {
				m.updateWorld()
			})
		}
	}
}
