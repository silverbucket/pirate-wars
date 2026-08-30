package resources

import "testing"

func TestWaveFrameCoordsInRange(t *testing.T) {
	for frame := 0; frame < WaveFrameCount; frame++ {
		col, row, ok := WaveFrameCoord(frame)
		if !ok {
			t.Fatalf("wave frame %d missing coordinate", frame)
		}
		if col < 0 || row < 0 {
			t.Fatalf("wave frame %d has negative coordinate (%d,%d)", frame, col, row)
		}
		if frame == 0 {
			if !tileRegionInBounds(col, row) {
				t.Fatalf("wave frame 1 (%d,%d) should always be in bounds", col, row)
			}
			continue
		}
		if HasExpandedTileset() && !tileRegionInBounds(col, row) {
			t.Fatalf("wave frame %d (%d,%d) out of bounds on expanded sheet", frame, col, row)
		}
	}
}

func TestGetWaveTileFallsBackOnShortSheet(t *testing.T) {
	if HasExpandedTileset() {
		t.Skip("short-sheet fallback only applies to 10-row tileset")
	}

	tile := GetWaveTile(2)
	if tile == nil {
		t.Fatal("wave tile should not be nil")
	}
	if isTileNearlyEmpty(tile) {
		t.Fatal("wave tile fallback should not be empty")
	}
}

func TestWaveAnimationCycles(t *testing.T) {
	waveFrameIndex = 0
	waveTickCounter = 0

	for i := 0; i < WaveTicksPerFrame; i++ {
		AdvanceWaveAnimation()
	}
	if CurrentWaveFrame() != 1 {
		t.Fatalf("wave frame = %d, want 1 after one cycle", CurrentWaveFrame())
	}
}
