package resources

import "testing"

func TestTilesetHeight(t *testing.T) {
	height := GetTilesetHeight()
	if height != TileSize*10 && height != TileSize*12 && height != TileSize*14 && height != TileSize*16 {
		t.Fatalf("tileset height = %d, want a 10, 12, 14, or 16 row sheet", height)
	}

	if height >= ExpandedTilesetHeight {
		if !HasExpandedTileset() {
			t.Fatal("12-row tileset should report expanded")
		}
	} else if HasExpandedTileset() {
		t.Fatal("10-row tileset should not report expanded")
	}

	if height >= SailingVisualsTilesetHeight {
		if !HasSailingVisualsTileset() {
			t.Fatal("14-row tileset should report sailing visuals")
		}
	} else if HasSailingVisualsTileset() {
		t.Fatal("short tileset should not report sailing visuals")
	}
}
