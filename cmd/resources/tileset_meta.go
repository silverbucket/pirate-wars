package resources

const ExpandedTilesetHeight = TileSize * 12
const SailingVisualsTilesetHeight = TileSize * 14

// Coast tile coordinates (col, row).
const (
	CoastRow = 10

	CoastNorthCol = 0
	CoastEastCol  = 1
	CoastSouthCol = 2
	CoastWestCol  = 3
	CoastNECol    = 4
	CoastSECol    = 5
	CoastSWCol    = 0
	CoastNWCol    = 1
	CoastSWRow    = 11
	CoastNWRow    = 11

	WaveFrame2Col = 2
	WaveFrame2Row = 11
	WaveFrame3Col = 3
	WaveFrame3Row = 11

	ExamineRingCol  = 4
	ExamineRingRow  = 11
	PlayerMarkerCol = 5
	PlayerMarkerRow = 11
)

// GetTilesetHeight returns the loaded tileset height in pixels.
func GetTilesetHeight() int {
	return getTileset().Bounds().Dy()
}

// HasExpandedTileset reports whether rows 10–11 are available.
func HasExpandedTileset() bool {
	return GetTilesetHeight() >= ExpandedTilesetHeight
}

// HasSailingVisualsTileset reports whether rows 12–13 (pennant/wake) are available.
func HasSailingVisualsTileset() bool {
	return GetTilesetHeight() >= SailingVisualsTilesetHeight
}

// AnimatedCoastTilesetHeight marks a sheet carrying rows 14–15: the second
// coast-foam frame, offset CoastFrameRowStride rows below the first.
const AnimatedCoastTilesetHeight = TileSize * 16

// CoastFrameRowStride is the row offset from a coast tile to its frame-B twin.
const CoastFrameRowStride = 4

// HasAnimatedCoastTileset reports whether the breathing-shoreline frames exist.
func HasAnimatedCoastTileset() bool {
	return GetTilesetHeight() >= AnimatedCoastTilesetHeight
}

// Pennant tile coordinates (col, row) — pennant points downwind.
const (
	PennantRow = 12

	PennantNCol  = 0
	PennantNECol = 1
	PennantECol  = 2
	PennantSECol = 3
	PennantSCol  = 4
	PennantSWCol = 5

	PennantWRow  = 13
	PennantWCol  = 0
	PennantNWCol = 1
	PennantNWRow = 13

	WakeFrame0Col = 2
	WakeFrame0Row = 13
	WakeFrame1Col = 3
	WakeFrame1Row = 13
)
