package player

import (
	"pirate-wars/cmd/terrain"
)

func Create(t *terrain.Terrain) terrain.Avatar {
	p := terrain.CreateAvatar(t.RandomPositionDeepWater(), '⏏', terrain.ColorScheme{"0", "#ffffff"})
	return p
}
