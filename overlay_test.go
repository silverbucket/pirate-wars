package main

import (
	"strings"
	"testing"

	"pirate-wars/cmd/world"
)

// TestOverlayRowsAndButtonsStayInsidePanel covers a review finding: the panel
// height was clamped to the viewport but every row was still laid out, so a
// screen taller than the clamp drew text and left live tap targets below the
// panel, over the scrim and the action bar.
func TestOverlayRowsAndButtonsStayInsidePanel(t *testing.T) {
	ViewType = world.ViewTypeHail
	gs := mainMapGameState()

	rows := make([]overlayRow, 0, 200)
	for i := 0; i < 200; i++ {
		rows = append(rows, overlayRow{text: "a very long hail from a very talkative captain"})
	}
	rows = append(rows, overlayRow{buttons: []overlayAction{{label: "Dismiss", do: func() {}}}})
	screen := overlayScreen{title: "Hail", rows: rows}

	panel, buttons := gs.overlayLayout(screen)
	if !panel.In(overlayRegion()) {
		t.Fatalf("panel %v escapes the overlay region %v", panel, overlayRegion())
	}
	for _, b := range buttons {
		if !b.rect.In(panel) {
			t.Fatalf("button %q at %v escapes the panel %v", b.label, b.rect, panel)
		}
	}

	_, fitted := gs.overlayPanel(screen)
	if len(fitted) >= len(screen.rows) {
		t.Fatal("an oversized screen should drop the rows that do not fit")
	}
	if len(fitted) == 0 {
		t.Fatal("even a truncated screen should render something")
	}
}

// TestTruncatedOverlayKeepsItsWayOut checks a clipped screen still shows the
// dismiss button, so the player is never trapped behind a long overlay.
func TestTruncatedOverlayKeepsItsWayOut(t *testing.T) {
	ViewType = world.ViewTypeHail
	gs := mainMapGameState()

	rows := make([]overlayRow, 0, 200)
	for i := 0; i < 200; i++ {
		rows = append(rows, overlayRow{text: "filler"})
	}
	rows = append(rows, overlayRow{buttons: []overlayAction{{label: "Dismiss", do: func() {}}}})

	_, buttons := gs.overlayLayout(overlayScreen{title: "Hail", rows: rows})
	found := false
	for _, b := range buttons {
		if b.label == "Dismiss" {
			found = true
		}
	}
	if !found {
		t.Fatal("a truncated overlay must keep its dismiss button")
	}
}

// TestWrapTextCountsRunesNotBytes covers a review finding: len() measures bytes,
// so a single non-ASCII character wrapped the line early. Town and captain names
// are generated text.
func TestWrapTextCountsRunesNotBytes(t *testing.T) {
	// Ten 2-byte runes: 20 bytes, 10 characters.
	word := strings.Repeat("é", 10)
	lines := wrapText(word+" "+word, 25)
	if len(lines) != 1 {
		t.Fatalf("21 characters should fit a 25-character line, wrapped to %d lines: %q", len(lines), lines)
	}
}

// TestOverlayNeverCoversTheActionBar keeps the modal clear of the command row,
// whose keys stay live while an overlay is open.
func TestOverlayNeverCoversTheActionBar(t *testing.T) {
	if overlayRegion().Max.Y > actionBarRect.Min.Y {
		t.Fatalf("overlay region %v runs into the action bar %v", overlayRegion(), actionBarRect)
	}

	for _, view := range []int{world.ViewTypeHelp, world.ViewTypeQuitConfirm, world.ViewTypeHail} {
		ViewType = view
		gs := mainMapGameState()
		panel, _ := gs.overlayPanel(gs.currentOverlayScreen())
		if panel.Overlaps(actionBarRect) {
			t.Fatalf("view %d: overlay panel %v overlaps the action bar %v", view, panel, actionBarRect)
		}
	}
}

// TestMerchantScreenFitsTheOverlay keeps the real dock screens inside the panel,
// which is where the clipping bug would actually bite a player.
func TestMerchantScreenFitsTheOverlay(t *testing.T) {
	for _, view := range []int{world.ViewTypeHelp, world.ViewTypeQuitConfirm} {
		ViewType = view
		gs := mainMapGameState()
		gs.refreshActionBar()

		panel, buttons := gs.overlayLayout(gs.currentOverlayScreen())
		for _, b := range buttons {
			if !b.rect.In(panel) {
				t.Fatalf("view %d: button %q at %v escapes panel %v", view, b.label, b.rect, panel)
			}
		}
	}
}
