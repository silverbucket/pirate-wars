package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
	"image/color"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/layout"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/world"
	"time"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

var ViewType = world.ViewTypeMainMap

type model struct {
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

func initModel(logger *zap.SugaredLogger) *model {
	m := model{}
	m.logger = logger
	m.world = world.Init(m.logger)
	m.towns = town.Init(m.world, m.logger)
	m.npcs = npc.Init(m.towns, m.world, m.logger)
	m.player = player.Create(m.world)
	return &m
}

func createSidebar() *fyne.Container {
	// Create the sidebar
	sidebar := container.NewVBox(
		widget.NewLabel("Info Panel"),
		widget.NewLabel("This is the right panel."),
		widget.NewLabel("Width: ~1/4 of window"),
	)
	sidebar.Resize(fyne.NewSize(float32(layout.InfoPane.Width), float32(layout.InfoPane.Height)))
	return sidebar
}

func createActionMenu() *fyne.Container {
	// action menu
	actionMenu := container.NewVBox(
		widget.NewLabel("Action Menu"),
	)
	actionMenu.Resize(fyne.NewSize(float32(layout.ActionMenu.Width), float32(layout.ActionMenu.Height)))
	return actionMenu
}

func (m *model) processTick() {
	m.npcs.CalcMovements()
	// convert to readonly type for display
	visibleNpcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
	visible := []entities.AvatarReadOnly{}
	for _, n := range visibleNpcs.GetList() {
		visible = append(visible, &n)
	}
}

func (m *model) updateWorld() {
	// get visible NPCs
	highlight := ExamineData.GetFocusedEntity()
	visible := []entities.AvatarReadOnly{}
	for _, n := range m.npcs.GetList() {
		visible = append(visible, &n)
	}

	m.world.Paint(m.player, visible, highlight)
}

// ⏅ ⏏ ⏚ ⏛ ⏡ ⪮ ⩯ ⩠ ⩟ ⅏
func main() {
	logger := createLogger()
	logger.Info("Starting...")

	w := app.New().NewWindow("Pirate Wars")
	w.Resize(fyne.NewSize(float32(layout.Window.Width), float32(layout.Window.Height)))
	w.SetFixedSize(true) // don't allow resizing for now

	logger.Info(fmt.Sprintf("Window Dimensions %+v", layout.Window))
	logger.Info(fmt.Sprintf("Viewable Area %+v", layout.ViewableArea))

	m := initModel(logger)
	// redrew minimap every time screen resizes
	//m.world.GenerateMiniMap()
	m.updateWorld()

	sidebar := createSidebar()
	actionMenu := createActionMenu()

	// Separators
	vertLine := canvas.NewRectangle(color.Gray{Y: 128})
	vertLine.Resize(fyne.NewSize(2, 718))
	horizLine := canvas.NewRectangle(color.Gray{Y: 128})
	horizLine.Resize(fyne.NewSize(924, 2))

	// Main layout
	content := container.NewBorder(
		nil,
		container.NewVBox(horizLine, actionMenu),
		nil,
		container.NewHBox(sidebar, vertLine),
		m.world.GetWorldGrid(),
	)

	w.SetContent(content)

	// Handle refresh signals from the goroutine in the main thread
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			// This runs on the main thread because ShowAndRun() processes it
			if ViewType == world.ViewTypeMainMap {
				m.processTick()
				time.AfterFunc(0, func() {
					m.updateWorld()
				})
			}
		}
	}()

	w.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		time.AfterFunc(0, func() {

			if ViewType == world.ViewTypeMainMap {
				m.processInput(key, sailingKeyMap)
			} else if ViewType == world.ViewTypeMiniMap {
				m.processInput(key, miniMapKeyMap)
			}

			if ViewType == world.ViewTypeMiniMap {
				m.world.ShowMinimapPopup(m.player.GetPos(), w)
			} else {
				m.world.HideMinimapPopup()
			}
		})
	})

	w.ShowAndRun()
	logger.Info("Exiting...")
}
