package resources

import (
	"image"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/sailing"
)

var pennantCols = []int{
	PennantNCol, PennantNECol, PennantECol, PennantSECol, PennantSCol, PennantSWCol, PennantWCol, PennantNWCol,
}
var pennantRows = []int{
	PennantRow, PennantRow, PennantRow, PennantRow, PennantRow, PennantRow, PennantWRow, PennantNWRow,
}

// PennantTileIndex returns the tileset column/row for a downwind-facing pennant.
func PennantTileIndex(facing common.Facing) (col, row int) {
	if !HasSailingVisualsTileset() {
		return -1, -1
	}
	i := int(facing)
	if i < 0 || i >= len(pennantCols) {
		return -1, -1
	}
	return pennantCols[i], pennantRows[i]
}

// GetPennantOverlay returns the wind pennant tile for the given downwind facing.
func GetPennantOverlay(facing common.Facing) image.Image {
	col, row := PennantTileIndex(facing)
	if col < 0 {
		return nil
	}
	return getSailingOverlayTile(col, row)
}

// WakeTileIndex returns column/row for wake animation frame 0 or 1.
func WakeTileIndex(frame int) (col, row int) {
	if !HasSailingVisualsTileset() {
		return -1, -1
	}
	if frame&1 == 0 {
		return WakeFrame0Col, WakeFrame0Row
	}
	return WakeFrame1Col, WakeFrame1Row
}

// GetWakeOverlay returns a wake sprite for the given animation frame.
func GetWakeOverlay(frame int) image.Image {
	col, row := WakeTileIndex(frame)
	if col < 0 {
		return nil
	}
	return getSailingOverlayTile(col, row)
}

var wakeFrameIndex int

// AdvanceWakeAnimation toggles the 2-frame wake cycle on game ticks.
func AdvanceWakeAnimation() {
	wakeFrameIndex = (wakeFrameIndex + 1) % 2
}

// CurrentWakeFrame returns the active wake animation frame (0 or 1).
func CurrentWakeFrame() int {
	return wakeFrameIndex
}

func getSailingOverlayTile(col, row int) image.Image {
	if !HasSailingVisualsTileset() {
		return nil
	}
	tile := extractTileAt(col, row)
	if tile == nil || !hasOpaquePixels(tile) {
		return nil
	}
	return tile
}

func hasOpaquePixels(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0x8000 {
				return true
			}
		}
	}
	return false
}

// WakeAftPosition returns the cell behind a ship for wake placement.
func WakeAftPosition(pos common.Coordinates, heading common.Facing) common.Coordinates {
	return sailing.WakeAftPosition(pos, heading)
}
