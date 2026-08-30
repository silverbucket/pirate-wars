package resources

import (
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
		facing       common.Facing
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

func TestSailingOverlayTransparentWhenPresent(t *testing.T) {
	if !HasSailingVisualsTileset() {
		t.Skip("sailing visuals tileset not present")
	}
	pennant := GetPennantOverlay(common.FacingN)
	if pennant == nil {
		t.Fatal("pennant overlay should not be nil with 448px tileset")
	}
	if !HasTransparentCenter(pennant) {
		t.Fatal("pennant overlay center should be transparent")
	}
	wake := GetWakeOverlay(0)
	if wake == nil {
		t.Fatal("wake overlay should not be nil with 448px tileset")
	}
	if !HasTransparentCenter(wake) {
		t.Fatal("wake overlay center should be transparent")
	}
}
