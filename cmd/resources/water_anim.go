package resources

import (
	"image"
	"image/color"
	"math"

	"pirate-wars/cmd/common"
)

// waterFrameCache holds the generated twinkle frames for deep and open water,
// keyed by terrain type and frame, so image identity stays stable for the
// Ebiten texture cache.
var waterFrameCache = map[[2]int]image.Image{}

// ResetWaterFrameCache drops the generated frames after a tileset override.
func ResetWaterFrameCache() {
	waterFrameCache = map[[2]int]image.Image{}
}

// GetWaterTile returns the animated tile for any water terrain. Shallow water
// keeps its authored wave frames from the tileset; deep and open water twinkle:
// their sparkle highlights are dimmed in rotating thirds so the open sea
// glitters instead of sitting still.
func GetWaterTile(tt common.TerrainType, frame int) image.Image {
	if tt == common.TerrainTypeShallowWater {
		return GetWaveTile(frame)
	}
	if !IsWaterTerrain(tt) {
		return GetTerrainTile(tt)
	}

	frame = ((frame % WaveFrameCount) + WaveFrameCount) % WaveFrameCount
	key := [2]int{int(tt), frame}
	if cached, ok := waterFrameCache[key]; ok {
		return cached
	}

	img := twinkleFrame(GetTerrainTile(tt), frame, WaveFrameCount)
	waterFrameCache[key] = img
	return img
}

// twinkleFrame copies base and fades the bright sparkle pixels whose group
// matches this frame toward the water's mean brightness. The same scheme
// generates the shallow-water frames in the tileset (scripts/build-tileset64),
// so all three waters twinkle at the same rhythm.
func twinkleFrame(base image.Image, frame, frameCount int) image.Image {
	b := base.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	var lumSum, lumMax float64
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, bl, _ := base.At(b.Min.X+x, b.Min.Y+y).RGBA()
			l := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(bl>>8)
			lumSum += l
			if l > lumMax {
				lumMax = l
			}
		}
	}
	mean := lumSum / float64(b.Dx()*b.Dy())
	// Adaptive, matching scripts/build-tileset64: only clearly-brighter-than-
	// the-sea pixels twinkle, whatever the tile's overall brightness.
	threshold := mean + math.Max(25, (lumMax-mean)*0.45)

	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, bl, a := base.At(b.Min.X+x, b.Min.Y+y).RGBA()
			r8, g8, b8 := float64(r>>8), float64(g>>8), float64(bl>>8)
			lum := 0.299*r8 + 0.587*g8 + 0.114*b8
			if lum >= threshold && (x/3+y/3+frame)%frameCount == 0 {
				// Fade this sparkle toward the mean water brightness.
				const keep = 0.35
				scale := keep + (1-keep)*mean/lum
				r8, g8, b8 = r8*scale, g8*scale, b8*scale
			}
			img.SetRGBA(x, y, color.RGBA{R: uint8(r8), G: uint8(g8), B: uint8(b8), A: uint8(a >> 8)})
		}
	}
	return img
}
