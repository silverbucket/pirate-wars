package window

import "testing"

func TestViewportTileCountsAtCellSize64(t *testing.T) {
	if CellSize != 64 {
		t.Fatalf("CellSize = %d, want 64", CellSize)
	}

	expectedCols := viewPortWidth / CellSize
	expectedRows := viewPortHeight / CellSize
	if ViewPort.Region.Cols != expectedCols {
		t.Fatalf("viewport cols = %d, want %d", ViewPort.Region.Cols, expectedCols)
	}
	if ViewPort.Region.Rows != expectedRows {
		t.Fatalf("viewport rows = %d, want %d", ViewPort.Region.Rows, expectedRows)
	}

	// Roughly 17x11 tiles visible after sidebar and action menu.
	if ViewPort.Region.Cols < 15 || ViewPort.Region.Cols > 19 {
		t.Fatalf("viewport cols = %d, want roughly 17", ViewPort.Region.Cols)
	}
	if ViewPort.Region.Rows < 9 || ViewPort.Region.Rows > 13 {
		t.Fatalf("viewport rows = %d, want roughly 11", ViewPort.Region.Rows)
	}
}

func TestViewportDimensionsDerivedFromLayout(t *testing.T) {
	if ViewPort.Dimensions.Width != viewPortWidth {
		t.Fatalf("viewport width = %d, want %d", ViewPort.Dimensions.Width, viewPortWidth)
	}
	if ViewPort.Dimensions.Height != viewPortHeight {
		t.Fatalf("viewport height = %d, want %d", ViewPort.Dimensions.Height, viewPortHeight)
	}
}
