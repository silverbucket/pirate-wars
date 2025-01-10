package common

const (
	LogFile         = "pirate-wars.log"
	WorldWidth      = 600 // Y
	WorldHeight     = 600 // X
	TotalTowns      = 20
	ViewWidth       = 75
	ViewHeight      = 50
	MiniMapFactor   = 11
	TotalNpcs       = 50
	MaxMovementCost = 9999
)

type ViewPort struct {
	width   int
	height  int
	topLeft int
}

type Coordinates struct {
	X int
	Y int
}

// Directions to explore (up, down, left, right)
var Directions = []Coordinates{
	{-1, 0},  // up
	{-1, -1}, // up left
	{-1, 1},  // up right
	{1, 0},   // down
	{1, -1},  // down left
	{1, 1},   // down right
	{0, -1},  // left
	{0, 1},   // right
}
