package world

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/entities"
	"pirate-wars/cmd/terrain"
	"pirate-wars/cmd/window"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (world *MapView) generateMinimapImage() {
	world.logger.Info("Generating minimap")
	cols := common.WorldCols
	rows := common.WorldRows
	cellWidth := float32(window.MiniMapArea.Width) / float32(cols)
	cellHeight := float32(window.MiniMapArea.Height) / float32(rows)

	world.minimap = world.createRawMapImage(cellWidth, cellHeight, cols, rows, window.MiniMapArea.Width, window.MiniMapArea.Height)
}

func (world *MapView) createRawMapImage(cellWidth, cellHeight float32, cols, rows int, imageWidth, imageHeight int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	fmt.Print(".")
	count := 0

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			for y := int(float32(r) * cellHeight); y < int(float32(r+1)*cellHeight); y++ {
				for x := int(float32(c) * cellWidth); x < int(float32(c+1)*cellWidth); x++ {
					img.Set(x, y, terrain.GetColor(world.terrain.Cells[c][r]))
					count++
					if count%2000 == 0 {
						fmt.Print(".")
					}
				}
			}
		}
	}
	fmt.Println("done")
	return img
}

func (world *MapView) getMinimapWithOverlays(pos common.Coordinates, entities entities.ViewableEntities) *image.RGBA {
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

func (world *MapView) ShowMinimapPopup(pos common.Coordinates, entities entities.ViewableEntities, w fyne.Window) {
	minimapPopup = widget.NewModalPopUp(
		container.NewStack(
			canvas.NewImageFromImage(world.getMinimapWithOverlays(pos, entities)),
		),
		w.Canvas(),
	)
	minimapPopup.Resize(fyne.NewSize(float32(window.MiniMapArea.Width), float32(window.MiniMapArea.Height)))
	minimapPopup.Move(
		fyne.NewPos(float32(window.Window.Width-window.MiniMapArea.Width)/2,
			float32(window.Window.Height-window.MiniMapArea.Height)/2),
	)
	minimapPopup.Show()
}

func (world *MapView) HideMinimapPopup() {
	if minimapPopup != nil {
		minimapPopup.Hide()
	}
}
