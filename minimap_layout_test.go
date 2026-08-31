package main

import (
	"testing"

	"pirate-wars/cmd/window"
)

// TestMinimapAndLegendClearTheActionBar covers a review finding.
//
// The chart raster is 700px and the map area above the action bar is 700px, so
// at 1:1 the chart ran to y=714 against a bar starting at y=700 — and the legend
// drawn below it landed inside the bar entirely, painting over its border and
// caption. The chart is now fitted to the space that exists.
func TestMinimapAndLegendClearTheActionBar(t *testing.T) {
	gs := mainMapGameState()
	bounds := gs.minimapBounds()

	if bounds.Overlaps(actionBarRect) {
		t.Fatalf("minimap and legend %v overlap the action bar %v", bounds, actionBarRect)
	}
	if bounds.Max.Y > actionBarRect.Min.Y {
		t.Fatalf("minimap bottom %d is below the action bar top %d", bounds.Max.Y, actionBarRect.Min.Y)
	}
	if bounds.Min.Y < 0 || bounds.Min.X < 0 {
		t.Fatalf("minimap %v starts off screen", bounds)
	}
	if bounds.Max.X > sidePanelRect.Min.X {
		t.Fatalf("minimap right edge %d runs under the side panel at %d", bounds.Max.X, sidePanelRect.Min.X)
	}
}

// TestMinimapFillsTheSpaceItHas keeps the fit from quietly shrinking the chart:
// it should still use most of the height available to it.
func TestMinimapFillsTheSpaceItHas(t *testing.T) {
	gs := mainMapGameState()
	bounds := gs.minimapBounds()

	available := overlayRegion().Dy()
	if used := bounds.Dy(); used < available*9/10 {
		t.Fatalf("minimap uses only %d of %d available px", used, available)
	}
	if window.MiniMapArea.Width != window.MiniMapArea.Height {
		t.Fatal("the fit assumes a square chart")
	}
}
