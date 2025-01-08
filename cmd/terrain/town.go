package terrain

import (
	"pirate-wars/cmd/common"
)

type HeatMap = [][]int

type Town struct {
	pos     common.Coordinates
	heatMap HeatMap
}

func (t *Town) GetX() int {
	return t.pos.X
}
func (t *Town) GetMiniMapX() int {
	return t.pos.X / common.MiniMapFactor
}

func (t *Town) GetY() int {
	return t.pos.Y
}

func (t *Town) GetMiniMapY() int {
	return t.pos.Y / common.MiniMapFactor
}
