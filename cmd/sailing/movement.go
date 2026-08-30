package sailing

import (
	"math/rand"
	"pirate-wars/cmd/common"
)

// ShouldMove rolls whether a ship advances one cell this tick given effective speed.
func ShouldMove(speed float64) bool {
	if speed <= 0 {
		return false
	}
	if speed >= 1 {
		return true
	}
	return rand.Float64() < speed
}

// WakeAftPosition returns the cell immediately behind the ship for wake overlay.
func WakeAftPosition(pos common.Coordinates, heading common.Facing) common.Coordinates {
	delta := common.FacingToDelta(heading)
	return common.Coordinates{X: pos.X - delta.X, Y: pos.Y - delta.Y}
}
