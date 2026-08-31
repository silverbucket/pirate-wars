package main

import (
	"image"
	"time"

	"pirate-wars/cmd/gfx"

	"github.com/hajimehoshi/ebiten/v2"
)

// noticeDuration is how long a transient message stays on screen.
const noticeDuration = 2500 * time.Millisecond

// notice is a short-lived message shown across the top of the viewport.
//
// Several commands used to fail silently — examining with nothing in sight,
// docking away from a town — which leaves the player unable to tell a broken
// key from an unmet condition.
type notice struct {
	text    string
	shownAt time.Time
}

// setNotice replaces any message on screen with text.
func (gs *GameState) setNotice(text string) {
	gs.notice = notice{text: text, shownAt: time.Now()}
}

// activeNotice returns the message still within its display window, or "".
func (gs *GameState) activeNotice() string {
	if gs.notice.text == "" || time.Since(gs.notice.shownAt) > noticeDuration {
		return ""
	}
	return gs.notice.text
}

// drawNotice renders the current message as a banner near the top of the map.
func (gs *GameState) drawNotice(screen *ebiten.Image) {
	text := gs.activeNotice()
	if text == "" {
		return
	}

	w := gfx.TextWidth(text) + 24
	h := gfx.LineHeight + 14
	x := viewportRect.Min.X + (viewportRect.Dx()-w)/2
	y := viewportRect.Min.Y + 16

	rect := image.Rect(x, y, x+w, y+h)
	fillRect(screen, rect, colorNoticeBg)
	strokeRect(screen, rect, colorHeading)
	drawText(screen, text, x+12, y+7, colorHeading)
}
