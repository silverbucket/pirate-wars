package player

import (
	"image/color"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/harbor"
	"pirate-wars/cmd/world"
)

// SpawnWorld optionally provides a harbor spawn when the painted region is active.
type SpawnWorld interface {
	RandomPositionDeepWater() common.Coordinates
}

// Create places the player; prefers harbor parking water when harbor world is available.
func Create(world *world.MapView, harborWorld *harbor.World) *entities.Avatar {
	pos := world.RandomPositionDeepWater()
	if harborWorld != nil {
		if spawn, ok := harborWorld.ClampSpawn(harbor.TownPos); ok {
			pos = spawn
		}
	}
	p := entities.CreateAvatar(pos, common.ShipWhite, color.White)
	return &p
}
