package sailing

import "pirate-wars/cmd/common"

// Occupancy tracks one ship per cell for soft-bump collision.
type Occupancy map[int]string

// NewOccupancy builds a position map from ship id → position pairs.
func NewOccupancy(ships map[string]common.Coordinates) Occupancy {
	o := make(Occupancy, len(ships))
	for id, pos := range ships {
		o[common.CoordToKey(pos)] = id
	}
	return o
}

// OccupantAt returns the ship id occupying pos, excluding selfID.
func (o Occupancy) OccupantAt(pos common.Coordinates, selfID string) (string, bool) {
	id, ok := o[common.CoordToKey(pos)]
	if !ok || id == selfID {
		return "", false
	}
	return id, true
}

func (o Occupancy) IsOccupied(pos common.Coordinates, selfID string) bool {
	id, ok := o[common.CoordToKey(pos)]
	return ok && id != selfID
}

// TryStep attempts a one-cell move; returns the final position and whether movement occurred.
// Soft bump: if the target cell is occupied, the ship stays put.
func TryStep(
	from common.Coordinates,
	heading common.Facing,
	selfID string,
	o Occupancy,
	passable func(common.Coordinates) bool,
) (common.Coordinates, bool) {
	delta := common.FacingToDelta(heading)
	target := common.Coordinates{X: from.X + delta.X, Y: from.Y + delta.Y}
	if !common.Inbounds(target) {
		return from, false
	}
	if passable != nil && !passable(target) {
		return from, false
	}
	if o.IsOccupied(target, selfID) {
		return from, false
	}
	return target, true
}
