package window

import "testing"

func TestViewportTileCountsAtCellSize32(t *testing.T) {
	if CellSize != 32 {
		t.Fatalf("CellSize = %d, want 32", CellSize)
	}

	expectedCols := viewPortWidth / CellSize
	expectedRows := viewPortHeight / CellSize
	if ViewPort.Region.Cols != expectedCols {
		t.Fatalf("viewport cols = %d, want %d", ViewPort.Region.Cols, expectedCols)
	}
	if ViewPort.Region.Rows != expectedRows {
		t.Fatalf("viewport rows = %d, want %d", ViewPort.Region.Rows, expectedRows)
	}

	// Roughly 24x18 tiles visible after sidebar and action menu.
	if ViewPort.Region.Cols < 22 || ViewPort.Region.Cols > 28 {
		t.Fatalf("viewport cols = %d, want roughly 24", ViewPort.Region.Cols)
	}
	if ViewPort.Region.Rows < 16 || ViewPort.Region.Rows > 24 {
		t.Fatalf("viewport rows = %d, want roughly 18", ViewPort.Region.Rows)
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
