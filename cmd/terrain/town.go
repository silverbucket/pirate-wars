package terrain

import (
	"pirate-wars/cmd/common"
)

type Town struct {
	pos common.Coordinates
}

var Towns []Town

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

func CreateTown(coords common.Coordinates, c rune) Town {
	t := Town{pos: coords}
	Towns = append(Towns, t)
	return t
}
