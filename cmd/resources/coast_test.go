package resources

import "testing"

func TestPickCoastTileEdges(t *testing.T) {
	cases := []struct {
		name      string
		neighbors WaterNeighbors
		wantCol   int
		wantRow   int
	}{
		{"north", WaterNeighbors{N: true}, CoastNorthCol, CoastRow},
		{"east", WaterNeighbors{E: true}, CoastEastCol, CoastRow},
		{"south", WaterNeighbors{S: true}, CoastSouthCol, CoastRow},
		{"west", WaterNeighbors{W: true}, CoastWestCol, CoastRow},
		{"northeast", WaterNeighbors{N: true, E: true}, CoastNECol, CoastRow},
		{"southeast", WaterNeighbors{S: true, E: true}, CoastSECol, CoastRow},
		{"southwest", WaterNeighbors{S: true, W: true}, CoastSWCol, CoastSWRow},
		{"northwest", WaterNeighbors{N: true, W: true}, CoastNWCol, CoastNWRow},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			col, row, ok := PickCoastTile(tc.neighbors)
			if !ok && HasExpandedTileset() {
				t.Fatalf("expected coast tile for %s", tc.name)
			}
			if !HasExpandedTileset() {
				if ok {
					t.Fatalf("coast tile should not be available on 10-row sheet")
				}
				return
			}
			if col != tc.wantCol || row != tc.wantRow {
				t.Fatalf("coast tile = (%d,%d), want (%d,%d)", col, row, tc.wantCol, tc.wantRow)
			}
		})
	}
}

func TestPickCoastTileCornerPriority(t *testing.T) {
	if !HasExpandedTileset() {
		t.Skip("expanded tileset not present")
	}

	col, row, ok := PickCoastTile(WaterNeighbors{N: true, E: true, S: true})
	if !ok {
		t.Fatal("expected coast tile")
	}
	if col != CoastNECol || row != CoastRow {
		t.Fatalf("corner priority = (%d,%d), want NE (%d,%d)", col, row, CoastNECol, CoastRow)
	}
}

func TestPickCoastTileNoWater(t *testing.T) {
	_, _, ok := PickCoastTile(WaterNeighbors{})
	if ok {
		t.Fatal("expected no coast tile without water neighbors")
	}
}
