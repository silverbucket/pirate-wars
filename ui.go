package main

import (
	"image"
	"image/color"
	"strings"

	"pirate-wars/cmd/gfx"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// button is a tap-friendly text button drawn by the Ebiten UI.
type button struct {
	label string
	rect  image.Rectangle
	// enabled false keeps the button on the bar but greys it out and ignores
	// taps, so an unavailable command does not reflow the row.
	enabled        bool
	disabledReason string
	hovered        bool
	action         func()
}

const (
	buttonHeight  = 34
	buttonPadding = 12
	buttonGap     = 6
	// Leaves room for the view label above the action bar buttons.
	actionBarLabelHeight = 20
)

var (
	colorPanel      = color.RGBA{R: 18, G: 20, B: 26, A: 255}
	colorPanelEdge  = color.RGBA{R: 90, G: 96, B: 110, A: 255}
	colorButton     = color.RGBA{R: 42, G: 48, B: 62, A: 255}
	colorButtonEdge = color.RGBA{R: 130, G: 140, B: 160, A: 255}
	colorScrim      = color.RGBA{R: 0, G: 0, B: 0, A: 205}
	colorText       = color.RGBA{R: 232, G: 232, B: 236, A: 255}
	colorTextDim    = color.RGBA{R: 150, G: 152, B: 160, A: 255}
	colorHeading    = color.RGBA{R: 240, G: 214, B: 140, A: 255}
	colorNoticeBg   = color.RGBA{R: 28, G: 32, B: 42, A: 240}
	colorWarn       = color.RGBA{R: 240, G: 180, B: 110, A: 255}

	// Hover and disabled fills. A pointer needs to know a rectangle is a target
	// before it clicks, and which targets are live.
	colorButtonHover     = color.RGBA{R: 64, G: 74, B: 96, A: 255}
	colorButtonHoverEdge = color.RGBA{R: 190, G: 200, B: 220, A: 255}
	colorButtonOff       = color.RGBA{R: 26, G: 29, B: 36, A: 255}
	colorButtonOffEdge   = color.RGBA{R: 68, G: 72, B: 84, A: 255}
	colorTextOff         = color.RGBA{R: 104, G: 108, B: 118, A: 255}
)

// fillRect paints a solid rectangle.
func fillRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	vector.DrawFilledRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), c, false)
}

// strokeRect outlines a rectangle with a one-pixel border.
func strokeRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	vector.StrokeRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), 1, c, false)
}

// drawText draws one line of UI text with its top-left at (x, y).
func drawText(dst *ebiten.Image, s string, x, y int, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, s, gfx.Face(), op)
}

// drawTextBlock draws multi-line text and returns the y below the last line.
func drawTextBlock(dst *ebiten.Image, s string, x, y int, c color.Color) int {
	for _, line := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		drawText(dst, line, x, y, c)
		y += gfx.LineHeight
	}
	return y
}

// drawButton renders a button in its normal, hovered or disabled state.
func drawButton(dst *ebiten.Image, b button) {
	fill, edge, text := colorButton, colorButtonEdge, colorText
	switch {
	case !b.enabled:
		fill, edge, text = colorButtonOff, colorButtonOffEdge, colorTextOff
	case b.hovered:
		fill, edge = colorButtonHover, colorButtonHoverEdge
	}

	fillRect(dst, b.rect, fill)
	strokeRect(dst, b.rect, edge)
	tx := b.rect.Min.X + (b.rect.Dx()-gfx.TextWidth(b.label))/2
	ty := b.rect.Min.Y + (b.rect.Dy()-gfx.LineHeight)/2
	drawText(dst, b.label, tx, ty, text)
}

// barSpec is one command's worth of action bar button.
type barSpec struct {
	label          string
	enabled        bool
	disabledReason string
	action         func()
}

// buttonRow lays buttons out left to right inside bounds. TestActionBarFitsAction
// MenuArea asserts every view's row fits, so the overflow break here is a
// backstop rather than the plan.
func buttonRow(bounds image.Rectangle, specs []barSpec) []button {
	buttons := make([]button, 0, len(specs))
	x := bounds.Min.X + buttonGap
	y := bounds.Min.Y + actionBarLabelHeight + (bounds.Dy()-actionBarLabelHeight-buttonHeight)/2
	for _, spec := range specs {
		w := gfx.TextWidth(spec.label) + buttonPadding*2
		if x+w > bounds.Max.X-buttonGap {
			break
		}
		buttons = append(buttons, button{
			label:          spec.label,
			rect:           image.Rect(x, y, x+w, y+buttonHeight),
			enabled:        spec.enabled,
			disabledReason: spec.disabledReason,
			action:         spec.action,
		})
		x += w + buttonGap
	}
	return buttons
}

// applyHover marks the button under the pointer and switches the cursor, so a
// mouse gets the same affordance feedback the keyboard gets from the legend.
func (gs *GameState) applyHover() {
	x, y := ebiten.CursorPosition()
	shape := ebiten.CursorShapeDefault
	for i := range gs.buttons {
		b := &gs.buttons[i]
		b.hovered = x >= b.rect.Min.X && x < b.rect.Max.X && y >= b.rect.Min.Y && y < b.rect.Max.Y
		if b.hovered && b.enabled {
			shape = ebiten.CursorShapePointer
		}
	}
	ebiten.SetCursorShape(shape)
}
