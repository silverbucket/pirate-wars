package harbor

import (
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/window"
)

// World adapts harbor mask passability for sailing within the harbor rect.
type World struct {
	mask *Mask
}

// NewWorld wraps a harbor mask for boat passability checks.
func NewWorld(mask *Mask) *World {
	return &World{mask: mask}
}

// InHarbor reports whether pos is inside the harbor region.
func (w *World) InHarbor(pos common.Coordinates) bool {
	return InRegion(pos)
}

// IsPassableByBoat uses the painted mask when inside the harbor rect.
// A nil receiver means the harbor art is not installed, so nothing is sailable
// here and callers fall back to the tilemap.
func (w *World) IsPassableByBoat(pos common.Coordinates) bool {
	if w == nil || w.mask == nil || !InRegion(pos) {
		return false
	}
	return w.mask.IsPassable(pos)
}

// IsDock reports dock overlay eligibility (green parking water).
func (w *World) IsDock(pos common.Coordinates) bool {
	if w == nil || w.mask == nil || !InRegion(pos) {
		return false
	}
	return w.mask.IsDock(pos)
}

// GetViewportRegion returns follow-cam region in world cells (for HUD/minimap compat).
func GetViewportRegion(pos common.Coordinates) window.Region {
	vp := window.Region{
		Cols: window.ViewPort.Region.Cols,
		Rows: window.ViewPort.Region.Rows,
		X:    pos.X - window.ViewPort.Region.Cols/2,
		Y:    pos.Y - window.ViewPort.Region.Rows/2,
	}
	if vp.X < Origin.X {
		vp.X = Origin.X
	} else if vp.X+vp.Cols > Origin.X+WorldCols {
		vp.X = Origin.X + WorldCols - vp.Cols
	}
	if vp.Y < Origin.Y {
		vp.Y = Origin.Y
	} else if vp.Y+vp.Rows > Origin.Y+WorldRows {
		vp.Y = Origin.Y + WorldRows - vp.Rows
	}
	return vp
}

// ClampSpawn returns a sailable cell inside the harbor, or false.
func (w *World) ClampSpawn(c common.Coordinates) (common.Coordinates, bool) {
	if w == nil || w.mask == nil {
		return common.Coordinates{}, false
	}
	if w.mask.IsPassable(c) {
		return c, true
	}
	// Spiral search from town anchor for blue/green water.
	lx, ly, ok := WorldToLocal(TownPos)
	if !ok {
		return common.Coordinates{}, false
	}
	for radius := 0; radius < 20; radius++ {
		for dy := -radius; dy <= radius; dy++ {
			for dx := -radius; dx <= radius; dx++ {
				if dx*dx+dy*dy > radius*radius {
					continue
				}
				wc := LocalToWorld(lx+dx, ly+dy)
				if w.mask.IsPassable(wc) {
					return wc, true
				}
			}
		}
	}
	return common.Coordinates{}, false
}
