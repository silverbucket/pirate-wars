package main

import (
	"image"
	"image/color"
	"math"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/gfx"
	"pirate-wars/cmd/sailing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	// colorHeadingNeedle matches the bow chevron drawn on the player's ship, so
	// the same cyan means "where you are pointed" in both places.
	colorHeadingNeedle = color.RGBA{R: 90, G: 230, B: 255, A: 255}
	colorWindNeedle    = color.RGBA{R: 240, G: 214, B: 140, A: 255}
	colorCompassFace   = color.RGBA{R: 26, G: 30, B: 38, A: 255}
)

// compassSize is the diameter of the side panel compass in pixels.
const compassSize = 108

// facingVector returns the unit screen-space direction for a facing, with
// diagonals normalised so all eight needles are the same length.
func facingVector(f common.Facing) (float64, float64) {
	d := common.FacingToDelta(f)
	x, y := float64(d.X), float64(d.Y)
	if l := math.Hypot(x, y); l > 0 {
		x, y = x/l, y/l
	}
	return x, y
}

// drawCompass draws the wind and heading needles on one dial.
//
// Point of sail is the angle between them and it drives a 45x speed range, but
// working it out from two 8-way sprites on a 32px tile is not something a player
// can do at a glance. On one dial the angle is the picture: needles together is
// a run, opposed is in irons.
func drawCompass(dst *ebiten.Image, bounds image.Rectangle, heading common.Facing, wind *sailing.Wind) {
	cx := float32(bounds.Min.X + bounds.Dx()/2)
	cy := float32(bounds.Min.Y + bounds.Dy()/2)
	r := float32(bounds.Dx()) / 2

	vector.DrawFilledCircle(dst, cx, cy, r, colorCompassFace, true)
	vector.StrokeCircle(dst, cx, cy, r, 1, colorPanelEdge, true)

	// Octant ticks, so the eight headings are countable rather than estimated.
	for i := 0; i < 8; i++ {
		vx, vy := facingVector(common.Facing(i))
		vector.StrokeLine(dst,
			cx+float32(vx)*(r-4), cy+float32(vy)*(r-4),
			cx+float32(vx)*r, cy+float32(vy)*r,
			1, colorPanelEdge, true)
	}
	drawText(dst, "N", int(cx)-3, bounds.Min.Y-gfx.LineHeight+2, colorTextDim)

	if wind != nil {
		// The wind needle points downwind, matching the pennant on the map.
		wx, wy := facingVector(wind.Facing)
		vector.StrokeLine(dst,
			cx-float32(wx)*(r-8), cy-float32(wy)*(r-8),
			cx+float32(wx)*(r-8), cy+float32(wy)*(r-8),
			3, colorWindNeedle, true)
		drawNeedleTip(dst, cx+float32(wx)*(r-8), cy+float32(wy)*(r-8), 4, colorWindNeedle)
	}

	hx, hy := facingVector(heading)
	vector.StrokeLine(dst, cx, cy, cx+float32(hx)*(r-8), cy+float32(hy)*(r-8), 3, colorHeadingNeedle, true)
	drawNeedleTip(dst, cx+float32(hx)*(r-8), cy+float32(hy)*(r-8), 5, colorHeadingNeedle)
	vector.DrawFilledCircle(dst, cx, cy, 3, colorHeadingNeedle, true)
}

// drawNeedleTip marks the pointing end of a needle with a bulb, so the two
// needles read as directions rather than as a bare cross.
func drawNeedleTip(dst *ebiten.Image, x, y, size float32, c color.Color) {
	vector.DrawFilledCircle(dst, x, y, size, c, true)
}
