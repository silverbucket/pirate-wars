package player

import (
	"image/color"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/world"
)

func Create(world *world.MapView) *entities.Avatar {
	p := entities.CreateAvatar(world.RandomPositionDeepWater(), '⏏', entities.ColorScheme{color.Black, color.White})
	return &p
}
