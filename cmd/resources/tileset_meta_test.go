package resources

import "testing"

func TestTilesetHeight(t *testing.T) {
	height := GetTilesetHeight()
	if height != 320 && height != 384 && height != 448 {
		t.Fatalf("tileset height = %d, want 320, 384, or 448", height)
	}

	if height >= ExpandedTilesetHeight {
		if !HasExpandedTileset() {
			t.Fatal("384px+ tileset should report expanded")
		}
	} else if HasExpandedTileset() {
		t.Fatal("320px tileset should not report expanded")
	}

	if height >= SailingVisualsTilesetHeight {
		if !HasSailingVisualsTileset() {
			t.Fatal("448px tileset should report sailing visuals")
		}
	} else if HasSailingVisualsTileset() {
		t.Fatal("sub-448 tileset should not report sailing visuals")
	}
}
