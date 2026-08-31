package resources

import (
	"image"
	"pirate-wars/cmd/common"
	"testing"
)

func TestSailingVisualsGating(t *testing.T) {
	height := GetTilesetHeight()
	if height >= SailingVisualsTilesetHeight {
		if !HasSailingVisualsTileset() {
			t.Fatal("448px tileset should enable sailing visuals")
		}
	} else {
		if HasSailingVisualsTileset() {
			t.Fatal("sub-448 tileset should not enable sailing visuals")
		}
	}
}

func TestPennantTileIndexGating(t *testing.T) {
	col, row := PennantTileIndex(common.FacingN)
	if HasSailingVisualsTileset() {
		if col != PennantNCol || row != PennantRow {
			t.Fatalf("pennant N = (%d,%d), want (%d,%d)", col, row, PennantNCol, PennantRow)
		}
	} else if col >= 0 || row >= 0 {
		t.Fatal("pennant index should be invalid without sailing visuals tileset")
	}
}

func TestPennantIndexFromFacing(t *testing.T) {
	if !HasSailingVisualsTileset() {
		t.Skip("sailing visuals tileset not present")
	}
	cases := []struct {
		facing           common.Facing
		wantCol, wantRow int
	}{
		{common.FacingN, PennantNCol, PennantRow},
		{common.FacingNE, PennantNECol, PennantRow},
		{common.FacingE, PennantECol, PennantRow},
		{common.FacingSE, PennantSECol, PennantRow},
		{common.FacingS, PennantSCol, PennantRow},
		{common.FacingSW, PennantSWCol, PennantRow},
		{common.FacingW, PennantWCol, PennantWRow},
		{common.FacingNW, PennantNWCol, PennantNWRow},
	}
	for _, tc := range cases {
		col, row := PennantTileIndex(tc.facing)
		if col != tc.wantCol || row != tc.wantRow {
			t.Fatalf("facing %v pennant = (%d,%d), want (%d,%d)", tc.facing, col, row, tc.wantCol, tc.wantRow)
		}
	}
}

func TestWakeTileIndexGating(t *testing.T) {
	col, row := WakeTileIndex(0)
	if HasSailingVisualsTileset() {
		if col != WakeFrame0Col || row != WakeFrame0Row {
			t.Fatalf("wake frame 0 = (%d,%d), want (%d,%d)", col, row, WakeFrame0Col, WakeFrame0Row)
		}
		col1, row1 := WakeTileIndex(1)
		if col1 != WakeFrame1Col || row1 != WakeFrame1Row {
			t.Fatalf("wake frame 1 = (%d,%d), want (%d,%d)", col1, row1, WakeFrame1Col, WakeFrame1Row)
		}
	} else if col >= 0 || row >= 0 {
		t.Fatal("wake index should be invalid without sailing visuals tileset")
	}
}

func TestWakeAftCoordinate(t *testing.T) {
	pos := common.Coordinates{X: 3, Y: 4}
	aft := WakeAftPosition(pos, common.FacingS)
	want := common.Coordinates{X: 3, Y: 3}
	if aft != want {
		t.Fatalf("wake aft = %v, want %v", aft, want)
	}
}

func TestWakeFrame1Extraction(t *testing.T) {
	if !HasSailingVisualsTileset() {
		t.Skip("sailing visuals tileset not present")
	}
	if GetWakeOverlay(1) == nil {
		t.Fatal("GetWakeOverlay(1) returned nil")
	}
}

func TestSailingOverlayTransparentWhenPresent(t *testing.T) {
	if !HasSailingVisualsTileset() {
		t.Skip("sailing visuals tileset not present")
	}
	pennant := GetPennantOverlay(common.FacingN)
	if pennant == nil {
		t.Fatal("pennant overlay should not be nil with 448px tileset")
	}
	if countOpaquePixels(pennant) == 0 {
		t.Fatal("pennant overlay should have opaque pixels")
	}
	if countTransparentPixels(pennant) == 0 {
		t.Fatal("pennant overlay should have transparent pixels")
	}
	wake0 := GetWakeOverlay(0)
	wake1 := GetWakeOverlay(1)
	if wake0 == nil && wake1 == nil {
		t.Fatal("at least one wake overlay frame should be present with 448px tileset")
	}
	if wake1 == nil {
		t.Fatal("wake frame 1 overlay should not be nil with 448px tileset")
	}
	if wake1 != nil {
		if countOpaquePixels(wake1) == 0 {
			t.Fatal("wake frame 1 should have opaque pixels")
		}
		if countTransparentPixels(wake1) == 0 {
			t.Fatal("wake frame 1 should have transparent pixels")
		}
	}
}

func countTransparentPixels(img image.Image) int {
	count := 0
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a <= 0x8000 {
				count++
			}
		}
	}
	return count
}
