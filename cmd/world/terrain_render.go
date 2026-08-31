package world

import (
	"image"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/resources"
)

func (world *MapView) AdvanceAnimation() {
	resources.AdvanceWaveAnimation()
	resources.AdvanceWakeAnimation()
}

func (world *MapView) isWaterAt(x, y int) bool {
	if x < 0 || y < 0 || x >= world.GetWidth() || y >= world.GetHeight() {
		return false
	}
	return resources.IsWaterTerrain(world.GetPositionType(common.Coordinates{X: x, Y: y}))
}

func (world *MapView) waterNeighbors(c common.Coordinates) resources.WaterNeighbors {
	return resources.WaterNeighbors{
		N: world.isWaterAt(c.X, c.Y-1),
		E: world.isWaterAt(c.X+1, c.Y),
		S: world.isWaterAt(c.X, c.Y+1),
		W: world.isWaterAt(c.X-1, c.Y),
	}
}

// typeAt returns the terrain type at (x, y), or fallback out of bounds.
func (world *MapView) typeAt(x, y int, fallback common.TerrainType) common.TerrainType {
	if x < 0 || y < 0 || x >= world.GetWidth() || y >= world.GetHeight() {
		return fallback
	}
	return world.terrain.Cells[x][y]
}

func (world *MapView) terrainTileAt(pos common.Coordinates) image.Image {
	tt := world.terrain.Cells[pos.X][pos.Y]

	if resources.IsWaterTerrain(tt) {
		// Edges shared with a different water depth are dithered toward it,
		// so the depth bands read as one shelving sea, not stacked rectangles.
		ctx := resources.WaterContext{
			N: world.typeAt(pos.X, pos.Y-1, tt),
			E: world.typeAt(pos.X+1, pos.Y, tt),
			S: world.typeAt(pos.X, pos.Y+1, tt),
			W: world.typeAt(pos.X-1, pos.Y, tt),
		}
		return resources.GetBlendedWaterTile(tt, resources.CurrentWaveFrame(), ctx)
	}

	if resources.IsCoastLandTerrain(tt) {
		coast := resources.GetCoastTileFrame(world.waterNeighbors(pos), resources.CurrentWaveFrame())
		if coast != nil {
			return coast
		}
	}

	return resources.GetTerrainTile(tt)
}
