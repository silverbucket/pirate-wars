package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	fyneLayout "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/zap"
	"pirate-wars/cmd/dialog"
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

type model struct {
	logger      *zap.SugaredLogger
	world       *world.MapView
	player      *entities.Avatar
	npcs        *npc.Npcs
	towns       *town.Towns
	initialized bool
	viewType    int
	action      int
}

//func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
//switch msg := msg.(type) {
//case tea.WindowSizeMsg:
//	m.logger.Info(fmt.Sprintf("Window size: %v", msg))
//	screen.SetWindowSize(msg.Width, msg.Height)
//	m.logger.Info(fmt.Sprintf("Info pane size %v", screen.InfoPaneSize))
//	m.screen = dialog.SetScreenStyle(msg.Width, msg.Height)
//	if !m.initialized {
//		// only run once at startup
//		m.world = world.Init(m.logger)
//		m.towns = town.Init(m.world, m.logger)
//		m.npcs = npc.Init(m.towns, m.world, m.logger)
//		m.player = player.Create(m.world)
//		m.initialized = true
//		m.logger.Info(fmt.Sprintf("Player initialized at: %+v, %v",
//			m.player.GetPos(), m.world.GetPositionType(m.player.GetPos())))
//	}
//	// redrew minimap every time screen resizes
//	m.world.GenerateMiniMap()
//	return m, nil
//}
//if !m.initialized {
//	return m, nil
//}
//
//if m.viewType == world.ViewTypeMiniMap {
//	return m.getInput(msg, miniMapKeyMap)
//} else if m.action == user_action.UserActionIdExamine {
//	return m.getInput(msg, examineKeyMap)
//} else {
//	return m.getInput(msg, sailingKeyMap)
//}
//}

//func Update(m model) string {
//	if !m.initialized {
//		return "Loading..."
//	}
//
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
//	if m.viewType == world.ViewTypeMiniMap {
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

func initModel(logger *zap.SugaredLogger) model {
	m := model{}
	m.logger = logger
	m.world = world.Init(m.logger)
	m.towns = town.Init(m.world, m.logger)
	m.npcs = npc.Init(m.towns, m.world, m.logger)
	m.player = player.Create(m.world)
	return m
}

func setLayout(grid *fyne.Container) *fyne.Container {
	// Create the info panel
	infoPanel := container.NewVBox(
		widget.NewLabel("Info Panel"),
		widget.NewLabel("This is the right panel."),
		widget.NewLabel("Width: ~1/4 of window"),
	)
	infoPanelContainer := container.New(fyneLayout.NewVBoxLayout(), infoPanel)

	// Combine grid and info panel in a horizontal split (90% grid, 10% info panel)
	topPortion := container.NewHSplit(grid, infoPanelContainer)
	topPortion.SetOffset(0.9) // 90% for grid, 10% for info panel

	// action menu
	actionMenu := container.NewVBox(
		widget.NewLabel("Action Menu"),
	)
	actionMenuContainer := container.New(fyneLayout.NewHBoxLayout(), actionMenu)

	return container.NewVBox(
		topPortion,
		actionMenuContainer,
	)
}

func processTick(m model) {
	m.npcs.CalcMovements()

	// convert to readonly type for display
	visibleNpcs := m.npcs.GetVisible(m.player.GetPos(), m.player.GetViewableRange())
	visible := []entities.AvatarReadOnly{}
	for _, n := range visibleNpcs.GetList() {
		visible = append(visible, &n)
	}
}

func updateView(m model) {
	// get visible NPCs
	highlight := ExamineData.GetFocusedEntity()
	visible := []entities.AvatarReadOnly{}
	for _, n := range m.npcs.GetList() {
		visible = append(visible, &n)
	}

	m.world.Paint(m.player, visible, highlight, world.ViewTypeMainMap)
}

func helpText(km KeyMap, cat int) string {
	r := ""
	f := true
	for _, k := range km {
		if k.cat != cat {
			continue
		}
		s := ""
		t := true
		for _, i := range k.key {
			if t {
				t = false
			} else {
				s += "/"
			}
			if i == "up" {
				i = "↑"
			} else if i == "down" {
				i = "↓"
			} else if i == "left" {
				i = "←"
			} else if i == "right" {
				i = "→"
			}
			s += fmt.Sprintf("%v", i)
		}

		if f {
			f = false
		} else {
			r += " • "
		}
		r += fmt.Sprintf("%v: %v", s, k.help)
	}
	return dialog.HelpStyle(r)
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

	// redrew minimap every time screen resizes
	//m.world.GenerateMiniMap()

	m := initModel(logger)
	updateView(m)

	w.SetContent(setLayout(m.world.GetFyneGrid()))

	// Handle refresh signals from the goroutine in the main thread
	go func() {
		for {
			time.Sleep(250 * time.Millisecond)
			// This runs on the main thread because ShowAndRun() processes it
			processTick(m)
			time.AfterFunc(0, func() {
				updateView(m)
			})
		}
	}()

	w.ShowAndRun()

	//if _, err := tea.NewProgram(model{
	//	logger:   logger,
	//	viewType: world.ViewTypeMainMap,
	//	action:   user_action.UserActionIdNone,
	//}, tea.WithAltScreen()).Run(); err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	os.Exit(1)
	//}
	logger.Info("Exiting...")
}
