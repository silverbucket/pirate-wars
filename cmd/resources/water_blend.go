package resources

import (
	"image"
	"image/color"
	"image/draw"

	"pirate-wars/cmd/common"
)

// WaterContext is the terrain type on each side of a water cell.
type WaterContext struct {
	N, E, S, W common.TerrainType
}

type waterBlendKey struct {
	tt    common.TerrainType
	frame int
	ctx   WaterContext
}

var (
	waterBlendCache = map[waterBlendKey]image.Image{}
	waterMeanCache  = map[common.TerrainType]color.NRGBA{}
)

// ResetWaterBlendCache drops the blended tiles after a tileset override.
func ResetWaterBlendCache() {
	waterBlendCache = map[waterBlendKey]image.Image{}
	waterMeanCache = map[common.TerrainType]color.NRGBA{}
}

// GetBlendedWaterTile returns the water tile for tt with its edges pulled
// toward any neighbouring water of a different depth. The blend band is
// dithered rather than smooth, so the transition keeps the sheet's pixel
// grain instead of reading as an airbrushed gradient. Uniform surroundings
// return the plain animated tile.
func GetBlendedWaterTile(tt common.TerrainType, frame int, ctx WaterContext) image.Image {
	blendN := IsWaterTerrain(ctx.N) && ctx.N != tt
	blendE := IsWaterTerrain(ctx.E) && ctx.E != tt
	blendS := IsWaterTerrain(ctx.S) && ctx.S != tt
	blendW := IsWaterTerrain(ctx.W) && ctx.W != tt
	if !blendN && !blendE && !blendS && !blendW {
		return GetWaterTile(tt, frame)
	}

	// Normalise non-blending sides so equivalent cells share a cache entry.
	if !blendN {
		ctx.N = tt
	}
	if !blendE {
		ctx.E = tt
	}
	if !blendS {
		ctx.S = tt
	}
	if !blendW {
		ctx.W = tt
	}
	key := waterBlendKey{tt: tt, frame: frame, ctx: ctx}
	if cached, ok := waterBlendCache[key]; ok {
		return cached
	}

	base := GetWaterTile(tt, frame)
	img := image.NewRGBA(image.Rect(0, 0, TileSize, TileSize))
	draw.Draw(img, img.Bounds(), base, base.Bounds().Min, draw.Src)

	const band = 16.0
	for y := 0; y < TileSize; y++ {
		for x := 0; x < TileSize; x++ {
			var weight float64
			var toward color.NRGBA
			consider := func(on bool, dist float64, other common.TerrainType) {
				if !on || dist >= band {
					return
				}
				w := (1 - dist/band) * 0.55
				if w > weight {
					weight = w
					toward = waterMeanColor(other)
				}
			}
			consider(blendN, float64(y), ctx.N)
			consider(blendS, float64(TileSize-1-y), ctx.S)
			consider(blendW, float64(x), ctx.W)
			consider(blendE, float64(TileSize-1-x), ctx.E)
			if weight == 0 {
				continue
			}
			// Checker dither keeps the band textured.
			if (x+y)%2 == 1 {
				weight *= 0.4
			}
			p := img.RGBAAt(x, y)
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(p.R)*(1-weight) + float64(toward.R)*weight),
				G: uint8(float64(p.G)*(1-weight) + float64(toward.G)*weight),
				B: uint8(float64(p.B)*(1-weight) + float64(toward.B)*weight),
				A: p.A,
			})
		}
	}
	waterBlendCache[key] = img
	return img
}

// waterMeanColor is the average colour of a water type's base tile.
func waterMeanColor(tt common.TerrainType) color.NRGBA {
	if cached, ok := waterMeanCache[tt]; ok {
		return cached
	}
	tile := GetWaterTile(tt, 0)
	b := tile.Bounds()
	var r, g, bl, n float64
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			pr, pg, pb, _ := tile.At(x, y).RGBA()
			r += float64(pr >> 8)
			g += float64(pg >> 8)
			bl += float64(pb >> 8)
			n++
		}
	}
	c := color.NRGBA{R: uint8(r / n), G: uint8(g / n), B: uint8(bl / n), A: 255}
	waterMeanCache[tt] = c
	return c
}
