package world

import (
	"fmt"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/screen"
	"pirate-wars/cmd/terrain"
)

func (mm *MiniMapView) SetPositionType(c common.Coordinates, tt terrain.Type) {
	mm.grid[c.X][c.Y] = tt
}

func (world *MapView) GenerateMiniMap() {
	// Calculate MiniMap dimensions
	width := len(world.grid) / screen.MiniMapFactor
	height := len(world.grid[0]) / screen.MiniMapFactor
	world.logger.Info(fmt.Sprintf("generating mini-map with dimensions %v,%v", width, height))

	// Create new 2D slice
	miniMapGrid := make([][]terrain.Type, width+1)
	for i := range miniMapGrid {
		miniMapGrid[i] = make([]terrain.Type, height+1)
	}

	world.miniMap = MiniMapView{
		grid: miniMapGrid,
	}

	// Down-sample
	for i, row := range world.grid {
		for j, val := range row {
			// Calculate corresponding index in new slice
			c := common.Coordinates{
				X: i / screen.MiniMapFactor,
				Y: j / screen.MiniMapFactor,
			}
			// Assign original Type value
			world.miniMap.SetPositionType(c, val)
		}
	}
	for _, m := range world.mapItems {
		world.miniMap.SetPositionType(common.GetMiniMapScale(m.GetPos()), m.GetTerrainType())
	}
}
