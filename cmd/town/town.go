package town

import (
	"pirate-wars/cmd/common"
)

type Type struct {
	pos common.Coordinates
}

var List []Type

func (t *Type) GetX() int {
	return t.pos.X
}
func (t *Type) GetMiniMapX() int {
	return t.pos.X / common.MiniMapFactor
}

func (t *Type) GetY() int {
	return t.pos.Y
}

func (t *Type) GetMiniMapY() int {
	return t.pos.Y / common.MiniMapFactor
}

func Create(coords common.Coordinates, c rune) Type {
	t := Type{pos: coords}
	List = append(List, t)
	return t
}
