package player

import (
	"pirate-wars/cmd/npc"
	"pirate-wars/cmd/terrain"
)

func Create(world terrain.MapView) npc.Avatar {
	p := npc.CreateAvatar(world.RandomPositionDeepWater(), '⏏', npc.ColorScheme{"0", "#ffffff"})
	return p
}
