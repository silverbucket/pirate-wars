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

var facingToDelta = []Coordinates{
	{0, -1},  // N
	{1, -1},  // NE
	{1, 0},   // E
	{1, 1},   // SE
	{0, 1},   // S
	{-1, 1},  // SW
	{-1, 0},  // W
	{-1, -1}, // NW
}

var facingLabels = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

// FacingToDelta returns the one-cell movement delta for a facing.
func FacingToDelta(f Facing) Coordinates {
	if int(f) < 0 || int(f) >= len(facingToDelta) {
		return Coordinates{}
	}
	return facingToDelta[f]
}

// FacingFromDirection sets heading from a keyboard/direction delta without moving.
func FacingFromDirection(d Coordinates) (Facing, bool) {
	return FacingFromDelta(d.X, d.Y)
}

// RotateFacing steps around the compass by delta steps (-1 or +1 typical).
func RotateFacing(f Facing, delta int) Facing {
	n := (int(f) + delta) % 8
	if n < 0 {
		n += 8
	}
	return Facing(n)
}

// OppositeFacing returns the 180° opposite direction.
func OppositeFacing(f Facing) Facing {
	return RotateFacing(f, 4)
}

// FacingSeparation returns the minimal steps (0–4) between two facings on the compass.
func FacingSeparation(a, b Facing) int {
	diff := int(a) - int(b)
	if diff < 0 {
		diff = -diff
	}
	if diff > 4 {
		diff = 8 - diff
	}
	return diff
}

// FacingLabel returns a short compass label.
func FacingLabel(f Facing) string {
	if int(f) < 0 || int(f) >= len(facingLabels) {
		return "?"
	}
	return facingLabels[f]
}
