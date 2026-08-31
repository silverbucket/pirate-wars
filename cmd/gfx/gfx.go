// Package gfx bridges the game's image.Image tile caches to Ebiten textures
// and owns the UI text face.
package gfx

import (
	"image"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var textureCache = map[image.Image]*ebiten.Image{}

// Texture returns the Ebiten texture for img, uploading it on first use.
// Callers must pass cached (pointer-stable) images; the tile caches in
// cmd/resources and cmd/harbor guarantee this.
func Texture(img image.Image) *ebiten.Image {
	if img == nil {
		return nil
	}
	if tex, ok := textureCache[img]; ok {
		return tex
	}
	tex := ebiten.NewImageFromImage(img)
	textureCache[img] = tex
	return tex
}

// bitmapfont covers the arrow and interpunct glyphs the action bar legend uses
// for key hints, which a plain ASCII face would render as tofu boxes.
var uiFace = text.NewGoXFace(bitmapfont.Face)

// Face is the bitmap face used for all in-game text.
func Face() text.Face { return uiFace }

// LineHeight is the vertical advance for one line of UI text.
const LineHeight = 14

// TextWidth measures a single line in pixels.
func TextWidth(s string) int {
	w, _ := text.Measure(s, uiFace, LineHeight)
	return int(w)
}
