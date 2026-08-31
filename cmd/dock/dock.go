package dock

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/town"
)

type boatWorld interface {
	IsPassableByBoat(common.Coordinates) bool
}

// HarborDock reports whether pos is on green parking water in the harbor region.
type HarborDock interface {
	IsDock(common.Coordinates) bool
}

// CanDock reports whether the player can open the dock overlay at pos.
// Harbor: dock ON green mask pixels. Open sea: adjacent town on water (slice 2).
func CanDock(pos common.Coordinates, world boatWorld, towns *town.Towns, harbor HarborDock) bool {
	if harbor != nil && harbor.IsDock(pos) {
		return true
	}
	return towns != nil && towns.AdjacentTown(pos, world) != nil
}

// AdjacentTown returns the town for docking at pos.
// Harbor green uses the settlement town; otherwise adjacent procedural town.
func AdjacentTown(pos common.Coordinates, world boatWorld, towns *town.Towns, harbor HarborDock) *town.Town {
	if harbor != nil && harbor.IsDock(pos) && towns != nil {
		if t := towns.HarborSettlement(); t != nil {
			return t
		}
	}
	if towns == nil {
		return nil
	}
	return towns.AdjacentTown(pos, world)
}
