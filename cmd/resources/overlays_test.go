package resources

import (
	"image"
	"testing"
)

func TestOverlayTransparentCenter(t *testing.T) {
	overlays := []struct {
		name string
		img  func(int) image.Image
	}{
		{"player marker", GetPlayerMarkerOverlay},
		{"examine ring", GetExamineRingOverlay},
	}

	for _, tc := range overlays {
		t.Run(tc.name, func(t *testing.T) {
			img := tc.img(32)
			if img == nil {
				t.Fatal("overlay should not be nil")
			}
			if !HasTransparentCenter(img) {
				t.Fatal("overlay center should be transparent")
			}
		})
	}
}

func TestExpandedOverlayTilesUseAlpha(t *testing.T) {
	if !HasExpandedTileset() {
		t.Skip("expanded tileset not present")
	}

	for name, getter := range map[string]func(int) image.Image{
		"player marker":  GetPlayerMarkerOverlay,
		"examine ring":   GetExamineRingOverlay,
	} {
		t.Run(name, func(t *testing.T) {
			img := getter(32)
			if !HasTransparentCenter(img) {
				t.Fatal("tile overlay center should be transparent")
			}
			if countOpaquePixels(img) == 0 {
				t.Fatal("tile overlay should have opaque pixels")
			}
		})
	}
}

func countOpaquePixels(img image.Image) int {
	count := 0
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0x8000 {
				count++
			}
		}
	}
	return count
}
