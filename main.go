package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/gfx"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/harbor"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"

	"github.com/hajimehoshi/ebiten/v2"
	"go.uber.org/zap"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

const splashDuration = 2 * time.Second

var ViewType = world.ViewTypeMainMap

// Screen regions. The viewport keeps the historic 854×728 size; the action bar
// is drawn as an opaque strip over its bottom edge.
var (
	viewportRect  = image.Rect(0, 0, window.ViewPort.Dimensions.Width, window.ViewPort.Dimensions.Height)
	sidePanelRect = image.Rect(window.Window.Width-window.SidePanel.Width, 0, window.Window.Width, window.Window.Height)
	actionBarRect = image.Rect(0, window.Window.Height-window.ActionMenu.Height, window.Window.Width-window.SidePanel.Width, window.Window.Height)
)

type GameState struct {
	paused      bool
	initialized bool
	logger      *zap.SugaredLogger
	world       *world.MapView
	player      *entities.Avatar
	npcs        *npc.Npcs
	towns       *town.Towns
	sailingCfg  sailing.Config
	wind        *sailing.Wind
	economyCfg  economy.Config
	clock       *economy.Clock
	hold        player.Hold
	dockTown    *town.Town
	dockPage    dockPage
	tavernRumor string
	hailData    hail.Payload

	harborRenderer *harbor.Renderer
	harborWorld    *harbor.World
	sailingWorld   *combinedWorld

	buttons      []button
	startedAt    time.Time
	lastTickAt   time.Time
	splashImage  *ebiten.Image
	minimapTex   *ebiten.Image
	minimapDirty bool
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

	// Harbor art is optional. Without the verified PNGs the harbor world stays
	// nil and the game runs on the 32px tilemap.
	if assets, err := harbor.LoadAssets(""); err == nil {
		mask := harbor.NewMask(assets.Mask)
		gs.harborWorld = harbor.NewWorld(mask)
		gs.harborRenderer = harbor.NewRenderer(assets, mask)
		gs.logger.Info("Harbor region loaded (painted backdrop + mask)")
	} else {
		gs.logger.Warnf("Harbor assets unavailable, using tilemap: %v", err)
	}

	gs.sailingWorld = newCombinedWorld(gs.world, gs.harborWorld)
	if gs.harborWorld != nil {
		gs.player = player.Create(gs.world, gs.harborWorld)
	} else {
		gs.player = player.Create(gs.world, nil)
	}
	gs.initialized = true
	return &gs
}

func (gs *GameState) processTick() {
	if gs.paused {
		return
	}

	if ViewType == world.ViewTypeMainMap || ViewType == world.ViewTypeDock || ViewType == world.ViewTypeHail {
		gs.clock.Tick()
	}

	if ViewType == world.ViewTypeMainMap {
		gs.resolveSailingTick()
	}

	gs.world.AdvanceAnimation()
	gs.minimapDirty = true
}

func (gs *GameState) visibleNPCs() []entities.AvatarReadOnly {
	list := gs.npcs.GetList()
	visible := make([]entities.AvatarReadOnly, 0, len(list))
	for i := range list {
		visible = append(visible, &list[i])
	}
	return visible
}

func (gs *GameState) inPaintedHarbor() bool {
	return gs.harborRenderer != nil && harbor.InRegion(gs.player.GetPos())
}

// Update implements ebiten.Game.
func (gs *GameState) Update() error {
	if time.Since(gs.startedAt) < splashDuration {
		return nil
	}
	gs.paused = false

	if gs.buttons == nil {
		gs.refreshActionBar()
	}
	gs.handleInput()

	tick := gs.sailingCfg.TickDuration()
	if time.Since(gs.lastTickAt) >= tick {
		gs.lastTickAt = time.Now()
		gs.processTick()
		gs.refreshActionBar()
	}
	return nil
}

// Draw implements ebiten.Game.
func (gs *GameState) Draw(screen *ebiten.Image) {
	screen.Fill(colorPanel)

	if time.Since(gs.startedAt) < splashDuration {
		gs.drawSplash(screen)
		return
	}

	viewport := screen.SubImage(viewportRect).(*ebiten.Image)
	highlight := ExamineData.GetFocusedEntity()
	if gs.inPaintedHarbor() {
		gs.harborRenderer.Draw(viewport, gs.player, gs.visibleNPCs())
	} else {
		gs.world.Draw(viewport, gs.player, gs.visibleNPCs(), highlight, gs.wind.Facing)
	}

	gs.drawSidePanel(screen, highlight)
	gs.drawActionBar(screen)

	if ViewType == world.ViewTypeMiniMap {
		gs.drawMinimap(screen)
	}
	if ViewType == world.ViewTypeDock || ViewType == world.ViewTypeHail {
		gs.drawOverlay(screen)
	}
	if debugOverlayVisible {
		gs.drawDebugOverlay(screen)
	}
}

// Layout implements ebiten.Game.
func (gs *GameState) Layout(int, int) (int, int) {
	return window.Window.Width, window.Window.Height
}

func (gs *GameState) drawSplash(screen *ebiten.Image) {
	if gs.splashImage != nil {
		op := &ebiten.DrawImageOptions{}
		b := gs.splashImage.Bounds()
		op.GeoM.Translate(
			float64(window.Window.Width-b.Dx())/2,
			float64(window.Window.Height-b.Dy())/2,
		)
		screen.DrawImage(gs.splashImage, op)
		return
	}
	msg := "PIRATE WARS"
	drawText(screen, msg, (window.Window.Width-gfx.TextWidth(msg))/2, window.Window.Height/2, colorHeading)
}

func (gs *GameState) drawSidePanel(screen *ebiten.Image, examine entities.ViewableEntity) {
	fillRect(screen, sidePanelRect, colorPanel)
	strokeRect(screen, sidePanelRect, colorPanelEdge)

	x := sidePanelRect.Min.X + 8
	y := 10
	drawText(screen, "Ship Status", x, y, colorHeading)
	y += gfx.LineHeight + 4
	y = drawTextBlock(screen, shipStatusText(
		gs.player.GetLastSpeed(),
		gs.wind,
		gs.clock.TimeOfDay(),
		gs.hold.Gold,
		gs.hold.Cargo.Total(),
		gs.hold.Cargo.Capacity(),
	), x, y, colorText)

	y += gfx.LineHeight
	drawText(screen, "Examine", x, y, colorHeading)
	y += gfx.LineHeight + 4
	drawTextBlock(screen, examinePanelText(examine), x, y, colorTextDim)
}

// drawActionBar renders the cheat-sheet bar: the legend line naming the view plus
// the heading and admin keys, above one button per available command.
func (gs *GameState) drawActionBar(screen *ebiten.Image) {
	fillRect(screen, actionBarRect, colorPanel)
	strokeRect(screen, actionBarRect, colorPanelEdge)
	drawText(screen, gs.actionBarCaption(), actionBarRect.Min.X+6, actionBarRect.Min.Y+4, colorTextDim)
	for _, b := range gs.buttons {
		if b.rect.Overlaps(actionBarRect) {
			drawButton(screen, b)
		}
	}
}

func (gs *GameState) drawMinimap(screen *ebiten.Image) {
	// The minimap is a whole-world raster, so it is only re-uploaded on ticks.
	if gs.minimapTex == nil || gs.minimapDirty {
		var viewable entities.ViewableEntities
		towns := gs.towns.GetTowns()
		for i := range towns {
			viewable = append(viewable, &towns[i])
		}
		gs.minimapTex = ebiten.NewImageFromImage(gs.world.MinimapImage(gs.player.GetPos(), viewable))
		gs.minimapDirty = false
	}

	b := gs.minimapTex.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(
		float64(viewportRect.Min.X+(viewportRect.Dx()-b.Dx())/2),
		float64(viewportRect.Min.Y+(viewportRect.Dy()-b.Dy())/2),
	)
	screen.DrawImage(gs.minimapTex, op)
}

func (gs *GameState) drawDebugOverlay(screen *ebiten.Image) {
	drawTextBlock(screen, debugOverlayText(gs.player.GetPos(), gs.wind), 8, 8, colorHeading)
}

func main() {
	logger := createLogger()
	logger.Info("Starting...")
	logger.Info(fmt.Sprintf("Window Dimensions %+v", window.Window))
	logger.Info(fmt.Sprintf("Viewable Area %+v", window.ViewPort))

	gs := initGameState(logger)
	gs.startedAt = time.Now()
	gs.lastTickAt = time.Now()
	gs.splashImage = loadSplash()
	gs.refreshActionBar()

	ebiten.SetWindowSize(window.Window.Width, window.Window.Height)
	ebiten.SetWindowTitle("Pirate Wars")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(gs); err != nil {
		log.Fatalf("pirate-wars exited: %v", err)
	}
	logger.Info("Exiting...")
}
