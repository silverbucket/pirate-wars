package world

import (
	"image"
	"image/color"
	"image/draw"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/terrain"
	"pirate-wars/cmd/window"
)

// generateMinimapImage rasterises the whole world once, at world init.
func (world *MapView) generateMinimapImage() {
	world.logger.Info("Generating minimap")
	cols := common.WorldCols
	rows := common.WorldRows
	cellWidth := float32(window.MiniMapArea.Width) / float32(cols)
	cellHeight := float32(window.MiniMapArea.Height) / float32(rows)

	world.minimap = world.createRawMapImage(cellWidth, cellHeight, cols, rows, window.MiniMapArea.Width, window.MiniMapArea.Height)
}

// createRawMapImage paints one pixel block per world cell, coloured by terrain.
func (world *MapView) createRawMapImage(cellWidth, cellHeight float32, cols, rows int, imageWidth, imageHeight int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			for y := int(float32(r) * cellHeight); y < int(float32(r+1)*cellHeight); y++ {
				for x := int(float32(c) * cellWidth); x < int(float32(c+1)*cellWidth); x++ {
					img.Set(x, y, terrain.GetColor(world.terrain.Cells[c][r]))
				}
			}
		}
	}
	return img
}

// MinimapBase returns the cached terrain raster, without any markers on it.
//
// The markers are drawn over this on the GPU each frame instead of being baked
// in. The old path copied the whole 700x700 image and re-uploaded it as a
// texture on every tick — four times a second while the map is open — to move
// one dot by less than three pixels.
func (world *MapView) MinimapBase() *image.RGBA {
	return world.minimap
}

// MinimapCellSize returns the minimap pixels per world cell, so callers can place
// markers and the viewport outline in the same space as the raster.
func (world *MapView) MinimapCellSize() (float32, float32) {
	return float32(window.MiniMapArea.Width) / float32(common.WorldCols),
		float32(window.MiniMapArea.Height) / float32(common.WorldRows)
}

// MinimapImage renders the whole-world minimap with player and town markers.
func (world *MapView) MinimapImage(pos common.Coordinates, entities entities.ViewableEntities) *image.RGBA {
	cols := common.WorldCols
	rows := common.WorldRows

	// Create a copy of the base image
	img := image.NewRGBA(world.minimap.Rect)
	draw.Draw(img, img.Bounds(), world.minimap, image.Point{}, draw.Src)

	// Calculate pixel position of player on minimap
	cellWidth := float32(window.MiniMapArea.Width) / float32(cols)
	cellHeight := float32(window.MiniMapArea.Height) / float32(rows)

	// overlays can be anything that implements ViewableEntity (towns, player)
	overlays := []MinimapOverlay{}
	overlays = append(overlays, MinimapOverlay{pos: pos, color: color.White})

	for _, e := range entities {
		overlays = append(overlays, MinimapOverlay{pos: e.GetPos(), color: e.GetColor()})
	}

	dotSize := 5
	for _, item := range overlays {
		x := int(float32(item.pos.X) * cellWidth)
		y := int(float32(item.pos.Y) * cellHeight)
		for dy := -dotSize / 2; dy <= dotSize/2; dy++ {
			for dx := -dotSize / 2; dx <= dotSize/2; dx++ {
				px := x + dx
				py := y + dy
				if px >= 0 && px < window.MiniMapArea.Width && py >= 0 && py < window.MiniMapArea.Height {
					img.Set(px, py, item.color)
				}
			}
		}
	}

	return img
}
