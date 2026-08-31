package resources

import (
	"image"
	"pirate-wars/cmd/common"
)

// WaterNeighbors records orthogonal water adjacency for coast picking.
type WaterNeighbors struct {
	N, E, S, W bool
}

func IsWaterTerrain(tt common.TerrainType) bool {
	switch tt {
	case common.TerrainTypeDeepWater, common.TerrainTypeOpenWater, common.TerrainTypeShallowWater:
		return true
	default:
		return false
	}
}

// IsCoastLandTerrain returns land types that may show a coast edge tile.
func IsCoastLandTerrain(tt common.TerrainType) bool {
	switch tt {
	case common.TerrainTypeBeach, common.TerrainTypeLowland, common.TerrainTypeHighland, common.TerrainTypeRock:
		return true
	default:
		return false
	}
}

// PickCoastTile selects a coast tile for convex outside corners and edges.
// Corners take precedence over single edges.
func PickCoastTile(neighbors WaterNeighbors) (col, row int, ok bool) {
	if !HasExpandedTileset() {
		return 0, 0, false
	}

	switch {
	case neighbors.N && neighbors.E:
		return CoastNECol, CoastRow, true
	case neighbors.S && neighbors.E:
		return CoastSECol, CoastRow, true
	case neighbors.S && neighbors.W:
		return CoastSWCol, CoastSWRow, true
	case neighbors.N && neighbors.W:
		return CoastNWCol, CoastNWRow, true
	case neighbors.N:
		return CoastNorthCol, CoastRow, true
	case neighbors.E:
		return CoastEastCol, CoastRow, true
	case neighbors.S:
		return CoastSouthCol, CoastRow, true
	case neighbors.W:
		return CoastWestCol, CoastRow, true
	default:
		return 0, 0, false
	}
}

// GetCoastTile returns the coast tile image for the given neighbors, or nil.
func GetCoastTile(neighbors WaterNeighbors) image.Image {
	return GetCoastTileFrame(neighbors, 0)
}

// GetCoastTileFrame returns a frame of the animated shoreline: odd frames use
// the phase-shifted waterline rows, so the foam laps against the sand.
func GetCoastTileFrame(neighbors WaterNeighbors, frame int) image.Image {
	col, row, ok := PickCoastTile(neighbors)
	if !ok {
		return nil
	}
	if frame%2 != 0 && HasAnimatedCoastTileset() {
		row += CoastFrameRowStride
	}
	return extractTileAt(col, row)
}
