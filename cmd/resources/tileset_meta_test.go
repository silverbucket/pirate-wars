package resources

import "testing"

func TestTilesetHeight(t *testing.T) {
	height := GetTilesetHeight()
	if height != 320 && height != 384 {
		t.Fatalf("tileset height = %d, want 320 or 384", height)
	}

	if height == 384 {
		if !HasExpandedTileset() {
			t.Fatal("384px tileset should report expanded")
		}
	} else if HasExpandedTileset() {
		t.Fatal("320px tileset should not report expanded")
	}
}
