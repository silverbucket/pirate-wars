package player

import (
	"image/color"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/world"
)

func Create(world *world.MapView) *entities.Avatar {
	p := entities.CreateAvatar(world.RandomPositionDeepWater(), '⏏', entities.ColorScheme{Foreground: color.Black, Background: color.RGBA{255, 125, 125, 105}})
	return &p
}
