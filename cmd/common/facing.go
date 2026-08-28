package common

// Facing is the 8-way direction a ship sprite faces.
type Facing int

const (
	FacingN Facing = iota
	FacingNE
	FacingE
	FacingSE
	FacingS
	FacingSW
	FacingW
	FacingNW
)

var deltaToFacing = map[Coordinates]Facing{
	{0, -1}: FacingN,
	{1, -1}: FacingNE,
	{1, 0}:  FacingE,
	{1, 1}:  FacingSE,
	{0, 1}:  FacingS,
	{-1, 1}: FacingSW,
	{-1, 0}: FacingW,
	{-1, -1}: FacingNW,
}

// FacingFromDelta maps a movement delta to a facing direction.
// Returns the facing and true when the delta is one of the 8 compass directions.
func FacingFromDelta(dx, dy int) (Facing, bool) {
	facing, ok := deltaToFacing[Coordinates{dx, dy}]
	return facing, ok
}
