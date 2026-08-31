package resources

import (
	"testing"

	"pirate-wars/cmd/common"
)

// TestBowMarkerIsDistinctPerFacing is the measurement behind the marker.
//
// The 8-way hull sprites are near mirror images — east and west share over four
// fifths of their pixels — so with relative steering the art alone cannot tell
// the player which way the ship is pointed. Every marker must differ from every
// other, or it adds nothing.
func TestBowMarkerIsDistinctPerFacing(t *testing.T) {
	const size = 32
	seen := map[common.Facing]map[int]bool{}

	for f := common.FacingN; f <= common.FacingNW; f++ {
		img := GetBowMarkerOverlay(size, f)
		if img == nil {
			t.Fatalf("no bow marker for %s", common.FacingLabel(f))
		}
		pixels := map[int]bool{}
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				if _, _, _, a := img.At(x, y).RGBA(); a > 0x8000 {
					pixels[y*size+x] = true
				}
			}
		}
		if len(pixels) == 0 {
			t.Fatalf("bow marker for %s is empty", common.FacingLabel(f))
		}
		seen[f] = pixels
	}

	for a := common.FacingN; a <= common.FacingNW; a++ {
		for b := a + 1; b <= common.FacingNW; b++ {
			overlap := 0
			for p := range seen[a] {
				if seen[b][p] {
					overlap++
				}
			}
			smaller := len(seen[a])
			if len(seen[b]) < smaller {
				smaller = len(seen[b])
			}
			if overlap*2 > smaller {
				t.Fatalf("bow markers for %s and %s overlap on %d of %d pixels",
					common.FacingLabel(a), common.FacingLabel(b), overlap, smaller)
			}
		}
	}
}

// TestBowMarkerSitsAtTheBow checks the chevron leads the ship rather than
// sitting on the hull, where it would cover the art it is disambiguating.
func TestBowMarkerSitsAtTheBow(t *testing.T) {
	const size = 32
	for f := common.FacingN; f <= common.FacingNW; f++ {
		img := GetBowMarkerOverlay(size, f)
		d := common.FacingToDelta(f)

		sumX, sumY, n := 0, 0, 0
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				if _, _, _, a := img.At(x, y).RGBA(); a > 0x8000 {
					sumX += x
					sumY += y
					n++
				}
			}
		}
		cx, cy := float64(sumX)/float64(n)-float64(size)/2, float64(sumY)/float64(n)-float64(size)/2
		if float64(d.X)*cx+float64(d.Y)*cy <= 0 {
			t.Fatalf("bow marker for %s sits at (%.1f, %.1f), not ahead of centre",
				common.FacingLabel(f), cx, cy)
		}
	}
}

// TestBowMarkerIsCached keeps texture identity stable; the Ebiten texture cache
// keys on the image pointer, so a fresh image per frame would leak textures.
func TestBowMarkerIsCached(t *testing.T) {
	a := GetBowMarkerOverlay(32, common.FacingNE)
	b := GetBowMarkerOverlay(32, common.FacingNE)
	if a != b {
		t.Fatal("bow markers must be cached so the texture cache stays stable")
	}
}
