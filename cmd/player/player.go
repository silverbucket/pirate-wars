package player

import (
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/world"
)

func Create(world *world.MapView) *npc.Avatar {
	p := npc.CreateAvatar(world.RandomPositionDeepWater(), '⏏', npc.ColorScheme{"0", "#ffffff"})
	return &p
}
