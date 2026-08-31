package world

import (
	"image"

	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/gfx"
	"pirate-wars/cmd/resources"
	"pirate-wars/cmd/window"

	"github.com/hajimehoshi/ebiten/v2"
)

// pennantOffset lifts the wind pennant off the hull to the top-left of the
// player's cell. Drawn centred it sits directly over the sail, in the same cream
// and red as the sail art, and hides the pixels that tell the eight facings apart.
const (
	pennantOffsetX = -10
	pennantOffsetY = -11
)

// Draw renders the 32px tilemap viewport centred on avatar into dst.
func (world *MapView) Draw(
	dst *ebiten.Image,
	avatar entities.AvatarReadOnly,
	npcs []entities.AvatarReadOnly,
	highlight entities.ViewableEntity,
	windFacing common.Facing,
) {
	p := avatar.GetPos()
	h := highlight.GetPos()
	vpr := window.GetViewportRegion(p)

	ships := make(map[int]entities.AvatarReadOnly, len(npcs)+1)
	ships[common.CoordToKey(p)] = avatar
	for _, n := range npcs {
		ships[common.CoordToKey(n.GetPos())] = n
	}

	wakeCells := make(map[int]bool, len(ships))
	for _, ship := range ships {
		if ship.MovedThisTick() {
			aft := resources.WakeAftPosition(ship.GetPos(), ship.GetFacing())
			wakeCells[common.CoordToKey(aft)] = true
		}
	}

	highlightVisible := false
	if h.X >= 0 {
		if !highlight.IsHighlighted() {
			highlight.Highlight(true)
		}
		_, _, _, a := highlight.GetColor().RGBA()
		highlightVisible = a > 0
	}

	for x := 0; x < vpr.Cols; x++ {
		for y := 0; y < vpr.Rows; y++ {
			pos := common.Coordinates{X: vpr.X + x, Y: vpr.Y + y}
			if pos.X < 0 || pos.X >= common.WorldCols || pos.Y < 0 || pos.Y >= common.WorldRows {
				continue
			}
			sx := float64(x * window.CellSize)
			sy := float64(y * window.CellSize)

			blit(dst, world.terrainTileAt(pos), sx, sy)

			key := common.CoordToKey(pos)
			if ship, ok := ships[key]; ok {
				if !ship.IsHighlighted() || highlightVisible {
					blit(dst, ship.GetTileImage(), sx, sy)
				}
			}

			if h.X >= 0 && common.CoordsMatch(pos, h) && highlight.IsHighlighted() && highlightVisible {
				blit(dst, resources.GetExamineRingOverlay(window.CellSize), sx, sy)
			}
			if wakeCells[key] {
				blit(dst, resources.GetWakeOverlay(resources.CurrentWakeFrame()), sx, sy)
			}
			if common.CoordsMatch(pos, p) {
				blit(dst, resources.GetPlayerMarkerOverlay(window.CellSize), sx, sy)
				// The bow chevron is the player's heading readout on the map. It
				// goes on last so nothing can cover it.
				blit(dst, resources.GetBowMarkerOverlay(window.CellSize, avatar.GetFacing()), sx, sy)
				// Wind is map-wide, so one pennant on the player's own masthead
				// says everything a pennant on all 150 NPCs would. Offsetting it
				// clear of the hull keeps it from reading as part of the sail and
				// leaves the 8-way art unobscured.
				blit(dst, resources.GetPennantOverlay(windFacing), sx+pennantOffsetX, sy+pennantOffsetY)
			}
		}
	}
}

// blit draws a cached tile image at (x, y) through the Ebiten texture cache.
func blit(dst *ebiten.Image, src image.Image, x, y float64) {
	tex := gfx.Texture(src)
	if tex == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	dst.DrawImage(tex, op)
}
