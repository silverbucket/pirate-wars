package common

const (
	LogFile       = "pirate-wars.log"
	WorldWidth    = 600 // Y
	WorldHeight   = 600 // X
	TotalTowns    = 30
	ViewWidth     = 75
	ViewHeight    = 50
	MiniMapFactor = 11
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
