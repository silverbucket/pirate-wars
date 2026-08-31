// Command render-preview composites sample scenes from the live resources
// package, so tileset and animation changes can be reviewed without launching
// the game.
//
// Static mode writes one PNG per wave frame. With -animate it simulates a few
// sailing ticks — gliding ships, swell roll and bob, directional wakes,
// lapping coast foam, twinkling water — and writes a numbered frame sequence
// ready for `magick -delay 6 seq-*.png out.gif`.
//
// Usage:
//
//	go run ./scripts/render-preview [-tileset sheet.png] [-out dir] [-animate]
package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/resources"
)

const cell = resources.TileSize

// D=deep, O=open, S=shallow, B=beach, L=lowland, H=highland, R=rock, P=peak,
// T=town, G=ghost town, U=brush.
var scene = []string{
	"DDDOOOSSSSSSS",
	"DDOOOSSBBLLLS",
	"DOOSSBBLLTLLS",
	"DOOSBBLLLLLUS",
	"DOOSSBLHHLUSS",
	"DDOOSSBBLLSSS",
	"DDDOOSSSSSSSD",
	"DDDDOOOOOSSDD",
}

var terrainFor = map[byte]common.TerrainType{
	'D': common.TerrainTypeDeepWater, 'O': common.TerrainTypeOpenWater,
	'S': common.TerrainTypeShallowWater, 'B': common.TerrainTypeBeach,
	'L': common.TerrainTypeLowland, 'H': common.TerrainTypeHighland,
	'R': common.TerrainTypeRock, 'P': common.TerrainTypePeak,
	'T': common.TerrainTypeTown, 'G': common.TerrainTypeGhostTown,
	'U': common.TerrainTypeLowlandBrush,
}

type shipSpec struct {
	x, y   int
	ship   common.ShipType
	facing common.Facing
	moving bool
	phase  float64
}

// The moving ships' courses stay on water for the four simulated ticks.
var ships = []shipSpec{
	{1, 6, common.ShipWhite, common.FacingE, true, 0},
	{5, 0, common.ShipPirate, common.FacingE, true, 2.1},
	{10, 7, common.ShipRed, common.FacingW, true, 4.2},
	{1, 2, common.ShipBlue, common.FacingS, true, 1.3},
	{8, 5, common.ShipYellow, common.FacingSE, false, 3.4},
}

func typeAt(x, y int, fallback common.TerrainType) common.TerrainType {
	if y < 0 || y >= len(scene) || x < 0 || x >= len(scene[y]) {
		return fallback
	}
	return terrainFor[scene[y][x]]
}

func isWater(x, y int) bool {
	return resources.IsWaterTerrain(typeAt(x, y, common.TerrainTypeBeach))
}

func main() {
	tileset := flag.String("tileset", "", "optional tileset PNG to preview instead of the bundled one")
	outdir := flag.String("out", "/tmp/pw-preview", "output directory")
	animate := flag.Bool("animate", false, "write an animation frame sequence instead of stills")
	flag.Parse()

	if *tileset != "" {
		f, err := os.Open(*tileset)
		if err != nil {
			fatal(err)
		}
		img, err := png.Decode(f)
		f.Close()
		if err != nil {
			fatal(err)
		}
		resources.OverrideTileset(img)
	}
	if err := os.MkdirAll(*outdir, 0o755); err != nil {
		fatal(err)
	}

	if *animate {
		writeAnimation(*outdir)
		return
	}
	for frame := 0; frame < 3; frame++ {
		writeFrame(filepath.Join(*outdir, fmt.Sprintf("frame%d.png", frame)),
			frame, float64(frame), 1, 0)
	}
}

// writeAnimation simulates ticks of sailing and writes one PNG per subframe.
func writeAnimation(outdir string) {
	const (
		ticks      = 4
		subPerTick = 8
		tickSecs   = 0.63
	)
	for tick := 0; tick < ticks; tick++ {
		for sub := 0; sub < subPerTick; sub++ {
			t := float64(sub) / subPerTick
			simTime := (float64(tick) + t) * tickSecs
			writeFrame(
				filepath.Join(outdir, fmt.Sprintf("seq-%03d.png", tick*subPerTick+sub)),
				tick%3, simTime, t, tick,
			)
		}
	}
}

// writeFrame composites one full scene: terrain with blended water and
// lapping coast, wakes, then hulls riding the swell.
func writeFrame(path string, waveFrame int, simTime, tickT float64, ticksDone int) {
	img := image.NewRGBA(image.Rect(0, 0, len(scene[0])*cell, len(scene)*cell))

	for y, row := range scene {
		for x := range row {
			tt := terrainFor[row[x]]
			var tile image.Image
			if resources.IsWaterTerrain(tt) {
				tile = resources.GetBlendedWaterTile(tt, waveFrame, resources.WaterContext{
					N: typeAt(x, y-1, tt), E: typeAt(x+1, y, tt),
					S: typeAt(x, y+1, tt), W: typeAt(x-1, y, tt),
				})
			} else if resources.IsCoastLandTerrain(tt) {
				tile = resources.GetCoastTileFrame(resources.WaterNeighbors{
					N: isWater(x, y-1), E: isWater(x+1, y),
					S: isWater(x, y+1), W: isWater(x-1, y),
				}, waveFrame)
			}
			if tile == nil {
				tile = resources.GetTerrainTile(tt)
			}
			blitAt(img, tile, float64(x*cell), float64(y*cell))
		}
	}

	// Visual positions after ticksDone whole steps plus the current glide.
	visual := func(s shipSpec) (float64, float64) {
		if !s.moving {
			return float64(s.x), float64(s.y)
		}
		d := common.FacingToDelta(s.facing)
		steps := float64(ticksDone)
		return float64(s.x) + float64(d.X)*(steps+tickT), float64(s.y) + float64(d.Y)*(steps+tickT)
	}

	for _, s := range ships {
		if !s.moving {
			continue
		}
		vx, vy := visual(s)
		d := common.FacingToDelta(s.facing)
		wf := int(simTime*3) % resources.DirectionalWakeFrames
		blitAt(img, resources.GetDirectionalWake(s.facing, wf),
			(vx-float64(d.X)*0.9)*cell, (vy-float64(d.Y)*0.9)*cell)
	}

	for _, s := range ships {
		vx, vy := visual(s)
		roll := 0.045 * math.Sin(simTime*1.4+s.phase)
		bob := 1.6 * math.Sin(simTime*2.2+s.phase*1.7)
		hull := toNRGBA(resources.GetShipTile(s.ship, s.facing))
		blitAt(img, resources.RotateSprite(hull, roll), vx*cell, vy*cell+bob)
	}

	// Player chrome rides the first ship.
	px, py := visual(ships[0])
	blitAt(img, resources.GetPlayerMarkerOverlay(cell), px*cell, py*cell)
	blitAt(img, resources.GetBowMarkerOverlay(cell, ships[0].facing), px*cell, py*cell)

	f, err := os.Create(path)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		fatal(err)
	}
	fmt.Println("wrote", path)
}

func blitAt(dst *image.RGBA, src image.Image, x, y float64) {
	if src == nil {
		return
	}
	b := src.Bounds()
	r := image.Rect(int(x+0.5), int(y+0.5), int(x+0.5)+b.Dx(), int(y+0.5)+b.Dy())
	draw.Draw(dst, r, src, b.Min, draw.Over)
}

func toNRGBA(src image.Image) *image.NRGBA {
	dst := image.NewNRGBA(src.Bounds())
	draw.Draw(dst, dst.Bounds(), src, src.Bounds().Min, draw.Src)
	return dst
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
