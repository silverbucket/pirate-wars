package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/economy"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/gfx"
	"pirate-wars/cmd/hail"
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/player"
	"pirate-wars/cmd/sailing"
	"pirate-wars/cmd/town"
	"pirate-wars/cmd/window"
	"pirate-wars/cmd/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"go.uber.org/zap"
)

const BASE_LOG_LEVEL = zap.DebugLevel
const DEV_MODE = true

const splashDuration = 2 * time.Second

var ViewType = world.ViewTypeMainMap

// Screen regions. The viewport keeps the historic 854×728 size; the action bar
// is drawn as an opaque strip over its bottom edge.
//
// The bar spans the full window width and the side panel stops above it, rather
// than the bar stopping at the panel. The panel's lower half is empty below the
// compass, and the bar was within 2px of overflowing once Hail joined it —
// buttonRow drops anything that does not fit, silently.
var (
	viewportRect  = image.Rect(0, 0, window.ViewPort.Dimensions.Width, window.ViewPort.Dimensions.Height)
	actionBarRect = image.Rect(0, window.Window.Height-window.ActionMenu.Height, window.Window.Width, window.Window.Height)
	sidePanelRect = image.Rect(window.Window.Width-window.SidePanel.Width, 0, window.Window.Width, actionBarRect.Min.Y)
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

	buttons     []button
	startedAt   time.Time
	lastTickAt  time.Time
	splashImage *ebiten.Image
	minimapTex  *ebiten.Image

	// notice is the transient banner shown for commands that would otherwise
	// fail silently.
	notice notice
	// turnedThisTick rations the helm to one octant per sailing tick.
	turnedThisTick bool
	// alongsideNpcID is the ship the player has come up against, offered as a
	// Hail on the action bar; lastHailedNpcID stops the prompt repeating for a
	// ship already hailed.
	alongsideNpcID  string
	lastHailedNpcID string
	// quitting ends the run loop after the Ctrl+Q confirmation.
	quitting       bool
	quitReturnView int
	// forceDockable lets tests exercise the bar with a town in reach without
	// generating a world to stand beside.
	forceDockable bool
}

// initGameState builds the world, towns, NPCs and player for a new voyage.
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
	gs.initialized = true
	return &gs
}

// processTick advances one sailing tick: the clock, ship movement, and the
// terrain animation. Only the main map moves ships, so the overlays are pauses.
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
}

// visibleNPCs adapts the NPC list to the read-only view the renderer takes.
func (gs *GameState) visibleNPCs() []entities.AvatarReadOnly {
	list := gs.npcs.GetList()
	visible := make([]entities.AvatarReadOnly, 0, len(list))
	for i := range list {
		visible = append(visible, &list[i])
	}
	return visible
}

// Update implements ebiten.Game.
func (gs *GameState) Update() error {
	if gs.showingSplash() {
		// Any key or click skips the title card. A fixed unskippable wait is a
		// toll charged on every launch, and it is paid most often by whoever is
		// restarting the game the most.
		if len(pressedKeyNames()) > 0 || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			gs.startedAt = time.Now().Add(-splashDuration)
		}
		return nil
	}
	gs.paused = false

	if gs.buttons == nil {
		gs.refreshActionBar()
	}
	gs.handleInput()

	if gs.quitting {
		// Termination unwinds Ebiten cleanly and lets main flush the logger,
		// which the old os.Exit in the quit handler skipped.
		return ebiten.Termination
	}

	tick := gs.sailingCfg.TickDuration()
	if time.Since(gs.lastTickAt) >= tick {
		gs.lastTickAt = time.Now()
		gs.processTick()
		gs.refreshActionBar()
	}
	gs.applyHover()
	return nil
}

// showingSplash reports whether the title card still owns the screen.
func (gs *GameState) showingSplash() bool {
	return time.Since(gs.startedAt) < splashDuration
}

// Draw implements ebiten.Game.
func (gs *GameState) Draw(screen *ebiten.Image) {
	screen.Fill(colorPanel)

	if gs.showingSplash() {
		gs.drawSplash(screen)
		return
	}

	viewport := screen.SubImage(viewportRect).(*ebiten.Image)
	highlight := ExamineData.GetFocusedEntity()
	gs.world.Draw(viewport, gs.player, gs.visibleNPCs(), highlight, gs.wind.Facing)

	gs.drawSidePanel(screen, highlight)
	gs.drawActionBar(screen)

	if ViewType == world.ViewTypeMiniMap {
		gs.drawMinimap(screen)
	}
	if hasOverlay() {
		gs.drawOverlay(screen)
	}
	gs.drawNotice(screen)
	if debugOverlayVisible {
		gs.drawDebugOverlay(screen)
	}
}

// Layout implements ebiten.Game.
func (gs *GameState) Layout(int, int) (int, int) {
	return window.Window.Width, window.Window.Height
}

// drawSplash paints the title card, or a text fallback when the art is missing.
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
	hint := "press any key"
	drawText(screen, hint, (window.Window.Width-gfx.TextWidth(hint))/2, window.Window.Height/2+gfx.LineHeight*2, colorTextDim)
}

// drawSidePanel renders the ship's instruments: the sailing readout, the compass
// showing heading against wind, and the examine block when something is focused.
func (gs *GameState) drawSidePanel(screen *ebiten.Image, examine entities.ViewableEntity) {
	fillRect(screen, sidePanelRect, colorPanel)
	strokeRect(screen, sidePanelRect, colorPanelEdge)

	x := sidePanelRect.Min.X + 8
	y := 10
	drawText(screen, "Ship Status", x, y, colorHeading)
	y += gfx.LineHeight + 4

	status := gs.newShipStatus()
	for _, line := range status.statusLines() {
		c := colorText
		if line.warn {
			c = colorWarn
		}
		drawText(screen, line.label, x, y, colorTextDim)
		drawText(screen, line.value, x+48, y, c)
		y += gfx.LineHeight
	}

	if reason := status.stallReason(); reason != "" {
		y += 4
		for _, line := range wrapText(reason, 25) {
			drawText(screen, line, x, y, colorWarn)
			y += gfx.LineHeight
		}
	}

	y += 10 + gfx.LineHeight
	compass := image.Rect(x+16, y, x+16+compassSize, y+compassSize)
	drawCompass(screen, compass, status.Heading, gs.wind)
	y = compass.Max.Y + gfx.LineHeight

	drawText(screen, "Heading", x, y, colorHeadingNeedle)
	drawText(screen, "Wind", x+80, y, colorWindNeedle)
	y += gfx.LineHeight + 6

	// The examine block used to render its headings over empty values whenever
	// nothing was focused, spending four lines of a 170px panel on blanks.
	if examine.GetName() != "" {
		drawText(screen, "Examine", x, y, colorHeading)
		y += gfx.LineHeight + 4
		drawTextBlock(screen, examinePanelText(examine), x, y, colorTextDim)
	}
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

// drawMinimap renders the world chart: the cached terrain raster, town markers,
// the player, and an outline of the area currently on screen.
func (gs *GameState) drawMinimap(screen *ebiten.Image) {
	// The terrain never changes, so the raster uploads once for the whole run.
	if gs.minimapTex == nil {
		gs.minimapTex = ebiten.NewImageFromImage(gs.world.MinimapBase())
	}

	b := gs.minimapTex.Bounds()
	originX := viewportRect.Min.X + (viewportRect.Dx()-b.Dx())/2
	originY := viewportRect.Min.Y + (viewportRect.Dy()-b.Dy())/2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(originX), float64(originY))
	screen.DrawImage(gs.minimapTex, op)

	mapRect := image.Rect(originX, originY, originX+b.Dx(), originY+b.Dy())
	strokeRect(screen, mapRect, colorPanelEdge)

	cw, ch := gs.world.MinimapCellSize()
	marker := func(pos common.Coordinates, size int, c color.Color) {
		x := originX + int(float32(pos.X)*cw)
		y := originY + int(float32(pos.Y)*ch)
		fillRect(screen, image.Rect(x-size/2, y-size/2, x-size/2+size, y-size/2+size), c)
	}

	for _, t := range gs.towns.GetTowns() {
		marker(t.GetPos(), 5, t.GetColor())
	}

	// The chart is otherwise scale-less: without this the player cannot tell how
	// much of the world one screen covers, or which way the view extends.
	vpr := window.GetViewportRegion(gs.player.GetPos())
	strokeRect(screen, image.Rect(
		originX+int(float32(vpr.X)*cw),
		originY+int(float32(vpr.Y)*ch),
		originX+int(float32(vpr.X+vpr.Cols)*cw),
		originY+int(float32(vpr.Y+vpr.Rows)*ch),
	), colorHeadingNeedle)

	marker(gs.player.GetPos(), 7, color.White)

	drawText(screen, "Your ship", mapRect.Min.X+6, mapRect.Max.Y+6, colorText)
	drawText(screen, "Towns", mapRect.Min.X+90, mapRect.Max.Y+6, colorHeading)
	drawText(screen, "On screen", mapRect.Min.X+160, mapRect.Max.Y+6, colorHeadingNeedle)
}

// drawDebugOverlay prints window, viewport and world geometry behind F3.
func (gs *GameState) drawDebugOverlay(screen *ebiten.Image) {
	drawTextBlock(screen, debugOverlayText(gs.player.GetPos(), gs.wind), 8, 8, colorHeading)
}

// main boots the game and hands the loop to Ebiten.
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

	if err := ebiten.RunGame(gs); err != nil && !errors.Is(err, ebiten.Termination) {
		log.Fatalf("pirate-wars exited: %v", err)
	}
	logger.Info("Exiting...")
}
