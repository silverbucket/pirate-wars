package resources

import (
	"image"
	"pirate-wars/cmd/common"
)

const (
	WaveFrameCount    = 3
	WaveTicksPerFrame = 3
)

var waveCoords = []image.Point{
	{X: 2, Y: 1},  // frame 1 — shallow water with wave line
	{X: 2, Y: 11}, // frame 2
	{X: 3, Y: 11}, // frame 3
}

var (
	waveFrameIndex  int
	waveTickCounter int
)

// AdvanceWaveAnimation advances the shallow-water wave cycle on game ticks.
func AdvanceWaveAnimation() {
	waveTickCounter++
	if waveTickCounter >= WaveTicksPerFrame {
		waveTickCounter = 0
		waveFrameIndex = (waveFrameIndex + 1) % WaveFrameCount
	}
}

// CurrentWaveFrame returns the active wave animation frame index.
func CurrentWaveFrame() int {
	return waveFrameIndex
}

// WaveFrameCoord returns the tileset coordinate for a wave frame index.
func WaveFrameCoord(frame int) (col, row int, ok bool) {
	if frame < 0 || frame >= len(waveCoords) {
		return 0, 0, false
	}
	coord := waveCoords[frame]
	return coord.X, coord.Y, true
}

// GetWaveTile returns the animated shallow-water tile for the given frame.
func GetWaveTile(frame int) image.Image {
	col, row, ok := WaveFrameCoord(frame)
	if !ok {
		col, row, _ = WaveFrameCoord(0)
	}

	if tileRegionInBounds(col, row) {
		if tile := extractTileAt(col, row); tile != nil {
			return tile
		}
	}

	return GetTerrainTile(common.TerrainTypeShallowWater)
}
