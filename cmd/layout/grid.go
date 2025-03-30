package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"math"
)

var viewableAreaGridCols = 66

// gridLayoutNoPadding defines a grid layout with no padding between cells
type gridLayoutNoPadding struct {
	cols int
}

// newGridLayoutNoPadding creates a new instance of the custom grid layout
func newGridLayoutNoPadding(cols int) fyne.Layout {
	return &gridLayoutNoPadding{cols: cols}
}

// Layout arranges the cells in a grid with no padding
func (g *gridLayoutNoPadding) Layout(cells []fyne.CanvasObject, size fyne.Size) {
	if g.cols < 1 {
		g.cols = 1 // Ensure at least one column
	}

	// Count visible cells
	visibleCount := 0
	for _, obj := range cells {
		if obj.Visible() {
			visibleCount++
		}
	}

	if visibleCount == 0 {
		return
	}

	// Add a small overlap to eliminate gaps (e.g., 1 pixel)
	overlap := float32(1.0)

	// Calculate rows based on visible cells and columns
	rows := (visibleCount + g.cols - 1) / g.cols // Ceiling division

	// Calculate cell size: width = available width / cols, height = available height / rows
	// Use the smaller value to ensure square cells
	cellWidth := size.Width / float32(g.cols)
	cellHeight := size.Height / float32(rows)
	cellSize := float32(math.Min(float64(cellWidth), float64(cellHeight)))

	// Position each square cell without padding
	i := 0
	for row := 0; row < rows && i < visibleCount; row++ {
		for col := 0; col < g.cols && i < visibleCount; col++ {
			if !cells[i].Visible() {
				i++
				continue
			}

			x := float32(col) * cellSize
			y := float32(row) * cellSize
			cells[i].Move(fyne.NewPos(x, y))
			//cells[i].Resize(fyne.NewSize(cellSize, cellSize))
			cells[i].Resize(fyne.NewSize(cellSize+overlap, cellSize+overlap))
			i++
		}
	}
}

// MinSize calculates the minimum size required for the layout
func (g *gridLayoutNoPadding) MinSize(cells []fyne.CanvasObject) fyne.Size {
	if g.cols < 1 {
		g.cols = 1
	}

	// Find the maximum min size of individual cells (assuming square)
	maxCellSize := float32(0)
	visibleCount := 0
	for _, obj := range cells {
		if obj.Visible() {
			minSize := obj.MinSize()
			cellSize := float32(math.Max(float64(minSize.Width), float64(minSize.Height)))
			if cellSize > maxCellSize {
				maxCellSize = cellSize
			}
			visibleCount++
		}
	}

	if visibleCount == 0 {
		return fyne.NewSize(0, 0)
	}

	// Calculate rows based on visible objects
	rows := (visibleCount + g.cols - 1) / g.cols
	return fyne.NewSize(maxCellSize*float32(g.cols), maxCellSize*float32(rows))
}

func CreateGridContainer(cells []fyne.CanvasObject) *fyne.Container {
	return container.New(
		newGridLayoutNoPadding(viewableAreaGridCols),
		cells...,
	)
}

func GetCellList() []fyne.CanvasObject {
	cg := GetColGridDimensions()
	return make([]fyne.CanvasObject, cg.Width*cg.Height)
}

func GetColGridDimensions() Dimensions {
	return Dimensions{
		Width: viewableAreaGridCols,
		//Height: int(float32(viewableAreaGridCols) * 0.75),
		Height: int(float32(viewableAreaGridCols) * (float32(ViewableArea.Height) / float32(ViewableArea.Width))),
	}
}
