package resources

import (
	"image"
	"image/color"
	"math"

	"pirate-wars/cmd/common"
)

// Bow marker colours. Cyan appears nowhere in the tileset palette — not in the
// blue water, the brown hulls, or the cream sails — so the marker stays legible
// wherever the ship happens to be.
var (
	bowFill    = color.RGBA{R: 90, G: 230, B: 255, A: 255}
	bowOutline = color.RGBA{R: 10, G: 20, B: 30, A: 255}
)

type bowKey struct {
	size   int
	facing common.Facing
}

var bowMarkerCache = map[bowKey]image.Image{}

// GetBowMarkerOverlay returns a chevron drawn at the bow of a ship facing f.
//
// The 8-way hull sprites are close to mirror images of each other — east and
// west differ by under a fifth of their pixels — so relative steering needs an
// unambiguous heading cue that does not depend on reading the art. The marker is
// generated in code rather than authored, so it stays correct for any cell size
// and cannot drift out of sync with the tileset.
func GetBowMarkerOverlay(size int, f common.Facing) image.Image {
	key := bowKey{size: size, facing: f}
	if cached, ok := bowMarkerCache[key]; ok {
		return cached
	}

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	d := common.FacingToDelta(f)
	if d.X == 0 && d.Y == 0 {
		bowMarkerCache[key] = img
		return img
	}

	// Unit direction, normalised so diagonals do not overshoot the tile corner.
	dx, dy := float64(d.X), float64(d.Y)
	if d.X != 0 && d.Y != 0 {
		dx *= 0.7071
		dy *= 0.7071
	}

	half := float64(size) / 2
	reach := half - float64(size)/16
	tipX, tipY := half+dx*reach, half+dy*reach

	// Base sits behind the tip, spread perpendicular to the heading.
	length := float64(size) / 5
	width := float64(size) / 7
	baseX, baseY := tipX-dx*length, tipY-dy*length
	px, py := -dy, dx

	tri := [3][2]float64{
		{tipX, tipY},
		{baseX + px*width, baseY + py*width},
		{baseX - px*width, baseY - py*width},
	}

	// Outline first, then the fill inset by a pixel, so the chevron keeps its
	// edge over both the pale sails and the dark water.
	fillTriangle(img, expandTriangle(tri, 1.6), bowOutline)
	fillTriangle(img, tri, bowFill)

	bowMarkerCache[key] = img
	return img
}

// expandTriangle pushes each vertex out from the centroid by grow pixels.
func expandTriangle(t [3][2]float64, grow float64) [3][2]float64 {
	cx := (t[0][0] + t[1][0] + t[2][0]) / 3
	cy := (t[0][1] + t[1][1] + t[2][1]) / 3
	var out [3][2]float64
	for i, v := range t {
		dx, dy := v[0]-cx, v[1]-cy
		l := dx*dx + dy*dy
		if l == 0 {
			out[i] = v
			continue
		}
		l = math.Sqrt(l)
		out[i] = [2]float64{v[0] + dx/l*grow, v[1] + dy/l*grow}
	}
	return out
}

// fillTriangle rasterises t into img with a half-plane test per pixel.
func fillTriangle(img *image.RGBA, t [3][2]float64, c color.Color) {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			d0 := edge(t[0], t[1], fx, fy)
			d1 := edge(t[1], t[2], fx, fy)
			d2 := edge(t[2], t[0], fx, fy)
			if (d0 >= 0 && d1 >= 0 && d2 >= 0) || (d0 <= 0 && d1 <= 0 && d2 <= 0) {
				img.Set(x, y, c)
			}
		}
	}
}

// edge is the signed area of the triangle (a, b, p), used as a half-plane test.
func edge(a, b [2]float64, px, py float64) float64 {
	return (b[0]-a[0])*(py-a[1]) - (b[1]-a[1])*(px-a[0])
}
