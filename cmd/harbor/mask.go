package harbor

import (
	"image"
	"image/color"
	"pirate-wars/cmd/common"
)

// CellKind classifies mask pixels.
type CellKind int

const (
	KindSail  CellKind = iota // blue — open water / shallows / harbor mouth
	KindBlock                 // red — land, rocks, walls, buildings, pier planks
	KindDock                  // green — parking water; dock overlay when ship is here
)

// Mask samples the harbor collision / dock mask.
type Mask struct {
	img image.Image
}

// NewMask wraps a decoded 1536×1024 mask image.
func NewMask(img image.Image) *Mask {
	return &Mask{img: img}
}

// KindAtPixel classifies a harbor-local pixel.
func (m *Mask) KindAtPixel(px, py int) CellKind {
	if px < 0 || py < 0 || px >= PixelWidth || py >= PixelHeight {
		return KindBlock
	}
	r, g, b, _ := m.img.At(px, py).RGBA()
	r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
	return classifyRGB(r8, g8, b8)
}

// KindAtWorld samples the mask at the center of world cell c.
func (m *Mask) KindAtWorld(c common.Coordinates) CellKind {
	px, py, ok := CellCenterPixel(c)
	if !ok {
		return KindBlock
	}
	return m.KindAtPixel(px, py)
}

// IsPassable reports whether a boat may enter world cell c.
func (m *Mask) IsPassable(c common.Coordinates) bool {
	k := m.KindAtWorld(c)
	return k == KindSail || k == KindDock
}

// IsDock reports whether the player may open the dock overlay at c (on green).
func (m *Mask) IsDock(c common.Coordinates) bool {
	return m.KindAtWorld(c) == KindDock
}

// MaskCellKind returns the dominant kind for a 64×64 mask cell (for tests / debug).
func (m *Mask) MaskCellKind(col, row int) CellKind {
	if col < 0 || col >= MaskCols || row < 0 || row >= MaskRows {
		return KindBlock
	}
	cx := col*MaskCell + MaskCell/2
	cy := row*MaskCell + MaskCell/2
	return m.KindAtPixel(cx, cy)
}

func classifyRGB(r, g, b uint8) CellKind {
	// Official mask: blue sail, red stop, green parking water.
	if g > 100 && g > r+20 && g > b+20 {
		return KindDock
	}
	if r > 120 && r > g+30 && r > b+30 {
		return KindBlock
	}
	if b > 80 && b > r {
		return KindSail
	}
	// Fallback: treat unknown as blocked.
	return KindBlock
}

// RGBA helpers for tests.
func rgba(r, g, b uint8) color.Color {
	return color.RGBA{R: r, G: g, B: b, A: 255}
}
