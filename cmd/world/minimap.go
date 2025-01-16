package world

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/screen"
	"pirate-wars/cmd/terrain"
)

func (mm *MiniMapView) SetPositionType(c common.Coordinates, tt terrain.TerrainType) {
	mm.grid[c.X][c.Y] = tt
}

func (world *MapView) GenerateMiniMap() {
	// Calculate MiniMap dimensions
	height := len(world.grid) / screen.MiniMapFactor
	width := len(world.grid[0]) / screen.MiniMapFactor

	// Create new 2D slice
	world.miniMap = MiniMapView{
		grid: make([][]terrain.TerrainType, height+1),
	}
	for i := range world.miniMap.grid {
		world.miniMap.grid[i] = make([]terrain.TerrainType, width+1)
	}

	// Down-sample
	for i, row := range world.grid {
		for j, val := range row {
			// Calculate corresponding index in new slice
			c := common.Coordinates{
				X: i / screen.MiniMapFactor,
				Y: j / screen.MiniMapFactor,
			}
			// Assign original TerrainType value
			world.miniMap.SetPositionType(c, val)
		}
	}
	for _, m := range world.mapItems {
		world.miniMap.SetPositionType(common.GetMiniMapScale(m.GetPos()), m.GetTerrainType())
	}
}
