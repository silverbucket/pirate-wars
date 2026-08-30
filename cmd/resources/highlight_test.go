package resources

import (
	"image"
	"testing"
)

func TestCompositeWithHighlightDrawsFrame(t *testing.T) {
	base := image.NewRGBA(image.Rect(0, 0, 32, 32))
	composited := CompositeWithHighlight(base, 32)
	if composited == nil {
		t.Fatal("CompositeWithHighlight should not return nil")
	}

	// Top-left corner should pick up the highlight frame.
	_, _, _, a := composited.At(0, 0).RGBA()
	if a == 0 {
		t.Fatal("highlight composite should draw a visible frame")
	}
}

func TestGetPlayerMarkerOverlayVisible(t *testing.T) {
	marker := GetPlayerMarkerOverlay(32)
	if marker == nil {
		t.Fatal("player marker should not be nil")
	}

	if countOpaquePixels(marker) == 0 {
		t.Fatal("player marker should draw visible pixels")
	}
	if !HasTransparentCenter(marker) {
		t.Fatal("player marker center should be transparent")
	}
}
