package world

import (
	"image"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/resources"
)

func (world *MapView) AdvanceAnimation() {
	resources.AdvanceWaveAnimation()
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

func (world *MapView) terrainTileAt(pos common.Coordinates) image.Image {
	tt := world.terrain.Cells[pos.X][pos.Y]

	if tt == common.TerrainTypeShallowWater {
		return resources.GetWaveTile(resources.CurrentWaveFrame())
	}

	if resources.IsCoastLandTerrain(tt) {
		if coast := resources.GetCoastTile(world.waterNeighbors(pos)); coast != nil {
			return coast
		}
	}

	return resources.GetTerrainTile(tt)
}
