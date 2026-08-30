package dock

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/town"
)

type boatWorld interface {
	IsPassableByBoat(common.Coordinates) bool
}

// CanDock reports whether the player can open the dock overlay at pos.
func CanDock(pos common.Coordinates, world boatWorld, towns *town.Towns) bool {
	return towns != nil && towns.AdjacentTown(pos, world) != nil
}

// AdjacentTown returns the town adjacent to pos on water, if any.
func AdjacentTown(pos common.Coordinates, world boatWorld, towns *town.Towns) *town.Town {
	if towns == nil {
		return nil
	}
	return towns.AdjacentTown(pos, world)
}
