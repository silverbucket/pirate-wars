package player

import (
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/resources"
	"pirate-wars/cmd/world"
)

func Create(world *world.MapView) *entities.Avatar {
	p := entities.CreateAvatar(world.RandomPositionDeepWater(), resources.GetShipTile(common.ShipWhite), color.White)
	return &p
}
