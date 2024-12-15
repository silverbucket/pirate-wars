package common

type Coordinates struct {
	X int
	Y int
}

const (
	WorldWidth       = 600
	WorldHeight      = 600
	TotalTowns       = 30
	ViewWidth        = 75
	ViewHeight       = 50
	MiniMapFactor    = 11
	TypeDeepWater    = 0
	TypeOpenWater    = 1
	TypeShallowWater = 2
	TypeBeach        = 3
	TypeLowland      = 4
	TypeHighland     = 5
	TypeRock         = 6
	TypePeak         = 7
	TypeTown         = 8
)

type ViewPort struct {
	width   int
	height  int
	topLeft int
}
