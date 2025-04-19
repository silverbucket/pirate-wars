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

	if ViewType == world.ViewTypeMainMap {
		m.world.Paint(m.player, visible, highlight)
	}
}

// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
func main() {
	logger := createLogger()
	logger.Info("Starting...")

	app := app.New()
	app.Settings().SetTheme(&customDarkTheme{})
	w := app.NewWindow("Pirate Wars")

	logger.Info(fmt.Sprintf("Window Dimensions %+v", window.Window))
	logger.Info(fmt.Sprintf("Viewable Area %+v", window.ViewPort))

	gameState := initGameState(logger)
	mainContent := gameState.world.GetViewPort()
	SidePanel = gameState.createSidePanel()
	ActionMenu = gameState.createActionMenu()

	// Main layout
	viewportBg := canvas.NewRectangle(color.Transparent)
	viewportBg.Resize(fyne.NewSize(float32(window.ViewPort.Dimensions.Width), float32(window.ViewPort.Dimensions.Height)))

	content := container.NewBorder(
		nil,
		ActionMenu,
		nil,
		SidePanel,
		container.NewStack(viewportBg, mainContent),
	)
	w.SetContent(content)
	w.Resize(fyne.NewSize(float32(window.Window.Width), float32(window.Window.Height)))
	w.SetFixedSize(true) // don't allow resizing for now

	go gameState.gameLoop()

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		gameState.handleKeyPress(key)
		if ViewType == world.ViewTypeMiniMap {
			gameState.world.ShowMinimapPopup(gameState.player.GetPos(), w)
		} else {
			gameState.world.HideMinimapPopup()
		}
	})

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
