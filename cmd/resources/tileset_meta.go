package resources

const ExpandedTilesetHeight = 384

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

	ExamineRingCol   = 4
	ExamineRingRow   = 11
	PlayerMarkerCol  = 5
	PlayerMarkerRow  = 11
)

// GetTilesetHeight returns the loaded tileset height in pixels.
func GetTilesetHeight() int {
	return getTileset().Bounds().Dy()
}

// HasExpandedTileset reports whether rows 10–11 are available.
func HasExpandedTileset() bool {
	return GetTilesetHeight() >= ExpandedTilesetHeight
}
