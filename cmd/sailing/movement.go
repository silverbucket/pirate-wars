package sailing

import (
	"pirate-wars/cmd/common"
)

// AccumulateSpeed adds tick speed to progress; returns true when a step is due.
func AccumulateSpeed(progress *float64, tickSpeed float64) bool {
	if tickSpeed <= 0 {
		return false
	}
	*progress += tickSpeed
	if *progress >= 1.0 {
		*progress -= 1.0
		return true
	}
	return false
}

// WakeAftPosition returns the cell immediately behind the ship for wake overlay.
func WakeAftPosition(pos common.Coordinates, heading common.Facing) common.Coordinates {
	delta := common.FacingToDelta(heading)
	return common.Coordinates{X: pos.X - delta.X, Y: pos.Y - delta.Y}
}
