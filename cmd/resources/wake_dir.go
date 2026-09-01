package resources

import (
	"image"
	"image/color"
	"math"

	"pirate-wars/cmd/common"
)

// DirectionalWakeFrames is the churn cycle length for the directional wake.
const DirectionalWakeFrames = 2

var directionalWakeCache [8][DirectionalWakeFrames]image.Image

// GetDirectionalWake returns the V-shaped wake for a ship on the given
// heading: two foam arms spreading aft from the stern, with churn between
// them. The base art is drawn for a north-facing ship (wake opening
// southward) and rotated like the hull sprites, so the wake always trails
// the ship's true course.
func GetDirectionalWake(facing common.Facing, frame int) image.Image {
	f := int(facing)
	if f < 0 || f >= 8 {
		return nil
	}
	frame = ((frame % DirectionalWakeFrames) + DirectionalWakeFrames) % DirectionalWakeFrames
	if cached := directionalWakeCache[f][frame]; cached != nil {
		return cached
	}

	base := drawWakeNorth(frame)
	img := RotateSprite(base, float64(f)*math.Pi/4)
	directionalWakeCache[f][frame] = img
	return img
}

// drawWakeNorth paints the wake for a north-bound ship on a TileSize canvas.
// The stern sits near the top of the tile (the tile is placed one cell aft of
// the hull), and the V opens downward, widening and fading with distance.
func drawWakeNorth(frame int) *image.NRGBA {
	s := TileSize
	img := image.NewNRGBA(image.Rect(0, 0, s, s))
	cx := float64(s) / 2
	spread := 0.26 // arm slope: px of half-width per px aft

	for y := 0; y < s; y++ {
		// dist runs 0 at the stern line (top of tile) to 1 at full trail;
		// the squared fade dies off well before the tile edge so the wake
		// reads as disturbed water, not a jet stream.
		dist := float64(y) / float64(s)
		halfW := 3 + spread*float64(y)
		fade := (1 - dist) * (1 - dist)
		for x := 0; x < s; x++ {
			dx := math.Abs(float64(x) - cx)
			if dx > halfW+1.5 {
				continue
			}
			var a float64
			edge := halfW - dx
			switch {
			case edge >= -1.5 && edge < 2:
				// Foam arm broken into travelling dashes, with real gaps
				// between them so the trail sparkles rather than streams.
				dash := math.Sin(float64(y)*0.55 + float64(frame)*math.Pi)
				if dash > -0.15 {
					a = fade * (70 + 45*dash)
				}
			case edge >= 2:
				// Sparse churn specks between the arms.
				n := math.Sin(float64(x)*1.1+float64(y)*0.7+float64(frame)*2.1) *
					math.Sin(float64(y)*0.9-float64(x)*0.5)
				if n > 0.72 {
					a = fade * 55 * n
				}
			}
			if a > 4 {
				if a > 130 {
					a = 130
				}
				img.SetNRGBA(x, y, color.NRGBA{R: 235, G: 248, B: 252, A: uint8(a)})
			}
		}
	}
	return img
}

// RotateSprite rotates a square sprite clockwise about its centre, sampling
// bilinearly with alpha-premultiplied weights. Exported for the art-preview
// tooling; the game itself rotates on the GPU.
func RotateSprite(src *image.NRGBA, rad float64) *image.NRGBA {
	b := src.Bounds()
	n := b.Dx()
	dst := image.NewNRGBA(image.Rect(0, 0, n, n))
	sin, cos := math.Sin(rad), math.Cos(rad)
	half := float64(n) / 2
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			dx := float64(x) + 0.5 - half
			dy := float64(y) + 0.5 - half
			sx := dx*cos + dy*sin + half
			sy := -dx*sin + dy*cos + half
			dst.SetNRGBA(x, y, bilinearNRGBA(src, sx, sy))
		}
	}
	return dst
}

func bilinearNRGBA(src *image.NRGBA, fx, fy float64) color.NRGBA {
	b := src.Bounds()
	n := b.Dx()
	x0 := int(math.Floor(fx - 0.5))
	y0 := int(math.Floor(fy - 0.5))
	wx := (fx - 0.5) - float64(x0)
	wy := (fy - 0.5) - float64(y0)

	var r, g, bl, a float64
	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			x, y := x0+dx, y0+dy
			if x < 0 || y < 0 || x >= n || y >= n {
				continue
			}
			w := (1 - math.Abs(float64(dx)-wx)) * (1 - math.Abs(float64(dy)-wy))
			p := src.NRGBAAt(b.Min.X+x, b.Min.Y+y)
			aw := w * float64(p.A)
			r += float64(p.R) * aw
			g += float64(p.G) * aw
			bl += float64(p.B) * aw
			a += aw
		}
	}
	if a < 1 {
		return color.NRGBA{}
	}
	return color.NRGBA{
		R: uint8(r/a + 0.5), G: uint8(g/a + 0.5), B: uint8(bl/a + 0.5),
		A: uint8(a + 0.5),
	}
}
