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
	label  string
	rect   image.Rectangle
	action func()
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
)

func fillRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	vector.DrawFilledRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), c, false)
}

func strokeRect(dst *ebiten.Image, r image.Rectangle, c color.Color) {
	vector.StrokeRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), 1, c, false)
}

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

func drawButton(dst *ebiten.Image, b button) {
	fillRect(dst, b.rect, colorButton)
	strokeRect(dst, b.rect, colorButtonEdge)
	tx := b.rect.Min.X + (b.rect.Dx()-gfx.TextWidth(b.label))/2
	ty := b.rect.Min.Y + (b.rect.Dy()-gfx.LineHeight)/2
	drawText(dst, b.label, tx, ty, colorText)
}

// buttonRow lays buttons out left to right inside bounds, wrapping is not needed
// because the action bar is sized for the widest view.
func buttonRow(bounds image.Rectangle, labels []string, actions []func()) []button {
	buttons := make([]button, 0, len(labels))
	x := bounds.Min.X + buttonGap
	y := bounds.Min.Y + actionBarLabelHeight + (bounds.Dy()-actionBarLabelHeight-buttonHeight)/2
	for i, label := range labels {
		w := gfx.TextWidth(label) + buttonPadding*2
		if x+w > bounds.Max.X-buttonGap {
			break
		}
		buttons = append(buttons, button{
			label:  label,
			rect:   image.Rect(x, y, x+w, y+buttonHeight),
			action: actions[i],
		})
		x += w + buttonGap
	}
	return buttons
}
