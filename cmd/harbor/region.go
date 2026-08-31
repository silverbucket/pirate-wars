package harbor

import "pirate-wars/cmd/common"

// Painted harbor region dimensions (pixels).
const (
	PixelWidth  = 1536
	PixelHeight = 1024
	MaskCell    = 64
	MaskCols    = PixelWidth / MaskCell  // 24
	MaskRows    = PixelHeight / MaskCell // 16
)

// Harbor occupies this many 32px world cells on the main 800×800 map.
const (
	WorldCols = PixelWidth / 32  // 48
	WorldRows = PixelHeight / 32 // 32
)

// Origin is the top-left world cell of the harbor rect on the main map.
var Origin = common.Coordinates{X: 200, Y: 280}

// TownPos is the settlement anchor in world cells (green parking water).
var TownPos = common.Coordinates{X: Origin.X + 26, Y: Origin.Y + 12}

// InRegion reports whether world cell c lies inside the harbor rect.
func InRegion(c common.Coordinates) bool {
	return c.X >= Origin.X && c.X < Origin.X+WorldCols &&
		c.Y >= Origin.Y && c.Y < Origin.Y+WorldRows
}

// WorldToLocal converts an in-region world cell to harbor-local cell coords.
func WorldToLocal(c common.Coordinates) (lx, ly int, ok bool) {
	if !InRegion(c) {
		return 0, 0, false
	}
	return c.X - Origin.X, c.Y - Origin.Y, true
}

// LocalToWorld converts harbor-local cell coords to world cells.
func LocalToWorld(lx, ly int) common.Coordinates {
	return common.Coordinates{X: Origin.X + lx, Y: Origin.Y + ly}
}

// CellCenterPixel returns the center of a world cell in harbor pixel space.
func CellCenterPixel(c common.Coordinates) (px, py int, ok bool) {
	lx, ly, ok := WorldToLocal(c)
	if !ok {
		return 0, 0, false
	}
	return lx*32 + 16, ly*32 + 16, true
}
